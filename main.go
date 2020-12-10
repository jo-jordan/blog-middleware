package main

import (
	"github.com/lzjlxebr/blog-middleware/server"
	"log"
	"net/http"
)

func main() {
	router := server.NewRouter()

	log.Fatal(http.ListenAndServe("127.0.0.1:8888", router))
}
