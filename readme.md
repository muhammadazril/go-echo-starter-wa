# Event Driven Architecture

I'm just trying to use Kafka as backbone for several service, just proof of concept.

## Profiler

This is some place you want to store your user/client profile, validate it if this is correct clients

## GoWaMessage

    ```
        $ go to profile dirrectory, then run
        go run main.go restapi
        to start the API, and
        go run main.go worker 
        to start the Worker
        in config file can set path the whatsapp session to be placed.
    ```