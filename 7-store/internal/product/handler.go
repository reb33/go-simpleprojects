package product

import (
	"demo-store/configs"
	"demo-store/internal/model"
	"demo-store/pkg/middleware"
	"demo-store/pkg/request"
	"demo-store/pkg/response"
	"log"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	repo *ProductRepository
}

func NewProductHandler(router *http.ServeMux, repo *ProductRepository, config *configs.Configs) {
	handler := &ProductHandler{
		repo: repo,
	}

	router.Handle("GET /product/{id}", middleware.IsAuthed(handler.GetProduct(), config))
	router.HandleFunc("POST /product", handler.CreateProduct())
	router.HandleFunc("PATCH /product/{id}", handler.UpdateProduct())
	router.HandleFunc("DELETE /product/{id}", handler.DeleteProduct())
}

func (handler *ProductHandler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string); ok {
			log.Printf("Phone: %s", phone)
		} else {
			log.Printf("Phone not found")
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
		}
		product, err := handler.repo.Get(uint64(id))
		if err != nil {
			if err == ErrProductNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, http.StatusOK, product)
	}
}

func (handler *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := request.HandleBody[CreateProductRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		product, err := handler.repo.Create(&model.Product{
			Name:        request.Name,
			Description: request.Description,
			Images:      request.Images,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		response.Json(w, http.StatusCreated, product)
	}
}

func (handler *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
		}
		request, err := request.HandleBody[UpdateProductRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		product, err := handler.repo.Update(uint64(id), &model.Product{
			Name:        request.Name,
			Description: request.Description,
			Images:      request.Images,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		response.Json(w, http.StatusOK, product)
	}
}

func (handler *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
		}
		product, deleted, err := handler.repo.Delete(uint64(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if !deleted {
			http.Error(w, "Product not found", http.StatusNotFound)
		}
		response.Json(w, http.StatusOK, product)
	}
}
