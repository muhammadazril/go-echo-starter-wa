package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dongri/phonenumber"
	"github.com/rimantoro/event_driven/profiler/entities/gowa"
)

type usecase struct {
	Repo           gowa.Repository
	contextTimeout time.Duration
}

func NewUsecase(repo gowa.Repository, to time.Duration) gowa.Usecase {
	return &usecase{
		Repo:           repo,
		contextTimeout: to,
	}
}

func (u *usecase) SendMessage(c context.Context, number string, message string) (string, error) {

	number = phonenumber.Parse(number, "ID")
	fmt.Println("isi no hp dari usecase ", number)
	return u.Repo.SendMessage(c, number, message)
}
