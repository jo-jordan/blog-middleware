package main

import (
	http2 "blog-middleware/http"
	"log"
	"net/http"
)

func main() {
	router := http2.NewRouter()

	log.Fatal(http.ListenAndServe("127.0.0.1:8888", router))
}
