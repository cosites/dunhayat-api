package port

import "context"

type Product struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	InStock int    `json:"in_stock"`
}

type ProductService interface {
	GetProductByID(ctx context.Context, productID string) (*Product, error)
	UpdateStock(ctx context.Context, productID string, quantity int) error
}
