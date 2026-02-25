package db

import (
	"net/http"
	"strings"
	"time"
)

// - Filters -

type ProductsFilters struct {
	Search string   `json:"search" validate:"max=100"`
	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (productsFilters ProductsFilters) Parse(r *http.Request) (ProductsFilters, error) {
	queryString := r.URL.Query()

	// Search
	search := queryString.Get("search")
	if search != "" {
		productsFilters.Search = search
	}

	// Badges
	badges := queryString.Get("badges")
	if badges != "" {
		productsFilters.Badges = strings.Split(badges, ",")
	}

	// Since
	since := queryString.Get("since")
	if since != "" {
		productsFilters.Since = parseTime(since)
	}

	// Until
	until := queryString.Get("until")
	if until != "" {
		productsFilters.Until = parseTime(until)
	}

	return productsFilters, nil
}

func parseTime(strTime string) string {
	t, err := time.Parse(time.DateTime, strTime)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}

type MenuFilters struct {
	Search string   `json:"search" validate:"max=100"`
	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
}

func (menuFilters MenuFilters) Parse(r *http.Request) (MenuFilters, error) {
	queryString := r.URL.Query()

	// Search
	search := queryString.Get("search")
	if search != "" {
		menuFilters.Search = search
	}

	// Badges
	badges := queryString.Get("badges")
	if badges != "" {
		menuFilters.Badges = strings.Split(badges, ",")
	}

	return menuFilters, nil
}

type ShopFilters struct {
	Search string   `json:"search" validate:"max=100"`
	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
}

func (shopFilters ShopFilters) Parse(r *http.Request) (ShopFilters, error) {
	queryString := r.URL.Query()

	// Search
	search := queryString.Get("search")
	if search != "" {
		shopFilters.Search = search
	}

	// Badges
	badges := queryString.Get("badges")
	if badges != "" {
		shopFilters.Badges = strings.Split(badges, ",")
	}

	return shopFilters, nil
}
