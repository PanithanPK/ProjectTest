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

func (s *UserService) Register(ctx context.Context, req user.RegisterRequest) (*user.MeResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.BankAccount = strings.TrimSpace(req.BankAccount)

	if req.Username == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.BankAccount == "" {
		return nil, errors.New("missing required fields")
	}
	if !isTenDigits(req.BankAccount) {
		return nil, ErrInvalidBankAccount
	}

	if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameTaken
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if _, err := s.repo.GetByBankAccount(ctx, req.BankAccount); err == nil {
		return nil, ErrBankAccountTaken
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Username:     req.Username,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BankAccount:  req.BankAccount,
		Credit:       1000,
	}

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return &user.MeResponse{
		ID:          id,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		BankAccount: u.BankAccount,
		Credit:      u.Credit,
	}, nil
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

func (s *UserService) Me(ctx context.Context, userID uint64) (*user.MeResponse, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.MeResponse{
		ID:          u.ID,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		BankAccount: u.BankAccount,
		Credit:      u.Credit,
	}, nil
}

func (s *UserService) Update(ctx context.Context, userID uint64, req user.UpdateRequest) (*user.MeResponse, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		u.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		u.LastName = strings.TrimSpace(*req.LastName)
	}
	if req.BankAccount != nil {
		bank := strings.TrimSpace(*req.BankAccount)
		if bank == "" {
			return nil, errors.New("bank_account cannot be empty")
		}
		if !isTenDigits(bank) {
			return nil, ErrInvalidBankAccount
		}
		existing, err := s.repo.GetByBankAccount(ctx, bank)
		if err == nil && existing.ID != u.ID {
			return nil, ErrBankAccountTaken
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		u.BankAccount = bank
	}
	if req.Password != nil {
		p := *req.Password
		if strings.TrimSpace(p) == "" {
			return nil, errors.New("password cannot be empty")
		}
		hash, err := utils.HashPassword(p)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = hash
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return &user.MeResponse{
		ID:          u.ID,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		BankAccount: u.BankAccount,
		Credit:      u.Credit,
	}, nil
}
