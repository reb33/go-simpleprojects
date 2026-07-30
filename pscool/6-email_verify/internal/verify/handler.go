package verify

import (
	"email-verify/configs"
	"email-verify/internal/service"
	"email-verify/pkg/request"
	"email-verify/pkg/response"
	"log"
	"net/http"
	"runtime/debug"
)

type VerifyHandler struct {
	*configs.Config
	*service.EmailService
}

type VeryfyHandlerDeps struct {
	*configs.Config
	*service.EmailService
}

func NewVerifyHandler(router *http.ServeMux, deps VeryfyHandlerDeps) {
	handler := &VerifyHandler{
		Config:       deps.Config,
		EmailService: deps.EmailService,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[SendRequestPayload](r)
		if err != nil {
			response.Json(w, 400, ResponsePayload{
				Status:  string(StatusError),
				Message: err.Error(),
			})
			return
		}
		err = handler.EmailService.SendEmail(payload.Email)
		if err != nil {
			log.Printf("ERROR: %v\nStack Trace:\n%s", err, debug.Stack())
			response.Json(w, 500, ResponsePayload{
				Status:  string(StatusError),
				Message: "Internal Error",
			})
			return
		}
		response.Json(w, 200, ResponsePayload{
			Status: string(StatusOk),
		})
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		err := handler.EmailService.VerifyEmail(hash)
		if err != nil {
			response.Json(w, 400, ResponsePayload{
				Status:  string(StatusError),
				Message: err.Error(),
			})
			return
		}
		response.Json(w, 200, ResponsePayload{
			Status: string(StatusOk),
		})
	}
}
