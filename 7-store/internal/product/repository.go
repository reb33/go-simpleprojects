package product

import (
	"demo-store/pkg/db"
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductDB struct {
	gorm.Model
	Name        string
	Description string
	Images      pq.StringArray `gorm:"type:text[]"`
}

func (ProductDB) TableName() string {
    return "products"
}

func NewProductDB(product Product) *ProductDB {
	return &ProductDB{
		Name:        product.Name,
		Description: product.Description,
		Images:      pq.StringArray(product.Images),
	}
}

type ProductRepository struct {
	*db.Db
}

func NewProductRepository(db *db.Db) *ProductRepository {
	return &ProductRepository{Db: db}
}

func (repo *ProductRepository) Create(product *Product) (*ProductDB, error) {
	productDB := NewProductDB(*product)
	result := repo.Db.Create(productDB)
	if result.Error != nil {
		return nil, result.Error
	}
	return productDB, nil
}

func (repo *ProductRepository) Update(id uint64, product *Product) (*ProductDB, error) {
	productDB := &ProductDB{
		Model: gorm.Model{ID: uint(id)},
		Name:        product.Name,
		Description: product.Description,
		Images:      pq.StringArray(product.Images),
	}
	result := repo.Db.Clauses(clause.Returning{}).Updates(productDB)
	if result.Error != nil {
		return nil, result.Error
	}
	return productDB, nil
}

func (repo *ProductRepository) Delete(id uint64) (*ProductDB, bool, error) {
	var productDB ProductDB
	result := repo.Db.Clauses(clause.Returning{}).Delete(&productDB, id)
	if result.Error != nil {
		return nil, false, result.Error
	}
	return &productDB, result.RowsAffected > 0, nil
}

func (repo *ProductRepository) Get(id uint64) (*ProductDB, error) {
	var productDB ProductDB
	result := repo.Db.Where(id).First(&productDB)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, result.Error
	}
	return &productDB, nil
}