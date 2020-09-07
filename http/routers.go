package http

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
		[]string{"POST", "OPTIONS"},
		"/github/sync",
		api.Pull,
	},
	Route{
		"Initial",
		[]string{"GET", "OPTIONS"},
		"/github/init",
		api.Init,
	},
}
