package postgres

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/models"
)

type PostgresProductStore struct {
	db *sql.DB
}

func NewPostgresProductStore(db *sql.DB) *PostgresProductStore {
	return &PostgresProductStore{db: db}
}

// - Methods -
func (store *PostgresProductStore) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, description, image_url, price, discount)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	err := store.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Price,
		product.Discount,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
