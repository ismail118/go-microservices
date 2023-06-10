package main

import (
	"broker/logs"
	"broker/models"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"strings"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_Broker(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testApp.Broker)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed wrong response code, want %d got %d", http.StatusOK, rr.Code)
	}
}

var dataTestHandleSubmission = []struct {
	name            string
	reqPayload      models.RequestPayload
	codeExpectation int
}{
	{
		name: "success-auth",
		reqPayload: models.RequestPayload{
			Action: "auth",
			Auth: models.AuthPayload{
				Email:    "some@email.com",
				Password: "verysecret",
			},
		},
		codeExpectation: http.StatusAccepted,
	},
	{
		name: "success-log",
		reqPayload: models.RequestPayload{
			Action: "log",
			Log: models.LogPayload{
				Name: "some log",
				Data: "some data",
			},
		},
		codeExpectation: http.StatusAccepted,
	},
	{
		name: "success-email",
		reqPayload: models.RequestPayload{
			Action: "mail",
			Mail: models.MailPayload{
				From:    "me@gmail.com",
				To:      "you@gmail.com",
				Subject: "some subject",
				Message: "some message",
			},
		},
		codeExpectation: http.StatusAccepted,
	},
	{
		name: "default",
		reqPayload: models.RequestPayload{
			Action: "default",
		},
		codeExpectation: http.StatusBadRequest,
	},
}

func Test_HandleSubmission(t *testing.T) {
	jsonToReturn := `
{
	"error": false,
	"message": "some message"
}
`

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(jsonToReturn))}
	})

	testApp.Client = client

	startRPCServer(t)

	for _, x := range dataTestHandleSubmission {
		jsonBody, _ := json.Marshal(x.reqPayload)
		sBody := string(jsonBody)

		req, _ := http.NewRequest(http.MethodPost, "/handle", strings.NewReader(sBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.HandleSubmission)
		handler.ServeHTTP(rr, req)

		if rr.Code != x.codeExpectation {
			t.Errorf("failed at %s wrong response code, want %d got %d", x.name, x.codeExpectation, rr.Code)
		}
	}
}

var dataTestLogViaGRPC = struct {
	name            string
	reqPayload      models.RequestPayload
	codeExpectation int
}{
	name: "success-log-grpc",
	reqPayload: models.RequestPayload{
		Action: "log",
		Log: models.LogPayload{
			Name: "some log",
			Data: "some data",
		},
	},
	codeExpectation: http.StatusAccepted,
}

func Test_LogViaGRPC(t *testing.T) {
	startGRPCServer(t)

	jsonBody, _ := json.Marshal(dataTestLogViaGRPC.reqPayload)
	sBody := string(jsonBody)
	req, _ := http.NewRequest(http.MethodPost, "/log-grpc", strings.NewReader(sBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.LogViaGRPC)
	handler.ServeHTTP(rr, req)

	if rr.Code != dataTestLogViaGRPC.codeExpectation {
		t.Errorf("failed wrong response code, want %d got %d", dataTestLogViaGRPC.codeExpectation, rr.Code)
	}

}

type RPCServer int
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	*resp = fmt.Sprintf("Process payload via RPC: %s", payload.Name)
	return nil
}
func startRPCServer(t *testing.T) {
	err := rpc.Register(new(RPCServer))
	if err != nil {
		t.Errorf("failed register testRpcServer")
	}
	go func() {
		log.Println("Starting rpc server on port", rpcPort)
		listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
		if err != nil {
			t.Errorf("failed listen")
		}
		defer listen.Close()

		for {
			rpcConn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			go rpc.ServeConn(rpcConn)
		}
	}()
}

type LogServer struct {
	logs.UnimplementedLogServiceServer
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

func startGRPCServer(t *testing.T) {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
		if err != nil {
			t.Errorf("failed listen grpc %s", err.Error())
		}

		srv := grpc.NewServer()

		logs.RegisterLogServiceServer(srv, &LogServer{})

		log.Printf("Grpc server started on port %s\n", gRpcPort)

		err = srv.Serve(lis)
		if err != nil {
			t.Errorf("failed to server grpc %s", err.Error())
		}
	}()
}
