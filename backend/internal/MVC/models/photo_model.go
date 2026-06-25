package models

import (
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// PhotoModel owns all persistence for the photo/feed domain. The DB handle is
// assigned by the service before any method is called.
type PhotoModel struct {
	DB *gorm.DB
}

// Create inserts a new photo. On return, photo carries its generated ID and
// CreatedAt.
func (m PhotoModel) Create(photo *db.Photo) error {
	return m.DB.Create(photo).Error
}

// ListByUser returns a user's own photos, newest first.
func (m PhotoModel) ListByUser(userID uint) ([]db.Photo, error) {
	var photos []db.Photo
	err := m.DB.Where("user_id = ?", userID).Order("id desc").Find(&photos).Error
	return photos, err
}

// FeedItem is a photo joined with its author's display fields, used to build
// the public feed without a second query per row.
type FeedItem struct {
	db.Photo
	FirstName    string
	Username     string
	UserPhotoURL string
	UserAvatar   string
}

// Feed returns photos from everyone except excludeUserID, newest first.
// When cursor > 0 only rows with id < cursor are returned (keyset pagination),
// capped at limit rows.
func (m PhotoModel) Feed(excludeUserID uint, cursor uint, limit int) ([]FeedItem, error) {
	var items []FeedItem
	q := m.DB.
		Table("photos").
		Select("photos.*, users.first_name, users.username, users.photo_url AS user_photo_url, users.avatar_key AS user_avatar").
		Joins("JOIN users ON users.id = photos.user_id").
		Where("photos.user_id <> ?", excludeUserID)
	if cursor > 0 {
		q = q.Where("photos.id < ?", cursor)
	}
	err := q.Order("photos.id desc").Limit(limit).Scan(&items).Error
	return items, err
}

// FindByIDAndUser loads a single photo scoped to its owner, so a caller can
// read its S3 key before deleting. Returns gorm.ErrRecordNotFound when absent.
func (m PhotoModel) FindByIDAndUser(id, userID uint) (*db.Photo, error) {
	var p db.Photo
	if err := m.DB.Where("id = ? AND user_id = ?", id, userID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// DeleteByIDAndUser removes a single photo scoped to its owner and reports how
// many rows were affected (0 means it did not exist for this user).
func (m PhotoModel) DeleteByIDAndUser(id, userID uint) (int64, error) {
	res := m.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&db.Photo{})
	return res.RowsAffected, res.Error
}

// SetAvatarKey updates the user's avatar S3 key.
func (m PhotoModel) SetAvatarKey(userID uint, key string) error {
	return m.DB.Model(&db.User{}).Where("id = ?", userID).Update("avatar_key", key).Error
}
