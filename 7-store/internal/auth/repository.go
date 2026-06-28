package auth

import (
	"demo-store/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type UserDB struct {
	gorm.Model
	Phone     string
	SessionId string
	Code      string
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
	phoneDB, err := repo.GetByPhone(phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if phoneDB == nil {
		phoneDB = &UserDB{
			Phone:     phone,
			SessionId: sessionId,
			Code:      code,
		}
		if err := repo.db.Create(phoneDB).Error; err != nil {
			return nil, err
		}
		return phoneDB, nil
	}
	if err := repo.db.Save(&phoneDB).Error; err != nil {
		return nil, err
	}
	return phoneDB, nil
}

func (repo *AuthRepository) Update(id uint, phoneData User) (*UserDB, error) {
	var phoneDB UserDB
	if err := repo.db.First(&phoneDB, id).Error; err != nil {
		return nil, err
	}
	phoneDB.Phone = phoneData.Phone
	phoneDB.SessionId = phoneData.SessionId
	phoneDB.Code = phoneData.Code

	if err := repo.db.Save(&phoneDB).Error; err != nil {
		return nil, err
	}
	return &phoneDB, nil
}

func (repo *AuthRepository) GetByPhone(phone string) (*UserDB, error) {
	var phoneDB UserDB
	if err := repo.db.Where("phone = ?", phone).First(&phoneDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &phoneDB, nil
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
