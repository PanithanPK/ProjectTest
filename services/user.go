package services

import (
	"ProjectTest/config"
	"ProjectTest/modules/user"
	"ProjectTest/repositorys"
	"ProjectTest/utils"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrBankAccountTaken   = errors.New("bank_account already exists")
	ErrInvalidBankAccount = errors.New("bank_account must be 10 digits")
)

type UserService struct {
	repo   *repositorys.UserRepository
	jwtCfg config.JWTConfig
}

func NewUserService(repo *repositorys.UserRepository, jwtCfg config.JWTConfig) *UserService {
	return &UserService{repo: repo, jwtCfg: jwtCfg}
}

func (s *UserService) Login(ctx context.Context, req user.LoginRequest) (*user.TokenResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("missing required fields")
	}

	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPassword(u.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	exp := now.Add(s.jwtCfg.Expire)
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     exp.Unix(),
		"iat":     now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtCfg.Secret)
	if err != nil {
		return nil, err
	}

	return &user.TokenResponse{AccessToken: signed, TokenType: "Bearer", ExpiresIn: int64(s.jwtCfg.Expire.Seconds())}, nil
}
