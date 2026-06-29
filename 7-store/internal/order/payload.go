package order

type CreateOrderRequest struct {
	Products []string `json:"products"`
}
