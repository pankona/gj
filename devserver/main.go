package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve files from the "public" directory
	fs := http.FileServer(http.Dir("public"))

	// Register the file server as the handler for all requests
	http.Handle("/", fs)

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
