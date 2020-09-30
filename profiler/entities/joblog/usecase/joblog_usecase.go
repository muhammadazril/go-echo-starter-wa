package usecase

import (
	"context"
	"time"

	"github.com/rimantoro/event_driven/profiler/entities/joblog"
	joblog_model "github.com/rimantoro/event_driven/profiler/entities/joblog/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type joblogUsecase struct {
	joblogRepo     joblog.Repository
	contextTimeout time.Duration
}

func NewUsecase(r joblog.Repository, to time.Duration) joblog.Usecase {
	return &joblogUsecase{
		joblogRepo:     r,
		contextTimeout: to,
	}
}

func (u *joblogUsecase) InsertOne(c context.Context, m *joblog_model.JobLog) (*joblog_model.JobLog, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := u.joblogRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *joblogUsecase) FindOne(c context.Context, id string) (*joblog_model.JobLog, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.joblogRepo.FindOneJobID(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *joblogUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]joblog_model.JobLog, int64, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, count, err := u.joblogRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (u *joblogUsecase) UpdateOne(c context.Context, m *joblog_model.JobLog, id string) (*joblog_model.JobLog, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.joblogRepo.UpdateOne(ctx, m, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
