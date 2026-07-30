package main

import (
	"demo-store/configs"
	"demo-store/internal/auth"
	"demo-store/internal/order"
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

func App() http.Handler{
	config := configs.LoadConfigs()
	db := db.NewDb(config)
	jwt := jwt.NewJWT(config.Auth.Secret)

	// Repository
	productRepository := product.NewProductRepository(db)
	authRepository := auth.NewAuthRepository(db)
	orderRepository := order.NewOrderRepository(db)

	// Service
	authService := auth.NewService(jwt, authRepository)
	orderService := order.NewOrderService(order.OrderServiceDeps{
		OrderRepository: orderRepository,
		UserRepository:       authRepository,
		ProductRepository:    productRepository,
	})

	// Handlers
	router := http.NewServeMux()

	product.NewProductHandler(router, productRepository, config)
	auth.NewHandler(router, authService)
	order.NewOrderHandler(router, &order.OrderHandlerDeps{
		OrderService: orderService,
		Config:       config,
	})
	return middleware.Logs(router)
}

func main() {
	
	app := App()
	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	fmt.Println("Server listen on port 8080")
	server.ListenAndServe()
}
