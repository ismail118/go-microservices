package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WriteLog(t *testing.T) {
	reqBody := map[string]interface{}{
		"name": "some name",
		"data": "some data",
	}
	reqJson, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/log", bytes.NewReader(reqJson))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testApp.WriteLog)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
	}
}
