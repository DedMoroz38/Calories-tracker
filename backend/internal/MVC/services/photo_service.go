package services

import (
	"context"
	"fmt"
	"time"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/db"
	"calorie-counter/internal/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PhotoService is the use-case layer for the photo/feed domain.
type PhotoService struct {
	DB *gorm.DB
}

// presignTTL is how long a generated view URL stays valid. The frontend loads
// images immediately, so a short window is plenty.
const presignTTL = time.Hour

// maxImageBytes caps a single upload at 8 MB.
const maxImageBytes = 8 << 20

// allowedImageTypes maps accepted content types to a file extension.
var allowedImageTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
	"image/gif":  "gif",
}

// validateImage checks the declared content type and size.
func validateImage(contentType string, size int) (ext string, apiErr *errors.APIError) {
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return "", errors.BadRequest("unsupported image type; use jpeg, png, webp, or gif")
	}
	if size == 0 {
		return "", errors.BadRequest("empty file")
	}
	if size > maxImageBytes {
		return "", errors.BadRequest("image too large (max 8 MB)")
	}
	return ext, nil
}

// CreatePhoto stores an uploaded image in S3 and records it for the user.
func (s PhotoService) CreatePhoto(userID uint, contentType string, data []byte) (*dto.PhotoResponse, *errors.APIError) {
	if !storage.IsEnabled() {
		return nil, errors.Internal("photo storage is not configured")
	}
	ext, apiErr := validateImage(contentType, len(data))
	if apiErr != nil {
		return nil, apiErr
	}

	key := fmt.Sprintf("photos/%d/%s.%s", userID, uuid.NewString(), ext)
	if err := storage.Default.Upload(context.Background(), key, contentType, data); err != nil {
		return nil, errors.Internal("could not upload image")
	}

	pm := models.PhotoModel{DB: s.DB}
	photo := &db.Photo{UserID: userID, Key: key, CreatedAt: time.Now().UTC()}
	if err := pm.Create(photo); err != nil {
		// Best-effort cleanup of the orphaned object.
		_ = storage.Default.Delete(context.Background(), key)
		return nil, errors.Internal("could not save photo")
	}

	url, err := storage.Default.PresignGet(context.Background(), key, presignTTL)
	if err != nil {
		return nil, errors.Internal("could not generate image url")
	}
	return &dto.PhotoResponse{
		ID:        photo.ID,
		URL:       url,
		CreatedAt: photo.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ListMine returns the current user's own photos with presigned URLs.
func (s PhotoService) ListMine(userID uint) ([]dto.PhotoResponse, *errors.APIError) {
	if !storage.IsEnabled() {
		return nil, errors.Internal("photo storage is not configured")
	}
	pm := models.PhotoModel{DB: s.DB}
	photos, err := pm.ListByUser(userID)
	if err != nil {
		return nil, errors.Internal("could not load photos")
	}

	out := make([]dto.PhotoResponse, 0, len(photos))
	for _, p := range photos {
		url, err := storage.Default.PresignGet(context.Background(), p.Key, presignTTL)
		if err != nil {
			continue // skip a single bad object rather than failing the page
		}
		out = append(out, dto.PhotoResponse{
			ID:        p.ID,
			URL:       url,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}

// Feed returns a paginated page of other users' photos, newest first.
func (s PhotoService) Feed(userID uint, cursor uint, limit int) (*dto.FeedResponse, *errors.APIError) {
	if !storage.IsEnabled() {
		return nil, errors.Internal("photo storage is not configured")
	}
	if limit <= 0 || limit > 50 {
		limit = 12
	}

	pm := models.PhotoModel{DB: s.DB}
	items, err := pm.Feed(userID, cursor, limit)
	if err != nil {
		return nil, errors.Internal("could not load feed")
	}

	out := make([]dto.FeedItemResponse, 0, len(items))
	var last uint
	for _, it := range items {
		url, err := storage.Default.PresignGet(context.Background(), it.Key, presignTTL)
		if err != nil {
			continue
		}
		out = append(out, dto.FeedItemResponse{
			ID:           it.ID,
			URL:          url,
			UserID:       it.UserID,
			AuthorName:   authorName(it.FirstName, it.Username),
			AuthorAvatar: avatarURL(it.UserAvatar, it.UserPhotoURL),
			CreatedAt:    it.CreatedAt.Format(time.RFC3339),
		})
		last = it.ID
	}

	// Only advertise a next cursor when the page was full; a short page is the end.
	var next uint
	if len(items) == limit {
		next = last
	}
	return &dto.FeedResponse{Items: out, NextCursor: next}, nil
}

// DeletePhoto removes a photo (S3 object + row) owned by the user.
func (s PhotoService) DeletePhoto(userID, id uint) *errors.APIError {
	pm := models.PhotoModel{DB: s.DB}
	photo, err := pm.FindByIDAndUser(id, userID)
	if err != nil {
		return errors.NotFound("photo not found")
	}
	if storage.IsEnabled() {
		_ = storage.Default.Delete(context.Background(), photo.Key)
	}
	if _, err := pm.DeleteByIDAndUser(id, userID); err != nil {
		return errors.Internal("could not delete photo")
	}
	return nil
}

// SetAvatar uploads a new avatar image and points the user's profile at it.
func (s PhotoService) SetAvatar(userID uint, contentType string, data []byte) (string, *errors.APIError) {
	if !storage.IsEnabled() {
		return "", errors.Internal("photo storage is not configured")
	}
	ext, apiErr := validateImage(contentType, len(data))
	if apiErr != nil {
		return "", apiErr
	}

	key := fmt.Sprintf("avatars/%d/%s.%s", userID, uuid.NewString(), ext)
	if err := storage.Default.Upload(context.Background(), key, contentType, data); err != nil {
		return "", errors.Internal("could not upload avatar")
	}

	pm := models.PhotoModel{DB: s.DB}
	if err := pm.SetAvatarKey(userID, key); err != nil {
		_ = storage.Default.Delete(context.Background(), key)
		return "", errors.Internal("could not save avatar")
	}

	url, err := storage.Default.PresignGet(context.Background(), key, presignTTL)
	if err != nil {
		return "", errors.Internal("could not generate avatar url")
	}
	return url, nil
}

// avatarURL presigns an uploaded avatar key, falling back to the Telegram photo
// URL (already an absolute URL) when no avatar has been uploaded.
func avatarURL(avatarKey, telegramURL string) string {
	if avatarKey != "" && storage.IsEnabled() {
		if url, err := storage.Default.PresignGet(context.Background(), avatarKey, presignTTL); err == nil {
			return url
		}
	}
	return telegramURL
}

// authorName prefers the first name, then username, then a generic label.
func authorName(firstName, username string) string {
	switch {
	case firstName != "":
		return firstName
	case username != "":
		return username
	default:
		return "Someone"
	}
}
