package api

import (
	"context"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store/models"
)

var productIdParam = "productId"

// ----- CREATE -----

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required,min=1,max=150"` // Required
	Description string  `json:"description" validate:"max=2500"`        // Default ""
	ImageURL    string  `json:"image_url" validate:"omitempty"`         // Default "" // TODO: aggiungi (anche in Update struct) "url" al validate se l'url per accedere alle foto soddisfa la validazione del validator
	Price       float64 `json:"price" validate:"required,gt=0"`         // Required
	Discount    float64 `json:"discount" validate:"gte=0,lte=1"`        // Default 0
}

func (app *App) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	// Create product from payload data
	product := &models.Product{
		Name:        payload.Name,
		Description: payload.Description,
		ImageURL:    payload.ImageURL,
		Price:       payload.Price,
		Discount:    payload.Discount,
	}

	// Create product
	if err := app.store.Product.Create(ctx, product); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return product
	if err := app.jsonResponse(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- GET -----

func (app *App) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Get product
	product, err := app.getProduct(ctx, productId)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return product
	if err := app.jsonResponse(w, http.StatusOK, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- UPDATE -----

type UpdateProductPayload struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=150"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=2500"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Discount    *float64 `json:"discount,omitempty" validate:"omitempty,gte=0,lte=1"`
}

func (app *App) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Get product
	product, err := app.getProduct(ctx, productId)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

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

	// Update fields in product (only the ones specified in the payload)
	if payload.Name != nil {
		product.Name = *payload.Name
	}

	if payload.Description != nil {
		product.Description = *payload.Description
	}

	if payload.ImageURL != nil {
		product.ImageURL = *payload.ImageURL
	}

	if payload.Price != nil {
		product.Price = *payload.Price
	}

	if payload.Discount != nil {
		product.Discount = *payload.Discount
	}

	// Update and return
	if err := app.store.Product.Update(r.Context(), product); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return product
	if err := app.jsonResponse(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- DELETE -----

func (app *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Product.Delete(ctx, productId); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content (product deleted)
	w.WriteHeader(http.StatusNoContent)
}

// ----- UTILS -----

func (app *App) getProduct(ctx context.Context, productId int64) (*models.Product, error) {
	product, err := app.store.Product.GetById(ctx, productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
