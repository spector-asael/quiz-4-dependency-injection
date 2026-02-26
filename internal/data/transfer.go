package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"github.com/spector-asael/lab4-crud/internal/validator"
)

type TransferModel struct {
	DB *sql.DB
}

type TransferInput struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        float64
}

// ValidateTransfer validates a transfer request
func ValidateTransfer(v *validator.Validator, input TransferInput) {
	v.Check(input.FromAccountID > 0, "from_account_id", "must be provided")
	v.Check(input.ToAccountID > 0, "to_account_id", "must be provided")
	v.Check(input.Amount > 0, "amount", "must be greater than 0")
}

// Create performs the transfer in a single DB transaction
func (m TransferModel) Create(input TransferInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var fromGL, toGL int64
	if err := tx.QueryRowContext(ctx, `SELECT gl_account_id FROM accounts WHERE id=$1`, input.FromAccountID).Scan(&fromGL); err != nil {
		tx.Rollback()
		return errors.New("sender account not found")
	}
	if err := tx.QueryRowContext(ctx, `SELECT gl_account_id FROM accounts WHERE id=$1`, input.ToAccountID).Scan(&toGL); err != nil {
		tx.Rollback()
		return errors.New("receiver account not found")
	}

	var journalID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO journal_entries (reference_type_id, reference_id, description)
		VALUES (1, $1, 'Transfer')
		RETURNING id
	`, input.FromAccountID).Scan(&journalID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO ledger_entries (gl_account_id, journal_entry_id, credit)
		VALUES ($1, $2, $3)
	`, fromGL, journalID, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO ledger_entries (gl_account_id, journal_entry_id, debit)
		VALUES ($1, $2, $3)
	`, toGL, journalID, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}