package product

type CreateProductRequest struct{
	Name        string
	Description string
	Images      []string
}

type UpdateProductRequest struct{
	Name        string
	Description string
	Images      []string
}