package services

import (
	"ProjectTest/modules/user"
	"ProjectTest/repositorys"
	"ProjectTest/utils"
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrBankAccountTaken   = errors.New("bank_account already exists")
	ErrInvalidBankAccount = errors.New("bank_account must be 10 digits")
)

type UserService struct {
	repo *repositorys.UserRepository
}

func (s *UserService) Login(ctx context.Context, req user.LoginRequest) error {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return errors.New("missing required fields")
	}

	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}

	if !utils.CheckPassword(u.PasswordHash, req.Password) {
		return ErrInvalidCredentials
	}

	return nil
}
