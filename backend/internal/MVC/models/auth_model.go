package models

import (
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// AuthModel owns all persistence for the user/auth domain. The DB handle is
// assigned by the service (am.DB = as.DB) before any method is called.
type AuthModel struct {
	DB *gorm.DB
}

// FirstOrCreateByTelegramID looks up a user by Telegram id and creates one with
// the supplied profile fields if none exists. On return, user is populated with
// the persisted row (including its generated ID).
func (m AuthModel) FirstOrCreateByTelegramID(user *db.User) error {
	return m.DB.
		Where(db.User{TelegramID: user.TelegramID}).
		Attrs(db.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			PhotoURL:  user.PhotoURL,
		}).
		FirstOrCreate(user).Error
}

// FindByID loads a single user by primary key.
func (m AuthModel) FindByID(id uint) (*db.User, error) {
	var user db.User
	if err := m.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
