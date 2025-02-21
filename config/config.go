package config

import (
	"errors"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var ErrEmptySubject = errors.New("tasks, results and health subject names required in config")

type Config struct {
//	Env 		   string 		 `env:"env"`
	NatsURL		   string		 `env:"NATS_URL"`
	NumWorkers 	   int 			 `env:"NUM_WORKERS"`
	TasksSubject   string 		 `env:"TASKS_SUBJECT"`
	ResultsSubject string 		 `env:"RESULTS_SUBJECT"`
	HealthSubject  string 		 `env:"HEALTH_SUBJECT"`
	HealthTick 	   time.Duration `env:"HEALTH_TICK"`
}


func LoadConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.TasksSubject == "" || cfg.ResultsSubject == "" || cfg.HealthSubject == "" {
		return nil, ErrEmptySubject
	}


	if cfg.NumWorkers == 0 {
		cfg.NumWorkers = 4
	}

	if cfg.HealthTick.String() == "" {
		cfg.HealthTick = 30 * time.Second
	}

	return &cfg, nil
}
