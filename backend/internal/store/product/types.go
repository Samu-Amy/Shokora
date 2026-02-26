package product

// TODO: es. struct per prodotti menu/shop (solo lato store-service, poi usa payloads)
type ShopProduct struct {
	Product
	// Price float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
	// Discount float64 (magari tolto da Product normale e messo nelle "versioni" menu/shop ?)
}
