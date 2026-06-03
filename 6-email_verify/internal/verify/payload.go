package verify

type SendRequestPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type Status string

const (
	StatusOk    Status = "ok"
	StatusError Status = "error"
)

type ResponsePayload struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
