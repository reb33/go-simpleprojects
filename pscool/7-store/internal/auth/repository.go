package auth

import (
	"demo-store/internal/model"
	"demo-store/internal/order"
	"demo-store/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type UserDB struct {
	gorm.Model
	Phone     string
	SessionId string
	Code      string
	Orders    []order.OrderDB `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (UserDB) TableName() string {
	return "users"
}

type AuthRepository struct {
	db *db.Db
}

func NewAuthRepository(db *db.Db) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) Upsert(phone, sessionId, code string) (*UserDB, error) {
	var phoneDB UserDB
	err := repo.db.Where("phone = ?", phone).First(&phoneDB).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	phoneDB.Phone = phone
	phoneDB.SessionId = sessionId
	phoneDB.Code = code

	if err := repo.db.Save(&phoneDB).Error; err != nil {
		return nil, err
	}
	return &phoneDB, nil
}

func (repo *AuthRepository) Update(id uint, phoneData model.User, code string) (*UserDB, error) {
	var phoneDB UserDB
	if err := repo.db.First(&phoneDB, id).Error; err != nil {
		return nil, err
	}
	phoneDB.Phone = phoneData.Phone
	phoneDB.SessionId = phoneData.SessionId
	phoneDB.Code = code

	if err := repo.db.Save(&phoneDB).Error; err != nil {
		return nil, err
	}
	return &phoneDB, nil
}

func (repo *AuthRepository) GetByPhone(phone string) (*model.User, error) {
	var phoneDB UserDB
	if err := repo.db.Where("phone = ?", phone).First(&phoneDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	user := model.User{
		Id:        int(phoneDB.ID),
		Phone:     phoneDB.Phone,
		SessionId: phoneDB.SessionId,
	}
	return &user, nil
}

func (repo *AuthRepository) GetBySessionId(sessionId string) (*UserDB, error) {
	var phoneDB UserDB
	if err := repo.db.Where("session_id = ?", sessionId).First(&phoneDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &phoneDB, nil
}
