package api

import (
	"context"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payload"
	"github.com/Samu-Amy/Shokora/internal/store"
)

// ----- CREATE -----

func (app *App) createProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payload.CreateProductPayload

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
	product := &store.Product{
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

func (app *App) getProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Get product
	product, err := app.getProductById(ctx, productId)
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

func (app *App) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Get product
	product, err := app.getProductById(ctx, productId)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	// Get payload data
	var payload payload.UpdateProductPayload
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

func (app *App) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get product id
	productId, err := app.getIdFromParam(r, productIdParam)
	if err != nil {
		app.badRequestError(w, r, err)
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

func (app *App) getProductById(ctx context.Context, productId int64) (*store.Product, error) {
	product, err := app.store.Product.GetById(ctx, productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
