package main

import (
	"log"
	"os"

	"github.com/rimantoro/event_driven/profiler/shared/svc"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 1 {
		switch args[0] {
		case "consumer":
			svc.StartConsumer("general")
		case "producer":
			log.Print("Stream producer here")
		case "restapi":
			svc.StartRestAPI()
		case "cli":
			svc.DoCliProcess(args)
		case "migrate":
			log.Print("migration mode here")
		}
	} else {
		log.Print("helper command here")
	}
}
