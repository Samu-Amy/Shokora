package store

import (
	"context"
)

// ----- GET -----

func (store *PostgresProductStore) GetMenuProducts(ctx context.Context, queryPaginationOptions QueryPaginationOptions, menuFilters MenuFilters) ([]Product, error) {
	// TODO: (forse si possono togliere i dati che non servono (version, update_at, ecc.)
	// TODO: aggiungi dati con JOIN (badges, ingredients, ecc.) -> fai
	// TODO: dividere/raggruppare per categoria (es. "menuSectionId"?) e come ordinarli (?)

	// TODO: aggiungi opzioni per scelta su quale parametro usare per sorting

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

	queryCtx, cancel := context.WithTimeout(ctx, long_query_timeout) //TODO: va bene?
	defer cancel()

	rows, err := store.db.QueryContext(
		queryCtx,
		query,
		menuFilters.Search,
		queryPaginationOptions.Limit,
		queryPaginationOptions.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var menu_products []Product

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
			return nil, err
		}

		menu_products = append(menu_products, product)
	}

	return menu_products, nil
}
