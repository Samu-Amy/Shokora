package store

import (
	"context"
)

// ----- GET -----

func (store *PostgresProductStore) GetMenuProducts(ctx context.Context, queryPaginationOptions QueryPaginationOptions) ([]Product, error) {
	// TODO: (forse si possono togliere i dati che non servono (version, update_at, ecc.)
	// TODO: aggiungi dati con JOIN (badges, ingredients, ecc.) -> fai
	// TODO: dividere/raggruppare per categoria (es. "menuSectionId"?) e come ordinarli (?)
	query := `
		SELECT id, name, description, image_url, price, discount, version, created_at, updated_at
		FROM products
	`

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var menu_products []Product

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

		menu_products = append(menu_products, product)
	}

	return menu_products, nil
}
