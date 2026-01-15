package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

func (app *App) getMenuProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// - Pagination, filters and sorting -

	// Define default values
	queryPaginationOptions := store.QueryPaginationOptions{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	// Parse values
	queryPaginationOptions, err := queryPaginationOptions.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(queryPaginationOptions); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// - Filters -

	// Define default values
	menuFilters := store.MenuFilters{
		Search: "",
		Badges: make([]string, 0),
	}

	// Parse values
	menuFilters, err = menuFilters.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(menuFilters); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// - Query -

	// Get menu products
	products, err := app.store.Product.GetMenuProducts(ctx, queryPaginationOptions, menuFilters)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return product
	if err := app.jsonResponse(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
