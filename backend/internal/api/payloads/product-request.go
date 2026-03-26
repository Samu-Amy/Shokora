package payloads

type CreateProductReq struct {
	Name        string  `json:"name" validate:"required,min=1,max=150,valid-chars"` // Required
	Description string  `json:"description" validate:"max=2500,valid-chars"`        // Default ""
	ImageURL    string  `json:"image_url" validate:"omitempty,url,valid-chars"`     // Default "" // TODO: aggiungi (anche in Update struct) "url" al validate se l'url per accedere alle foto soddisfa la validazione del validator
	Price       float64 `json:"price" validate:"required,gt=0"`                     // Required
	Discount    float64 `json:"discount" validate:"gte=0,lte=1"`                    // Default 0
}

type UpdateProductReq struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=150,valid-chars"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=2500,valid-chars"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty,url,valid-chars"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Discount    *float64 `json:"discount,omitempty" validate:"omitempty,gte=0,lte=1"`
}
