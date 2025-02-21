package models

import "time"

type Health struct {
	WorkerID int `json:"worker_id"`
	Healthy bool `json:"healthy"`
	TaskCount uint64 `json:"task_count"`
	ReportTime time.Time `json:"report_time"`
}
