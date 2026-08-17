package domain

import (
	"crypto/sha256"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

var pasetoV2 = paseto.NewV2()

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(userID int32, timeLimit time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(timeLimit),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return status.Error(codes.PermissionDenied, "Unauthorized Token")
	}
	return nil
}

func CreateToken(userID int32, timeLimit time.Duration) (string, error) {
	payload, err := NewPayload(userID, timeLimit)
	if err != nil {
		return "", err
	}
	key := os.Getenv("PASETO_SECRET")
	byteK := sha256.Sum256([]byte(key))
	var keyByte []byte = byteK[:]
	return pasetoV2.Encrypt(keyByte, payload, nil)
}

func VerifyToken(token string) (*Payload, error) {
	key := os.Getenv("PASETO_SECRET")
	byteK := sha256.Sum256([]byte(key))
	var keyByte []byte = byteK[:]
	payload := &Payload{}

	err := pasetoV2.Decrypt(token, keyByte, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func CreateRefreshToken(userID int32, timeLimit time.Duration) (string, error) {
	payload, err := NewPayload(userID, timeLimit)
	if err != nil {
		return "", err
	}
	key := os.Getenv("PASETO_REFRESH_SECRET")
	byteK := sha256.Sum256([]byte(key))
	var keyByte []byte = byteK[:]
	return pasetoV2.Encrypt(keyByte, payload, nil)
}

func VerifyRefreshToken(token string) (*Payload, error) {
	key := os.Getenv("PASETO_REFRESH_SECRET")
	byteK := sha256.Sum256([]byte(key))
	var keyByte []byte = byteK[:]
	payload := &Payload{}

	err := pasetoV2.Decrypt(token, keyByte, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
