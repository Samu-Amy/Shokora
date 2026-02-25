package product

type Product struct {
	Id          int64   `json:"id"` // Generated
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`   // Defalut 0
	Version     int     `json:"version"`    // Defalut 0
	CreatedAt   string  `json:"created_at"` // Default now()
	UpdatedAt   string  `json:"updated_at"` // Default now()
}

// TODO: es. struct per prodotti menu/shop
type ShopProduct struct {
	Product
	// Price float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
	// Discount float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
}
