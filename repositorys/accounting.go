package repositorys

import (
	"context"
	"database/sql"
	"time"
)

type AccountingRepository struct {
	db *sql.DB
}

func NewAccountingRepository(db *sql.DB) *AccountingRepository {
	return &AccountingRepository{db: db}
}

func (r *AccountingRepository) CreateTransfer(ctx context.Context, tx *sql.Tx, fromUserID uint64, toUserID uint64, amount int64) (uint64, error) {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO transfers (from_user_id, to_user_id, amount)
		VALUES (?, ?, ?)
	`, fromUserID, toUserID, amount)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

type TransferRow struct {
	ID              uint64
	FromBankAccount string
	ToBankAccount   string
	Amount          int64
	CreatedAt       time.Time
	FromUserID      uint64
	ToUserID        uint64
}

func (r *AccountingRepository) ListTransfers(ctx context.Context, userID uint64, start *time.Time, end *time.Time) ([]TransferRow, error) {
	query := `
		SELECT t.id, fu.bank_account, tu.bank_account, t.amount, t.created_at, t.from_user_id, t.to_user_id
		FROM transfers t
		JOIN users fu ON fu.id = t.from_user_id
		JOIN users tu ON tu.id = t.to_user_id
		WHERE (t.from_user_id = ? OR t.to_user_id = ?)
	`
	args := []any{userID, userID}

	if start != nil {
		query += " AND t.created_at >= ?"
		args = append(args, *start)
	}
	if end != nil {
		query += " AND t.created_at <= ?"
		args = append(args, *end)
	}
	query += " ORDER BY t.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []TransferRow
	for rows.Next() {
		var rrow TransferRow
		if err := rows.Scan(&rrow.ID, &rrow.FromBankAccount, &rrow.ToBankAccount, &rrow.Amount, &rrow.CreatedAt, &rrow.FromUserID, &rrow.ToUserID); err != nil {
			return nil, err
		}
		res = append(res, rrow)
	}
	return res, rows.Err()
}
