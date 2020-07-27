package svc

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	_clientStream "github.com/rimantoro/event_driven/profiler/entities/client/delivery/stream"
	_clientMdl "github.com/rimantoro/event_driven/profiler/entities/client/model"
	_clientRepo "github.com/rimantoro/event_driven/profiler/entities/client/repository"
	_clientUcase "github.com/rimantoro/event_driven/profiler/entities/client/usecase"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
)

func StartConsumer(process string) {

	bannerWorker += `
Process Name       ` + process + `
`

	log.Println(bannerWorker)

	//======= INITIATE USECASES =========

	timeoutContext := time.Duration(30) * time.Second

	clientRepo := _clientRepo.NewClientRepository()
	clientUcase := _clientUcase.NewUsecase(clientRepo, timeoutContext)

	usecase := AllUsecaseStruct{
		ClientUsecase: clientUcase,
	}

	bootstrap.App.Logger.Info("Loaded Usecases", zap.Any("usecase", usecase))

	c := bootstrap.App.Kafka.Consumer

	defer func() {
		c.Close()
	}()

	switch process {
	case "general":
		_clientStream.ConsumeTestMessages(c, []string{_clientMdl.Topic})
	default:
		c.Close()
		bootstrap.App.Logger.Error("please define process name to be processed.")
		os.Exit(1)
	}
}
