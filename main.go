package main

import (
	"github.com/gorilla/mux"
	"github.com/lzjlxebr/blog-middleware/server"
	"log"
	"net/http"
)

func main() {
	router := server.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(server.Filter)

	log.Fatal(http.ListenAndServe("127.0.0.1:8888", router))
}
