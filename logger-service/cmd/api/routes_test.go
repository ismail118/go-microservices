package main

import (
	"github.com/go-chi/chi/v5"
	"testing"
)

func Test_route_exist(t *testing.T) {
	routes := testApp.routes()
	chiRoutes := routes.(chi.Routes)

	routesToSearch := []string{"/log"}
	for _, r := range routesToSearch {
		isRouteExist(t, r, chiRoutes)
	}
}

func isRouteExist(t *testing.T, route string, chiRoutes chi.Routes) {
	found := false

	for _, r := range chiRoutes.Routes() {
		if r.Pattern == route {
			found = true
		}
	}

	if !found {
		t.Errorf("failed can't found route %s", route)
	}
}
