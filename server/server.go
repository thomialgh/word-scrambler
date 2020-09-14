package server

import (
	"net/http"
	"time"
)

// Server -
func Server() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	server.ListenAndServe()
}
