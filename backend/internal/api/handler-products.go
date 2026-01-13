package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store/models"
	"github.com/Samu-Amy/Shokora/internal/store/postgres"
	"github.com/go-chi/chi/v5"
)

// ----- CREATE -----

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required,max=150"`
	Description string  `json:"description" validate:"required,max=2500"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price" validate:"gt=0"`
	Discount    float64 `json:"discount" validate:"gte=0,lte=1"`
}

func (app *App) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get payload data
	var payload CreateProductPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// TODO: setta valori di default e gestisci price (required) e discount (not required) (?)

	// Create product from payload data
	product := &models.Product{
		Name:        payload.Name,
		Description: payload.Description,
		ImageURL:    payload.ImageURL,
		Price:       payload.Price,
		Discount:    payload.Discount,
	}

	// if payload.Discount != nil {
	// 	product.Discount = payload.Discount
	// } else {
	// 	product.Discount = 0
	// }

	ctx := r.Context()

	// Create the product on db (and update product with missing data (id, created_at, updated_at) from db)
	if err := app.store.Product.Create(ctx, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send product data to frontend
	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// ----- GET -----

func (app *App) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get param and convert it
	idParam := chi.URLParam(r, "productId")

	productId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// Get product
	product, err := app.store.Product.GetById(ctx, productId)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Send product data to frontend
	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// ----- UPDATE -----

type UpdateProductPayload struct {
	ID          int64    `json:"id" validate:"required"`
	Name        *string  `json:"name,omitempty" validate:"max=150"`
	Description *string  `json:"description,omitempty" validate:"max=2500"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"gt=0"`
	Discount    *float64 `json:"discount,omitempty" validate:"gte=0,lte=1"`
}

func (app *App) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get payload data
	var payload UpdateProductPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// TODO: setta valori di default (?)

	// Create product from payload data
	// product := &models.Product{
	// 	Name:        payload.Name,
	// 	Description: payload.Description,
	// 	ImageURL:    payload.ImageURL,
	// 	Price:       payload.Price,
	// 	Discount:    payload.Discount,
	// }

	// ctx := r.Context()

	// // Create the product on db (and update product with missing data (id, created_at, updated_at) from db)
	// if err := app.store.Product.Create(ctx, product); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	// // Send product data to frontend
	// if err := writeJSON(w, http.StatusCreated, product); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

}

// ----- DELETE -----

func (app *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "productId")

	productId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Product.Delete(ctx, productId); err != nil {
		switch {
		case errors.Is(err, postgres.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
