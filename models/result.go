package models

import "time"

type Result struct {
	WorkerID 	int 	  `json:"worker_id"`
	TaskID 		string 	  `json:"task_id"`
	Output 		float64   `json:"output"`
	ProcessedAt time.Time `json:"processed_at"`
}
