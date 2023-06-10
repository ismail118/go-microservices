package main

import (
	"github.com/go-chi/chi/v5"
	"testing"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	routes := testApp.routes()
	chiRoutes := routes.(chi.Routes)

	routesForSearch := []string{"/", "/log-grpc", "/handle"}
	for _, r := range routesForSearch {
		isRouteExist(t, r, chiRoutes)
	}
}

func isRouteExist(t *testing.T, route string, chiRoute chi.Routes) {
	found := false

	for _, r := range chiRoute.Routes() {
		if r.Pattern == route {
			found = true
		}
	}

	if !found {
		t.Errorf("failed route %s not found", route)
	}
}
