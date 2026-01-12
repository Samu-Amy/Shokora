package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/Samu-Amy/Shokora/internal/store/models"
	"github.com/Samu-Amy/Shokora/internal/store/postgres"
	"github.com/go-chi/chi/v5"
)

// ----- CREATE -----

type CreateProductPayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
}

func CreateProduct(store *store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get payload data
		var payload CreateProductPayload

		if err := readJSON(w, r, &payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: validaizone (e setta valori di default (?))

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
		if err := store.Product.Create(ctx, product); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Send product data to frontend
		if err := writeJSON(w, http.StatusCreated, product); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func GetProduct(store *store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get param and convert it
		idParam := chi.URLParam(r, "productId")

		postId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := r.Context()

		// Get product
		product, err := store.Product.GetById(ctx, postId)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				writeJSONError(w, http.StatusNotFound, err.Error())
			default:
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		// Send product data to frontend
		if err := writeJSON(w, http.StatusCreated, product); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// TODO: eventualmente sposta in store
type UpdateProductPayload struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Discount    *float64 `json:"discount,omitempty"`
}
