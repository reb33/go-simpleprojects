package verify

import (
	"email-verify/configs"
	"net/http"
)

type VerifyHandler struct{
	*configs.Config
}

type VeryfyHandlerDeps struct {
	*configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VeryfyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /sent", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sent"))
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("verify"))
	}
}
