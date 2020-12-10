package server

import (
	"github.com/lzjlxebr/blog-middleware/api"
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
	Route{
		"Network - Get Remote IP",
		[]string{"GET", "OPTIONS"},
		"/network/ip",
		api.GetIP,
	},
}
