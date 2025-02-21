package main

import (
	"context"

	aNats "github.com/dusk-chancellor/dc-agent/adapters/nats"
	"github.com/dusk-chancellor/dc-agent/config"
	"github.com/dusk-chancellor/dc-agent/utils"
	"github.com/dusk-chancellor/dc-agent/worker"
	"github.com/dusk-chancellor/dc-agent/zaplog"
	"go.uber.org/zap"
)

func main() {
	// init logger
	log := zaplog.New()
	
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic("failed to load config", zap.Error(err))
	}

	// set up connection w/nats
	nc, err := aNats.Connect(cfg.NatsURL)
	if err != nil {
		log.Panic("failed to connect to nats", zap.Error(err), zap.String("nats url", cfg.NatsURL))
	}
	defer nc.Close()
	
	// init utilities
	utils := utils.New(log)

	// init context
	ctx := context.Background()

	// create worker pool
	pool := worker.NewWorkerPool(
		nc,
		log,
		utils,
		cfg.NumWorkers,
		cfg.TasksSubject,
		cfg.ResultsSubject,
		cfg.HealthSubject,
		cfg.HealthTick,
	)

	// start pool
	pool.Start(ctx)

	// keep main thread alive
	select {}
}
