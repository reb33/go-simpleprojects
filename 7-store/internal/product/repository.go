package product

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
)


type Product struct {
	gorm.Model
	Name        string
	Description string
	Images      pq.StringArray `gorm:"type:text[]"`
}
