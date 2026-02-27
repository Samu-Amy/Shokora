package product

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresProductStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresProductStore {
	return &PostgresProductStore{db: db}
}

// ----- CREATE -----

func (store *PostgresProductStore) Create(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO products (name, description, image_url, price, discount)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Price,
		product.Discount,
	).Scan(
		&product.Id,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return database.ParseDbError(err)
}

// ----- GET -----

// TODO: crea diversi metodi per ottenere dati (uno per gli utenti (da menu/shop/sito) quando guardano i dettagli di un prodotto e uno per gli admin quando vogliono vedere i dati (completi, con anche version, created_at, updated_at, ecc.) di un prodotto)
func (store *PostgresProductStore) GetById(ctx context.Context, productId int64) (*Product, error) {
	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var product Product

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		productId,
	).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
		&product.Discount,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return &product, database.ParseDbError(err)
}

func (store *PostgresProductStore) GetProducts(ctx context.Context, queryPaginationOptions database.QueryPaginationOptions, productsFilters database.ProductsFilters) ([]Product, error) {
	// TODO: sistema implementazione (come GetMenuProducts -> guarda i TODO lì)

	// For added safety I don't use the sort parameter directly (even if there's validation)
	sort := "ASC"
	if queryPaginationOptions.Sort == "desc" {
		sort = "DESC"
	}

	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
		WHERE (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
		ORDER BY name ` + sort + `
		LIMIT $2 OFFSET $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.LongQueryTimeout) //TODO: va bene?
	defer cancel()

	rows, err := store.db.QueryContext(
		queryCtx,
		query,
		productsFilters.Search,
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
			&product.Id,
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
			return nil, database.ParseDbError(err)
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

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	// TODO: modifica e rendi "modulare" per aggiornare solo ciò che serve (e magari fai un controllo più accurato sulle modifiche fatte e non solo sulla versione)

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Price,
		product.Discount,
		product.Id,
		product.Version,
	).Scan(
		&product.Version,
	)

	// TODO: fai check di esistenza del prodotto per distinguere tra not found e conflict (version)
	// sia in caso di prodotto (id) non trovato, sia in caso di versione vecchia

	// ...check if product exists...
	// return ErrVersionConlflict

	return database.ParseDbError(err)
}

// ----- DELETE -----

func (store *PostgresProductStore) Delete(ctx context.Context, productId int64) error {
	query := `DELETE FROM products WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, productId))
}

// TODO: guarda immagine cheatsheet sql
