package auth

import (
	"errors"
	"net/http"
	"some-project/pkg/request"
	"some-project/pkg/response"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service, router *http.ServeMux) {
	handler := &Handler{
		Service: service,
	}
	router.HandleFunc("POST /auth", handler.Auth())
	router.HandleFunc("POST /verify", handler.Verify())
}

func (h *Handler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[AuthRequest](r)
		if err != nil {
			response.Json(w, http.StatusBadRequest, ErrResponse{Error: err.Error()})
			return
		}
		sessionId := h.Service.AddPhone(body.Phone)

		response.Json(w, http.StatusOK, AuthResponse{
			SessionId: sessionId,
		})
	}
}

func (h *Handler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[VerifyRequest](r)
		if err != nil {
			response.Json(w, http.StatusBadRequest, ErrResponse{Error: err.Error()})
			return
		}
		token, err := h.Service.Verify(body.SessionId, body.Code)
		if err != nil {
			if errors.Is(err, ErrInvalidCode) || errors.Is(err, ErrInvalidSessionId) {
				response.Json(w, http.StatusBadRequest, ErrResponse{Error: err.Error()})
				return
			}
			response.Json(w, http.StatusInternalServerError, ErrResponse{Error: err.Error()})
			return
		}
		response.Json(w, http.StatusOK, VerifyResponse{
			Token: token,
		})
	}
}
