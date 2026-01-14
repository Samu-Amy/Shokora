package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (store *PostgresProductStore) GetById(ctx context.Context, productId int64) (*models.Product, error) {
	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product

	err := store.db.QueryRowContext(
		ctx,
		query,
		productId,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
		&product.Discount,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (store *PostgresProductStore) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, image_url = $3, price = $4, discount = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
	`

	// TODO: modifica e rendi "modulare" per aggiornare solo ciò che serve (e magari fai un controllo più accurato sulle modifiche fatte e non solo sulla versione)

	err := store.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Price,
		product.Discount,
		product.ID,
		product.Version,
	).Scan(
		&product.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// TODO: fai check di esistenza del prodotto per distinguere tra not found e conflict (version)
			// sia in caso di prodotto (id) non trovato, sia in caso di versione vecchia

			// ...check if product exists...
			// return ErrVersionConlflict

			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (store *PostgresProductStore) Delete(ctx context.Context, productId int64) error {
	query := `DELETE FROM products WHERE id = $1`

	res, err := store.db.ExecContext(ctx, query, productId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Nothing deleted
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// TODO: guarda immagine cheatsheet sql
