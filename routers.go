package main

import (
	"blog-middleware/api"
	"net/http"
)

type Route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Sync From GitHub",
		[]string{"GET", "OPTIONS"},
		"/sync",
		api.Pull,
	},
	Route{
		"Initial",
		[]string{"GET", "OPTIONS"},
		"/init",
		api.Init,
	},
}
