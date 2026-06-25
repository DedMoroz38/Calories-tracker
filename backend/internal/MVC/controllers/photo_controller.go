package controllers

import (
	"io"
	"strconv"

	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PhotoController struct {
	BaseController
}

// readUpload pulls the multipart "file" field out of the request and returns
// its bytes and declared content type.
func (pc PhotoController) readUpload(c *fiber.Ctx) ([]byte, string, *errors.APIError) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, "", errors.BadRequest("missing 'file' upload")
	}
	f, err := fileHeader.Open()
	if err != nil {
		return nil, "", errors.BadRequest("could not read uploaded file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, "", errors.BadRequest("could not read uploaded file")
	}
	return data, fileHeader.Header.Get("Content-Type"), nil
}

func (pc PhotoController) service(c *fiber.Ctx) services.PhotoService {
	return services.PhotoService{DB: c.Locals("gorm").(*gorm.DB)}
}

// Create handles POST /api/v1/photos (multipart). Posts a new photo.
func (pc PhotoController) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	data, contentType, apiErr := pc.readUpload(c)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}

	photo, apiErr := pc.service(c).CreatePhoto(userID, contentType, data)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{Data: photo})
}

// Mine handles GET /api/v1/photos/me. The current user's own photos.
func (pc PhotoController) Mine(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	photos, apiErr := pc.service(c).ListMine(userID)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}
	return c.JSON(dto.BaseResponse{Data: photos})
}

// Delete handles DELETE /api/v1/photos/:id.
func (pc PhotoController) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return pc.Fail(c, errors.BadRequest("invalid photo id"))
	}
	if apiErr := pc.service(c).DeletePhoto(userID, uint(id)); apiErr != nil {
		return pc.Fail(c, apiErr)
	}
	return c.JSON(dto.BaseResponse{Message: "photo deleted"})
}

// Feed handles GET /api/v1/feed?cursor=&limit=. Other users' photos, newest
// first, keyset-paginated.
func (pc PhotoController) Feed(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var cursor uint
	if s := c.Query("cursor"); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return pc.Fail(c, errors.BadRequest("cursor must be a positive integer"))
		}
		cursor = uint(v)
	}
	limit := 12
	if s := c.Query("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			limit = v
		}
	}

	feed, apiErr := pc.service(c).Feed(userID, cursor, limit)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}
	return c.JSON(dto.BaseResponse{Data: feed})
}

// SetAvatar handles POST /api/v1/profile/avatar (multipart). Sets the avatar.
func (pc PhotoController) SetAvatar(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	data, contentType, apiErr := pc.readUpload(c)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}

	url, apiErr := pc.service(c).SetAvatar(userID, contentType, data)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}
	return c.JSON(dto.BaseResponse{Data: fiber.Map{"avatar_url": url}})
}
