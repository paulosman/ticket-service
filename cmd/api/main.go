package main

import (
	"log"
	"net/http"
	"time"

	"github.com/paulosman/ticket-service/handlers"
)

func main() {

	router := handlers.NewRouter()

	server := &http.Server{
		Addr:        "0.0.0.0:9000",
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
