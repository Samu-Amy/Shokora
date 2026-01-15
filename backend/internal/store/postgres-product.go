package store

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresProductStore struct {
	db *sql.DB
}

func NewPostgresProductStore(db *sql.DB) *PostgresProductStore {
	return &PostgresProductStore{db: db}
}

// ----- CREATE -----

func (store *PostgresProductStore) Create(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO products (name, description, image_url, price, discount)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

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

// ----- GET -----

// TODO: crea diversi metodi per ottenere dati (uno per gli utenti (da menu/shop/sito) quando guardano i dettagli di un prodotto e uno per gli admin quando vogliono vedere i dati (completi, con anche version, created_at, updated_at, ecc.) di un prodotto)
func (store *PostgresProductStore) GetById(ctx context.Context, productId int64) (*Product, error) {
	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var product Product

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

func (store *PostgresProductStore) GetProducts(ctx context.Context, queryPaginationOptions QueryPaginationOptions, productsFilters ProductsFilters) ([]Product, error) {
	// TODO: sistema implementazione (come GetMenuProducts -> guarda i TODO lì)

	// For added safety I don't use the sort parameter directly (even if there's validation)
	sort := "ASC"
	if queryPaginationOptions.Sort == "desc" {
		sort = "DESC"
	}

	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
		ORDER BY name ` + sort + `
		LIMIT $1 OFFSET $2
	`

	rows, err := store.db.QueryContext(
		ctx,
		query,
		queryPaginationOptions.Limit,
		queryPaginationOptions.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product

		err := rows.Scan(
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
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// ----- UPDATE -----

func (store *PostgresProductStore) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, image_url = $3, price = $4, discount = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

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

// ----- DELETE -----

func (store *PostgresProductStore) Delete(ctx context.Context, productId int64) error {
	query := `DELETE FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

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
