package main

import (
	"fmt"
	"log"
	"net/http"
)

const defaultPort = "80"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting server on port %s.", defaultPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", defaultPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
