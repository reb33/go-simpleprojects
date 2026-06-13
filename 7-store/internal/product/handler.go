package product

import "net/http"

type ProductHandler struct{}

func NewProductHandler(router *http.ServeMux) {
	handler := &ProductHandler{}

	router.HandleFunc("/product", handler.GetProduct())
}

func (handler *ProductHandler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Product"))
	}
}
