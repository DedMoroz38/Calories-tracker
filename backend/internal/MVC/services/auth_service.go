package services

import (
	stderrors "errors"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/config"
	"calorie-counter/internal/db"
	"calorie-counter/internal/util"

	"gorm.io/gorm"
)

// AuthService is the use-case layer for authentication. It embeds *gorm.DB
// (the embedded-field variant described in ARCHITECTURE_FLOW.md §5), so the
// controller assigns the handle with as.DB = c.Locals("gorm").(*gorm.DB).
type AuthService struct {
	*gorm.DB
}

// TelegramLogin verifies a Telegram Login Widget payload, upserts the user, and
// returns a signed JWT for the session.
func (s AuthService) TelegramLogin(req dto.TelegramAuthRequest) (string, *errors.APIError) {
	if err := util.VerifyTelegramAuth(req, config.Values.TelegramBotToken); err != nil {
		return "", errors.Unauthorized("invalid telegram authentication")
	}

	am := models.AuthModel{}
	am.DB = s.DB

	user := &db.User{
		TelegramID: req.ID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Username:   req.Username,
		PhotoURL:   req.PhotoURL,
	}
	if err := am.FirstOrCreateByTelegramID(user); err != nil {
		return "", errors.Internal("could not persist user")
	}

	token, err := util.GenerateJWT(user.ID, config.Values.JWTSecret)
	if err != nil {
		return "", errors.Internal("could not generate token")
	}
	return token, nil
}

// GetUser loads the profile for an authenticated user. It also checks whether
// the user has completed onboarding (i.e. has a UserProfile row).
func (s AuthService) GetUser(id uint) (*dto.MeResponse, *errors.APIError) {
	am := models.AuthModel{}
	am.DB = s.DB

	user, err := am.FindByID(id)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user not found")
		}
		return nil, errors.Internal("could not load user")
	}

	ps := ProfileService{DB: s.DB}
	onboarded := ps.IsOnboarded(id)

	return &dto.MeResponse{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		PhotoURL:   user.PhotoURL,
		CreatedAt:  user.CreatedAt,
		Onboarded:  onboarded,
	}, nil
}

