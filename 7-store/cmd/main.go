package main

import (
	"demo-store/configs"
	"demo-store/internal/auth"
	"demo-store/internal/product"
	"demo-store/pkg/db"
	"demo-store/pkg/jwt"
	"demo-store/pkg/middleware"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	config := configs.LoadConfigs()
	db := db.NewDb(config)
	jwt := jwt.NewJWT(config.Auth.Secret)
	productRepository := product.NewProductRepository(db)
	authRepository := auth.NewAuthRepository(db)

	authService := auth.NewService(jwt, authRepository)

	router := http.NewServeMux()

	product.NewProductHandler(router, productRepository)
	auth.NewHandler(router, authService)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logs(router),
	}

	fmt.Println("Server listen on port 8080")
	server.ListenAndServe()
}
