package usecase

import (
	"context"
	"time"

	"github.com/rimantoro/event_driven/profiler/entities/client"
	"github.com/rimantoro/event_driven/profiler/entities/client/model"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/helpers"
)

type usecase struct {
	streamRepo     client.StreamRepository
	contextTimeout time.Duration
}

func NewUsecase(repo client.StreamRepository, to time.Duration) client.Usecase {
	return &usecase{
		streamRepo:     repo,
		contextTimeout: to,
	}
}

func (u *usecase) PostStreamRegisterClient(c context.Context, id string) error {
	client := model.Client{
		ID:     id,
		Key:    helpers.RandString(32),
		Status: 1,
	}
	_, err := u.streamRepo.StreamForCreate(bootstrap.App.Kafka.Producer, &client)
	return err
}
