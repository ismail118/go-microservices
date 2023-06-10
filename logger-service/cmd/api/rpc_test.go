package main

import (
	"fmt"
	"testing"
)

func Test_rpc_writeLog(t *testing.T) {
	rpcSrv := NewRpcServer(testApp.Repo)
	var logReq RPCPayload = RPCPayload{
		Name: "some log",
		Data: "some data",
	}
	var res string
	err := rpcSrv.LogInfo(logReq, &res)
	if err != nil {
		t.Errorf("failed error %s", err.Error())
	}

	if res != fmt.Sprintf("Process payload via RPC: %s", logReq.Name) {
		t.Errorf("failed wrong response: %s", res)
	}
}
