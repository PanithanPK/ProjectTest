package config

import (
	"errors"
	"os"
	"time"
)

type JWTConfig struct {
	Secret []byte
	Expire time.Duration
}

func LoadJWTConfig() (JWTConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return JWTConfig{}, errors.New("JWT_SECRET is required")
	}

	expireRaw := os.Getenv("JWT_EXPIRE")
	if expireRaw == "" {
		expireRaw = "24h"
	}
	expire, err := time.ParseDuration(expireRaw)
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{Secret: []byte(secret), Expire: expire}, nil
}
