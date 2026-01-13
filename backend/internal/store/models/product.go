package models

import "context"

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	// Ingredients []string `json:ingredients` //? array di string, enums o id di prodotti
	// Badges    []string `json:badges` // TODO: string di enum invece che string (oppure si usano enums ma si salvano come string nel db) (?)
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type ProductRepository interface {
	Create(context.Context, *Product) error
	GetById(context.Context, int64) (*Product, error)
	Update(context.Context, int64) (*Product, error)
	Delete(context.Context, int64) error
}
