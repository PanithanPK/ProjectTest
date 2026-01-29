package repositorys

import (
	"ProjectTest/modules/user"
	"context"
	"database/sql"
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, first_name, last_name, bank_account, credit, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FirstName, &u.LastName, &u.BankAccount, &u.Credit, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByBankAccount(ctx context.Context, bankAccount string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, first_name, last_name, bank_account, credit, created_at, updated_at
		FROM users
		WHERE bank_account = ?
	`, bankAccount).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FirstName, &u.LastName, &u.BankAccount, &u.Credit, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, first_name, last_name, bank_account, credit)
		VALUES (?, ?, ?, ?, ?, ?)
	`, u.Username, u.PasswordHash, u.FirstName, u.LastName, u.BankAccount, u.Credit)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, first_name, last_name, bank_account, credit, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FirstName, &u.LastName, &u.BankAccount, &u.Credit, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = ?, first_name = ?, last_name = ?, bank_account = ?
		WHERE id = ?
	`, u.PasswordHash, u.FirstName, u.LastName, u.BankAccount, u.ID)
	return err
}
