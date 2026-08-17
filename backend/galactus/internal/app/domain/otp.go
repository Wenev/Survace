package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type OTPData struct {
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	Claimed   bool   `json:"claimed"`
}

func (o *OTPData) GetCreatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, o.CreatedAt)
}

func (o *OTPData) GetExpiresTime() (time.Time, error) {
	return time.Parse(time.RFC3339, o.ExpiresAt)
}

func (o *OTPData) IsExpired() bool {
	expiresTime, err := time.Parse(time.RFC3339, o.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expiresTime)
}

func (o *OTPData) TimeRemaining() time.Duration {
	expiresTime, err := time.Parse(time.RFC3339, o.ExpiresAt)
	if err != nil {
		return 0
	}
	return time.Until(expiresTime)
}

func (o *OTPData) MarkAsClaimed() {
	o.Claimed = true
}

func GenerateOTP() (*OTPData, error) {
	var maxInt int64 = 1000000
	n, err := rand.Int(rand.Reader, big.NewInt(maxInt))
	if err != nil {
		return nil, err
	}
	otpCode := fmt.Sprintf("%06d", n.Int64())
	now := time.Now()
	return &OTPData{
		Code:      otpCode,
		CreatedAt: now.Format(time.RFC3339),
		ExpiresAt: now.Add(time.Duration(10) * time.Minute).Format(time.RFC3339),
		Claimed:   false,
	}, nil
}
