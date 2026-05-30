package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	NewRandomHandler(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server started on port 8080")
	server.ListenAndServe()
}
