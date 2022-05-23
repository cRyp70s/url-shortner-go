package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Starting server.\n Listening on :9001\n")
	if err := http.ListenAndServe(":9001", createHandler()); err != nil {
		log.Fatal("failed to start server", err)
	}
}
