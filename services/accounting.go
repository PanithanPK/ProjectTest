package services

import (
	"ProjectTest/modules/accounting"
	"ProjectTest/repositorys"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type AccountingService struct {
	db       *sql.DB
	userRepo *repositorys.UserRepository
	repo     *repositorys.AccountingRepository
}

func NewAccountingService(db *sql.DB, userRepo *repositorys.UserRepository, repo *repositorys.AccountingRepository) *AccountingService {
	return &AccountingService{db: db, userRepo: userRepo, repo: repo}
}

func (s *AccountingService) Transfer(ctx context.Context, fromUserID uint64, req accounting.TransferRequest) (uint64, error) {
	req.BankAccount = strings.TrimSpace(req.BankAccount)
	if req.BankAccount == "" {
		return 0, errors.New("bank_account is required")
	}
	if !isTenDigits(req.BankAccount) {
		return 0, ErrInvalidBankAccount
	}
	if req.Amount <= 0 {
		return 0, ErrInvalidAmount
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	fromUser, err := s.userRepo.GetByIDForUpdate(ctx, tx, fromUserID)
	if err != nil {
		return 0, err
	}
	toUser, err := s.userRepo.GetByBankAccountForUpdate(ctx, tx, req.BankAccount)
	if err != nil {
		return 0, err
	}
	if toUser.ID == fromUser.ID {
		return 0, errors.New("cannot transfer to yourself")
	}

	if fromUser.Credit < req.Amount {
		return 0, ErrInsufficientBalance
	}

	if err := s.userRepo.UpdateCredit(ctx, tx, fromUser.ID, fromUser.Credit-req.Amount); err != nil {
		return 0, err
	}
	if err := s.userRepo.UpdateCredit(ctx, tx, toUser.ID, toUser.Credit+req.Amount); err != nil {
		return 0, err
	}

	transferID, err := s.repo.CreateTransfer(ctx, tx, fromUser.ID, toUser.ID, req.Amount)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return transferID, nil
}

func (s *AccountingService) TransferList(ctx context.Context, userID uint64, startDate string, endDate string) ([]accounting.TransferItem, error) {
	var start *time.Time
	var end *time.Time

	if strings.TrimSpace(startDate) != "" {
		t, err := parseDate(startDate, false)
		if err != nil {
			return nil, errors.New("invalid start_date")
		}
		start = &t
	}
	if strings.TrimSpace(endDate) != "" {
		t, err := parseDate(endDate, true)
		if err != nil {
			return nil, errors.New("invalid end_date")
		}
		end = &t
	}

	rows, err := s.repo.ListTransfers(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	items := make([]accounting.TransferItem, 0, len(rows))
	for _, r := range rows {
		direction := "out"
		if r.ToUserID == userID {
			direction = "in"
		}
		items = append(items, accounting.TransferItem{
			ID:              r.ID,
			FromBankAccount: r.FromBankAccount,
			ToBankAccount:   r.ToBankAccount,
			Amount:          r.Amount,
			CreatedAt:       r.CreatedAt.UTC().Format(time.RFC3339),
			Direction:       direction,
		})
	}
	return items, nil
}
