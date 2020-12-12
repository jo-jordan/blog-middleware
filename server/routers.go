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
		"Pull Notes From GitHub",
		[]string{"POST", "OPTIONS"},
		"/github/pull",
		api.Pull,
	},
	Route{
		"Get Notes",
		[]string{"GET", "OPTIONS"},
		"/blogs",
		api.BlogHandler,
	},
	Route{
		"Get Notes",
		[]string{"GET", "OPTIONS"},
		"/blogs/{category}",
		api.BlogHandler,
	},
	Route{
		"Get Notes",
		[]string{"GET", "OPTIONS"},
		"/blogs/{category}/{id}",
		api.BlogHandler,
	},
	Route{
		"Network - Get Remote IP",
		[]string{"GET", "OPTIONS"},
		"/network/ip",
		api.GetIP,
	},
}
