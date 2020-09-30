package joblog

import (
	"context"

	joblog_model "github.com/rimantoro/event_driven/profiler/entities/joblog/model"
)

type Usecase interface {
	InsertOne(ctx context.Context, u *joblog_model.JobLog) (*joblog_model.JobLog, error)
	FindOne(ctx context.Context, id string) (*joblog_model.JobLog, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]joblog_model.JobLog, int64, error)
	UpdateOne(ctx context.Context, joblog *joblog_model.JobLog, id string) (*joblog_model.JobLog, error)
}
