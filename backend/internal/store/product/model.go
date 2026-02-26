package product

type Product struct {
	Id          int64
	Name        string
	Description string
	ImageURL    string
	Price       float64
	Discount    float64
	Version     int
	CreatedAt   string
	UpdatedAt   string
}

// TODO: es. struct per prodotti menu/shop (solo lato store-service, poi usa payloads)
type ShopProduct struct {
	Product
	// Price float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
	// Discount float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
}
