package main

import (
	"email-verify/configs"
	"email-verify/internal/repository"
	"email-verify/internal/service"
	"email-verify/internal/verify"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	repo := repository.NewEmailRepository("store.json")
	service := service.NewEmailService(conf, repo)
	verify.NewVerifyHandler(router, verify.VeryfyHandlerDeps{
		Config: conf,
		EmailService: service,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
