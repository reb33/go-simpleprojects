package order

import "demo-store/pkg/di"

type OrderServiceDeps struct {
	OrderRepository   *OrderRepository
	ProductRepository di.IProductRepository
	UserRepository    di.IUserRepository
}

type OrderService struct {
	OrderRepository   *OrderRepository
	ProductRepository di.IProductRepository
	UserRepository    di.IUserRepository
}

func NewOrderService(deps OrderServiceDeps) *OrderService {
	return &OrderService{
		OrderRepository:   deps.OrderRepository,
		ProductRepository: deps.ProductRepository,
		UserRepository:    deps.UserRepository,
	}
}

func (s *OrderService) CreateOrder(phone string, productNames []string) (*OrderDB, error) {
	user, err := s.UserRepository.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	products, err := s.ProductRepository.GetAllByNames(productNames)
	if err != nil {
		return nil, err
	}
	return s.OrderRepository.Create(user, products)
}

func (s *OrderService) GetOrder(id int, phone string) (*OrderDB, error) {
	user, err := s.UserRepository.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	order, err := s.OrderRepository.Get(id)
	if err != nil {
		return nil, err
	}
	if int(order.UserID) != user.Id {
		return nil, ErrNotMatchOrderId
	}
	return order, nil
}

func (s *OrderService) GetAllOrders(phone string) ([]OrderDB, error) {
	user, err := s.UserRepository.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	return s.OrderRepository.GetAllByUser(user)
}