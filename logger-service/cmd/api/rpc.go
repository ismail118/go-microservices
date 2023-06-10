package main

import (
	"fmt"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct {
	Repo data.Repository
}

type RPCPayload struct {
	Name string
	Data string
}

func NewRpcServer(r data.Repository) *RPCServer {
	rpcSrv := new(RPCServer)
	rpcSrv.Repo = r
	return rpcSrv
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	entry := data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := r.Repo.Insert(entry)
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	*resp = fmt.Sprintf("Process payload via RPC: %s", payload.Name)

	return nil
}
