package handlers

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/Samu-Amy/Shokora/internal/store/models"
)

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

		// Create product from payload data
		product := &models.Product{
			Name:        payload.Name,
			Description: payload.Description,
			ImageURL:    payload.ImageURL,
			Price:       payload.Price,
			Discount:    payload.Discount,
		}
		// TODO: i dati non definiti che valori hanno (e come faccio valori opzionali (es. discount))?

		ctx := r.Context()

		if err := store.Product.Create(ctx, product); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Send product data
		if err := writeJSON(w, http.StatusCreated, product); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
