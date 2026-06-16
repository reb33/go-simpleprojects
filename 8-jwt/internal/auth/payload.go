package auth

type AuthRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type AuthResponse struct {
	SessionId string `json:"sessionId"`
}

type VerifyRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type VerifyResponse struct {
	Token string `json:"token"`
}

type ErrResponse struct {
	Error string `json:"error"`
}