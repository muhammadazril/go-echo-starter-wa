# Event Driven Architecture

I'm just trying to use Kafka as backbone for several service, just proof of concept.

## Profiler

This is some place you want to store your user/client profile, validate it if this is correct clients

## GoWaMessage

## This app use packages:
- [go-whatsapp](https://github.com/Rhymen/go-whatsapp) 
- [phonenumber normalize](https://github.com/dongri/phonenumber) 

    ```
        $ go to profile dirrectory, then run
        go run main.go restapi
        to start the API, and
        go run main.go worker 
        to start the Worker
        in config file can set path the whatsapp session to be placed.
    ```

## Run the app :
To start apps run API and worker in terminal.
Then hit url POST wa/messagejob with json body
{
"message":"type your message",
"number":"+6281234567890"
}
Worker will show QR code in terminal for the first time session to be stored.
![QR code](https://github.com/muhammadazril/go-echo-starter-wa/QRcode.png)
Scan the QR code in your Whatsapp application
![Scan the QR code](https://github.com/muhammadazril/go-echo-starter-wa/Scan_QR_code.jpg)