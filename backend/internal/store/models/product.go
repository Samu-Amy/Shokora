package models

import "context"

type ProductRepository interface {
	Create(context.Context) error
}
