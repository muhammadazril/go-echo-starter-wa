package svc

import (
	"log"
	"os"
	"os/signal"

	// "strings"

	"github.com/rimantoro/event_driven/profiler/shared/worker"
)

func WorkerStart() {

	log.Println(bannerWorker)

	pool := worker.NewWorker()
	defer func() {
		pool.Stop()
		worker.GetRedisPool().Close()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
