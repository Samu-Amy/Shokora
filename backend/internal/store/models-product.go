package store

import (
	"context"
)

type Product struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	// Ingredients []string `json:ingredients` //? array di string, enums o id di prodotti
	// Badges    []string `json:badges` // TODO: string di enum invece che string (oppure si usano enums ma si salvano come string nel db) (?)
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount"`
	Version   int     `json:"version"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// TODO: es. struct per prodotti menu/shop
type ShopProduct struct {
	Product
	// Price float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
	// Discount float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
}

// TODO: aggiungi nomi parametri/argomenti nei metodi di tutte le interfacce
type ProductRepository interface {
	Create(context.Context, *Product) error
	GetById(context.Context, int64) (*Product, error)
	GetProducts(context.Context, QueryPaginationOptions, ProductsFilters) ([]Product, error) // TODO: (anche per GetMenuProducts e GetShopProducts) modifica struct ritornata (aggiungendo anche badges ed altro)
	GetMenuProducts(context.Context, QueryPaginationOptions, MenuFilters) ([]Product, error)
	Update(context.Context, *Product) error
	Delete(context.Context, int64) error
}
