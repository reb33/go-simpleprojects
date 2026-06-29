package di

import "demo-store/internal/model"

type IUserRepository interface {
	GetByPhone(phone string) (*model.User, error)
}

type IProductRepository interface {
	GetAllByNames(names []string) ([]model.Product, error)
}
