package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"calorie-counter/internal/common/dto"
)

// VerifyTelegramAuth validates a Telegram Login Widget payload: it recomputes
// the HMAC-SHA256 over the sorted key=value pairs (secret = SHA256(botToken))
// and rejects data older than 24h.
func VerifyTelegramAuth(data dto.TelegramAuthRequest, botToken string) error {
	fields := map[string]string{
		"id":         strconv.FormatInt(data.ID, 10),
		"first_name": data.FirstName,
		"last_name":  data.LastName,
		"username":   data.Username,
		"photo_url":  data.PhotoURL,
		"auth_date":  strconv.FormatInt(data.AuthDate, 10),
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if fields[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, fields[k]))
	}
	dataCheckString := strings.Join(parts, "\n")

	h := sha256.New()
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedHash), []byte(data.Hash)) {
		return errors.New("invalid telegram auth hash")
	}

	if time.Now().Unix()-data.AuthDate > 86400 {
		return errors.New("telegram auth data expired")
	}

	return nil
}
