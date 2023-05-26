package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logsdb").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	*resp = fmt.Sprintf("Process payload via RPC: %s", payload.Name)

	return nil
}
