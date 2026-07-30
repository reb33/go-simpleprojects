package order

import (
	"demo-store/internal/model"
	"demo-store/internal/product"
	"demo-store/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type OrderDB struct {
	gorm.Model
	UserID   uint
	Products []product.ProductDB `gorm:"many2many:order_products;constraint:OnDelete:CASCADE;"`
}

func (OrderDB) TableName() string {
	return "orders"
}

type OrderRepository struct {
	db *db.Db
}

func NewOrderRepository(db *db.Db) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) Create(user *model.User, products []model.Product) (*OrderDB, error) {
	orderDB := &OrderDB{
		UserID: uint(user.Id),
	}
	repo.db.Create(&orderDB)
	
	var productDBs []product.ProductDB
	for _, p := range products {
		productDBs = append(productDBs, *product.NewProductDB(p))
	}
	err := repo.db.Model(&orderDB).Association("Products").Append(&productDBs)
	return orderDB, err
}

func (repo *OrderRepository) Get(orderId int) (*OrderDB, error) {
	var order OrderDB
	result := repo.db.Preload("Products").First(&order, orderId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, result.Error
	}
	// var products []model.Product
	// if err := repo.db.Model(&order).Association("Products").Find(&products); err != nil {
	// 	return nil, err
	// }
	return &order, nil
}

func (repo *OrderRepository) GetAllByUser(user *model.User) ([]OrderDB, error) {
	var orders []OrderDB
	result := repo.db.Preload("Products").Find(&orders, "user_id = ?", user.Id).Order("created_at desc")
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
