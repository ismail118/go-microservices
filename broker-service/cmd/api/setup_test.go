package main

import (
	"broker/event"
	"os"
	"testing"
)

var testApp Config

func TestMain(m *testing.M) {
	testEmitter := event.NewRabbitTestEmitter(nil)
	testApp.Emitter = testEmitter
	os.Exit(m.Run())
}
