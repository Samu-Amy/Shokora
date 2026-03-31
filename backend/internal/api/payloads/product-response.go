package payloads

// TODO: sistema (questa è solo di test)
type ProductRes struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
}
