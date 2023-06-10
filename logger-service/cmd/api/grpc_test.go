package main

import (
	"context"
	"log-service/logs"
	"testing"
)

func Test_grpc_writeLog(t *testing.T) {
	grpcTestSrv := LogServer{}
	grpcTestSrv.Repo = testApp.Repo

	req := &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: "some name",
			Data: "some data",
		},
	}
	res, err := grpcTestSrv.WriteLog(context.Background(), req)
	if err != nil {
		t.Errorf("failed error %s", err.Error())
	}

	if res.Result != "logged" {
		t.Errorf("failed wrong response, want %s got %s", "logged", res.Result)
	}
}
