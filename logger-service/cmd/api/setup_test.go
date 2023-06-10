package main

import (
	"log-service/data"
	"os"
	"testing"
)

var testApp Config

func TestMain(m *testing.M) {
	repo := data.NewMongoTestRepository(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}
