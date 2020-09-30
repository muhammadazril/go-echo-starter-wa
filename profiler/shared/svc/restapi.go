package svc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/models"
	"go.uber.org/zap"

	_clientHttp "github.com/rimantoro/event_driven/profiler/entities/client/delivery/http"
	_clientRepo "github.com/rimantoro/event_driven/profiler/entities/client/repository"
	_clientUcase "github.com/rimantoro/event_driven/profiler/entities/client/usecase"
	_gowaHttp "github.com/rimantoro/event_driven/profiler/entities/gowa/delivery/http"
	_gowaRepo "github.com/rimantoro/event_driven/profiler/entities/gowa/repository"
	_gowaUcase "github.com/rimantoro/event_driven/profiler/entities/gowa/usecase"
	_joblogHttp "github.com/rimantoro/event_driven/profiler/entities/joblog/delivery/http"
	_joblogRepo "github.com/rimantoro/event_driven/profiler/entities/joblog/repository"
	_joblogUcase "github.com/rimantoro/event_driven/profiler/entities/joblog/usecase"
)

func StartRestAPI() {

	log.Println(banner)

	//======= DISCONNECT ON EXIT =========

	defer func() {
		bootstrap.App.Redis.Close()
		if err := bootstrap.App.Mongo.Disconnect(context.Background()); err != nil {
			bootstrap.App.Logger.Error("fail disconnect mongo", zap.Error(err))
		}
		if err := bootstrap.App.Logger.Sync(); err != nil {
			bootstrap.App.Logger.Error("fail flushing zap logger", zap.Error(err))
		}
	}()

	//======= INITIATE USECASES =========

	timeoutContext := time.Duration(30) * time.Second

	clientRepo := _clientRepo.NewClientRepository()
	clientUcase := _clientUcase.NewUsecase(clientRepo, timeoutContext)
	gowaRepo := _gowaRepo.NewWaRepository()
	gowaUcase := _gowaUcase.NewUsecase(gowaRepo, timeoutContext)
	joblogRepo := _joblogRepo.NewMongoRepository(bootstrap.App.Mongo.Database(bootstrap.App.Config.GetString(models.ConfigMongoDatabase)))
	joblogUcase := _joblogUcase.NewUsecase(joblogRepo, timeoutContext)

	usecase := AllUsecaseStruct{
		ClientUsecase: clientUcase,
		GowaUsecase:   gowaUcase,
		JoblogUcase:   joblogUcase,
	}

	//======= INITIATE ECHO =========

	e := echo.New()
	e.Debug = true
	e.Server.Addr = fmt.Sprintf(":%d", bootstrap.App.Config.GetInt("app.port"))
	c := e.AcquireContext()
	c.Set("app_path", os.Getenv("APP_PATH"))

	//======= INITIATE HTTP HANDLER FOR EACH ENTITIES =========

	_clientHttp.NewHttpHandler(e, usecase.ClientUsecase)
	_gowaHttp.NewHttpHandler(e, usecase.GowaUsecase)
	// _joblogHttp.NewHttpHandler(e, usecase.JoblogUcase)
	_joblogHttp.NewJobLogHandler(e, usecase.JoblogUcase)

	//======= GRACEFULL SHUTDOWN FOR ECHO HTTP =========

	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
