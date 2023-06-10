package main

import (
	"broker/event"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	port     = "8080"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

type Config struct {
	Emitter event.Emitter
	Client  *http.Client
}

func main() {
	// try to connect to rabbitmq
	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitCon.Close()

	emitter := event.NewRabbitEmitter(rabbitCon)
	app := Config{
		Emitter: emitter,
		Client:  &http.Client{},
	}
	log.Printf("Start service on port %s\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMq not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("Connect to rabbitMq")
			break
		}
		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
	}

	return connection, nil
}
