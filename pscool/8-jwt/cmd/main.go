package main

import (
	"fmt"
	"net/http"
	"some-project/configs"
	"some-project/internal/auth"
	"some-project/pkg/jwt"
)

func main() {
	config := configs.LoadConfig()
	jwt := jwt.NewJWT(config.Secret)
	service := auth.NewService(jwt)
	

	mux := http.NewServeMux()
	auth.NewHandler(service, mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server started on port 8080")
	server.ListenAndServe()
}
