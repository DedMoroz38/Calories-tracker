package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"calorie-counter/internal/common/dto"
)

func TestJWTRoundTrip(t *testing.T) {
	const secret = "test-secret"

	token, err := GenerateJWT(42, secret)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	claims, err := ParseJWT(token, secret)
	if err != nil {
		t.Fatalf("ParseJWT: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("UserID = %d, want 42", claims.UserID)
	}
}

func TestParseJWTWrongSecret(t *testing.T) {
	token, _ := GenerateJWT(1, "right")
	if _, err := ParseJWT(token, "wrong"); err == nil {
		t.Fatal("expected error parsing token signed with a different secret")
	}
}

// signTelegram reproduces the widget's signing algorithm so the test can craft
// a payload that VerifyTelegramAuth must accept.
func signTelegram(data dto.TelegramAuthRequest, botToken string) string {
	dataCheckString := fmt.Sprintf("auth_date=%d\nfirst_name=%s\nid=%d",
		data.AuthDate, data.FirstName, data.ID)

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyTelegramAuth(t *testing.T) {
	const botToken = "123:abc"
	data := dto.TelegramAuthRequest{
		ID:        99,
		FirstName: "Egor",
		AuthDate:  time.Now().Unix(),
	}
	data.Hash = signTelegram(data, botToken)

	if err := VerifyTelegramAuth(data, botToken); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}

	tampered := data
	tampered.Hash = "deadbeef"
	if err := VerifyTelegramAuth(tampered, botToken); err == nil {
		t.Fatal("expected error for tampered hash")
	}

	expired := data
	expired.AuthDate = time.Now().Add(-48 * time.Hour).Unix()
	expired.Hash = signTelegram(expired, botToken)
	if err := VerifyTelegramAuth(expired, botToken); err == nil {
		t.Fatal("expected error for expired auth_date")
	}
}
