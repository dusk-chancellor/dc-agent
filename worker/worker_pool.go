package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dusk-chancellor/dc-agent/utils"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type WorkerPool struct {
	workers 	   []*Worker
	nc 			   *nats.Conn
	log			   *zap.Logger
	utils 		   utils.Utils

	numWorkers 	   int 			 // workers quantity
	tasksSubject   string 		 // for queue subscription
	resultsSubject string 		 // for publishing
	healthSubject  string		 // for checking health state
	healthTick 	   time.Duration // health tick frequency
}

func NewWorkerPool(
	nc *nats.Conn,
	log *zap.Logger,
	u utils.Utils,
	numWorkers int,
	tasks, results, health string,
	hTick time.Duration,
	) *WorkerPool {

	return &WorkerPool{
		workers: 		make([]*Worker, numWorkers),
		nc: 			nc,
		log:			log,
		utils: 			u,
		numWorkers: 	numWorkers,
		tasksSubject: 	tasks,
		resultsSubject: results,
		healthSubject: 	health,
		healthTick: 	hTick,
	}
}

// set up and start workers
func (wp *WorkerPool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 1; i <= wp.numWorkers; i++ {
		wp.workers[i] = NewWorker(
			i,
			wp.nc,
			wp.log,
			wp.utils,
			wp.tasksSubject,
			wp.resultsSubject,
			wp.healthSubject,
			"worker pool",
			wp.healthTick,
		)
		
		wg.Add(1)
		go wp.workers[i].Start(ctx, &wg)
	}

	// monitor for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT)

    go func() {
        <-sigChan
        wg.Wait()
        os.Exit(0)
    }()

    return
}
