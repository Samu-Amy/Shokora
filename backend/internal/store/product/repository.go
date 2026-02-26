package product

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type IProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetById(ctx context.Context, productId int64) (*Product, error)
	GetProducts(ctx context.Context, queryPaginationOptions database.QueryPaginationOptions, productsFilters database.ProductsFilters) ([]Product, error) // TODO: adatta per funzionare sia per Menu e Shop, modifica struct ritornata (aggiungendo anche badges ed altro)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, productId int64) error
}
