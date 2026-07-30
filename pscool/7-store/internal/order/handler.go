package order

import (
	"demo-store/configs"
	"demo-store/pkg/middleware"
	"demo-store/pkg/request"
	"demo-store/pkg/response"
	"log"
	"net/http"
	"strconv"
)

type OrderHandlerDeps struct {
	OrderService *OrderService
	Config       *configs.Configs
}

type OrderHandler struct {
	OrderService *OrderService
}

func NewOrderHandler(router *http.ServeMux, deps *OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderService: deps.OrderService,
	}
	router.Handle("POST /order", middleware.IsAuthed(handler.CreateOrder(), deps.Config))
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.GetOrder(), deps.Config))
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.GetAllOrders(), deps.Config))
}

func (handler *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			log.Printf("Phone: %s", phone)
		} else {
			http.Error(w, "Phone not found", http.StatusInternalServerError)
			return
		}

		body, err := request.HandleBody[CreateOrderRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order, err := handler.OrderService.CreateOrder(phone, body.Products)
		if err != nil {
			log.Printf("Error creating order: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, http.StatusCreated, order)
	}
}

func (handler *OrderHandler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			log.Printf("Phone: %s", phone)
		} else {
			http.Error(w, "Phone not found", http.StatusInternalServerError)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		order, err := handler.OrderService.GetOrder(id, phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, http.StatusOK, order)
	}
}

func (handler *OrderHandler) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			log.Printf("Phone: %s", phone)
		} else {
			http.Error(w, "Phone not found", http.StatusInternalServerError)
			return
		}
		orders, err := handler.OrderService.GetAllOrders(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, http.StatusOK, orders)
	}
}
	
