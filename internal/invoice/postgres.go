package invoice

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

var (
	ErrInvoiceNotFound = fmt.Errorf("invoice dost exist")
)

func NewPostgresRepo(p *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: p}
}

func (r *PostgresRepo) Create(ctx context.Context, invoice *Invoice) error {
	{

		query := `
	INSERT INTO invoices (
		id,
		user_id,
		status,
		amount,
		description,
		callback_url,
		pay_to_address,
		paid_by_address,
		overpayment,
		created_at,
		updated_at,       
		expired_at
	)
	VALUES (
		@id,
		@user_id,
		@status,
		@amount,
		@description,
		@callback_url,
		@pay_to_address,
		@paid_by_address,
		@overpayment,
		@created_at,
		@updated_at,       
		@expired_at
	)
	`

		args := pgx.NamedArgs{
			"id":              invoice.ID,
			"user_id":         invoice.UserID,
			"status":          invoice.Status,
			"amount":          invoice.Amount,
			"description":     invoice.Description,
			"callback_url":    invoice.CallbackURL,
			"pay_to_address":  invoice.PayToAddress,
			"paid_by_address": invoice.PaidByAddress,
			"overpayment":     invoice.Overpayment,
			"created_at":      invoice.CreatedAt,
			"updated_at":      invoice.UpdatedAt,
			"expired_at":      invoice.ExpiredAt,
		}

		_, err := r.pool.Exec(ctx, query, args)
		if err != nil {
			return err
		}

		return nil
	}

}

func (r *PostgresRepo) Update(ctx context.Context, invoice *Invoice) error {

	query := `
	UPDATE invoices
	SET
		amount = @amount,
		description = @description,
		updated_at = NOW()
	WHERE id = @id
	  AND user_id = @user_id
	  AND status = 'pending';
	`

	args := pgx.NamedArgs{
		"id":          invoice.ID,
		"user_id":     invoice.UserID,
		"amount":      invoice.Amount,
		"description": invoice.Description,
	}

	cmd, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrInvoiceNotFound
	}

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	var inv Invoice

	err := r.pool.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			status,
			amount,
			description,
			callback_url,
			pay_to_address,
			paid_by_address,
			overpayment,
			created_at,
			updated_at,
			expired_at
		FROM invoices
		WHERE id = $1
		  AND status <> 'deleted';
	`, id).Scan(
		&inv.ID,
		&inv.UserID,
		&inv.Status,
		&inv.Amount,
		&inv.Description,
		&inv.CallbackURL,
		&inv.PayToAddress,
		&inv.PaidByAddress,
		&inv.Overpayment,
		&inv.CreatedAt,
		&inv.UpdatedAt,
		&inv.ExpiredAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInvoiceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &inv, nil
}

func (r *PostgresRepo) ListByUser(ctx context.Context, userID uuid.UUID, p Pagination) (*[]Invoice, error) {

	offset := (p.Page - 1) * p.Limit

	query := `
	SELECT
		id,
		user_id,
		status,
		amount,
		description,
		callback_url,
		pay_to_address,
		paid_by_address,
		overpayment,
		created_at,
		updated_at,
		expired_at
	FROM invoices
	WHERE user_id = @user_id AND status <> 'deleted'
	ORDER BY created_at DESC
	LIMIT @limit
	OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"user_id": userID,
		"limit":   p.Limit,
		"offset":  offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]Invoice, 0, p.Limit)

	for rows.Next() {

		var invoice Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.UserID,
			&invoice.Status,
			&invoice.Amount,
			&invoice.Description,
			&invoice.CallbackURL,
			&invoice.PayToAddress,
			&invoice.PaidByAddress,
			&invoice.Overpayment,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&invoice.ExpiredAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &invoices, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, invoiceID, userID uuid.UUID) error {

	query := `
	UPDATE invoices
	SET
		status = @status,	
		updated_at = NOW()
	WHERE id = @id AND user_id = @user_id;
	`

	args := pgx.NamedArgs{
		"id":      invoiceID,
		"user_id": userID,
		"status":  StatusDeleted,
	}

	cmd, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrInvoiceNotFound
	}

	return nil
}

func (r *PostgresRepo) GetPending(ctx context.Context) (*[]Invoice, error) {

	query := `
	SELECT
		id,
		user_id,
		status,
		expired_at
	FROM invoices
	WHERE status = 'pending'
	ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]Invoice, 0)

	for rows.Next() {

		var invoice Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.UserID,
			&invoice.Status,
			&invoice.ExpiredAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &invoices, nil
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, invoice *Invoice) error {

	query := `
	UPDATE invoices
	SET
		status = @status,
		updated_at = NOW()
	WHERE id = @id;
	`

	args := pgx.NamedArgs{
		"status": invoice.Status,
		"id":     invoice.ID,
	}

	cmd, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrInvoiceNotFound
	}

	return nil
}
