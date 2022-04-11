package main

import (
	"log"
	"net/http"
)

const PORT = ":8081"

func main() {
	mux := routes()

	log.Println("starting web server on port", PORT)

	_ = http.ListenAndServe(PORT, mux) //run serve handling routes
}
