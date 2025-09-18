package adapter

import (
	"context"

	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/products/repository"
)

type OrdersProductAdapter struct {
	productRepo repository.ProductRepository
}

func NewOrdersProductAdapter(
	productRepo repository.ProductRepository,
) port.ProductPort {
	return &OrdersProductAdapter{
		productRepo: productRepo,
	}
}

func (s *OrdersProductAdapter) GetProductByID(
	ctx context.Context,
	productID string,
) (*port.Product, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	return &port.Product{
		ID:      product.ID,
		Name:    product.Name,
		Price:   product.Price,
		InStock: product.InStock,
	}, nil
}

func (s *OrdersProductAdapter) UpdateStock(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	return s.productRepo.UpdateStock(ctx, productID, quantity)
}
