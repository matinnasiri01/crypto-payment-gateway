package invoice

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type invoicePostgresRepo struct {
	pool *pgxpool.Pool
}

// todo if is impossible use one func to call to db also for users

var (
	ErrInvoiceNotFound = fmt.Errorf("invoice dost exist")
)

func NewPostgresRepo(p *pgxpool.Pool) *invoicePostgresRepo {
	return &invoicePostgresRepo{pool: p}
}

func (r *invoicePostgresRepo) Create(ctx context.Context, invoice *Invoice) error {
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

func (r *invoicePostgresRepo) Update(ctx context.Context, invoice *Invoice, userID uuid.UUID) error {

	query := `
	UPDATE invoices
	SET
		status = @status,
		amount = @amount,
		description = @description,
		callback_url = @callback_url,
		pay_to_address = @pay_to_address,
		paid_by_address = @paid_by_address,
		overpayment = @overpayment,		
		updated_at = NOW()
	WHERE id = @id AND user_id = @user_id;
	`

	args := pgx.NamedArgs{
		"id":              invoice.ID,
		"user_id":         userID,
		"status":          invoice.Status,
		"description":     invoice.Description,
		"callback_url":    invoice.CallbackURL,
		"pay_to_address":  invoice.PayToAddress,
		"paid_by_address": invoice.PaidByAddress,
		"overpayment":     invoice.Overpayment,
	}

	_, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
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

func (r *invoicePostgresRepo) ListByUser(ctx context.Context, userID uuid.UUID, p Pagination) (*[]Invoice, error) {

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

func (r *invoicePostgresRepo) Delete(ctx context.Context, invoiceID, userID uuid.UUID) error {

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
