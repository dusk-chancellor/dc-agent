package worker

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dusk-chancellor/dc-agent/models"
	"github.com/dusk-chancellor/dc-agent/utils"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Worker struct {
	id 			   int
	nc 			   *nats.Conn
	log 		   *zap.Logger
	utils		   utils.Utils

	tasksSubject   string 		 // for queue subscription
	resultsSubject string		 // for publishing
	healthSubject  string		 // for reporting health state
	queueGroup 	   string		 // common group of workers
	healthTick 	   time.Duration

	taskCount 	   atomic.Uint64 
	healthy 	   atomic.Bool
}

func NewWorker(
	id int,
	nc *nats.Conn,
	log *zap.Logger,
	u utils.Utils,
	tasks, results, wHealth, q string,
	hTick time.Duration,
	) *Worker {

	w := &Worker{
		id: 			id,
		nc: 			nc,
		log: 			log,
		utils: 			u,
		tasksSubject: 	tasks,
		resultsSubject: results,
		healthSubject: 	wHealth,
		queueGroup: 	q,
		healthTick: 	hTick,
	}

	return w
}

// subscribes to tasks subject, works on it & publishes to results subject;
// supports tick-based health report
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	sub, err := w.nc.QueueSubscribe(
		w.tasksSubject,
		w.queueGroup,
		func(msg *nats.Msg) {
			w.work(msg)
		},
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	ticker := time.NewTicker(w.healthTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.healthReport()
		}
	}
}

// worker main job
func (w *Worker) work(msg *nats.Msg) {
	// proccess message data into custom struct
	var task models.Task

	if err := json.Unmarshal(msg.Data, &task); err != nil {
		w.log.Error("failed to unmarshal message",
		zap.Error(err),
		zap.Int("worker id", w.id),
		zap.ByteString("msg", msg.Data),
	)
		return
	}

	w.taskCount.Add(1)

	// convert into postfix queue
	queue, err := w.utils.ToPostfix(task.Payload)
	if err != nil {
		w.log.Error("failed to convert to postfix",
		zap.Error(err),
		zap.Int("worker id", w.id),
		zap.String("payload", task.Payload),
	)
		return
	}

	// compute
	res, err := w.utils.Evaluate(queue)
	if err != nil {
		w.log.Error("failed to evaluate",
		zap.Error(err),
		zap.Int("worker id", w.id),
		zap.String("queue", queue),
	)
		return
	}

	// wrap up the result
	var result = models.Result{
		WorkerID: 	 w.id,
		TaskID: 	 task.ID,
		Output: 	 res,
		ProcessedAt: time.Now(),
	}

	resultData, err := json.Marshal(result)
	if err != nil {
		w.log.Error("failed to marshal result",
		zap.Error(err),
		zap.Int("worker id", w.id), 
		zap.Any("result", result),
	)
		return
	}

	// publish data to results subject
	if err := w.nc.Publish(w.resultsSubject, resultData); err != nil {
		w.log.Error("failed to publish result",
		zap.Error(err),
		zap.Int("worker id", w.id),
	)
		w.healthy.Store(false)
		return
	}
}

// periodic health report
func (w *Worker) healthReport() {
	// load states
	healthy := w.healthy.Load()
	taskCount := w.taskCount.Load()

	// wrap up
	var health = models.Health{
		WorkerID: w.id,
		Healthy: healthy,
		TaskCount: taskCount,
		ReportTime: time.Now(),
	}

	data, err := json.Marshal(health)
	if err != nil {
		w.log.Error("failed to marshal data",
		zap.Error(err),
		zap.Any("health", health),
	)
		return
	}

	// publish data to health subject
	if err := w.nc.Publish(w.healthSubject, data); err != nil {
		w.log.Error("failed to publish health state",
		zap.Error(err),
		zap.Int("worker id", w.id),
	)
		return
	}
}
