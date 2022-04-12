package main

import (
	"log"
	"net/http"

	"github.com/SzymonSkursrki/go-ws-example/internal/handlers"
)

const PORT = ":8080"

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("starting web server on port", PORT)

	_ = http.ListenAndServe(PORT, mux) //run serve handling routes
}
