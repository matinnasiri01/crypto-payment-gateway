package invoice

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type invoicePostgresRepo struct {
	pool *pgxpool.Pool
}

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
func (r *invoicePostgresRepo) Update(ctx context.Context, invoice *Invoice) error {
	return nil
}
func (r *invoicePostgresRepo) GetByID(ctx context.Context, id string) (*Invoice, error) {
	return nil, nil
}
