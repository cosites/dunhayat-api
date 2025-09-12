package usecase

import (
	"context"
	"fmt"

	"dunhayat-api/internal/products"
	"dunhayat-api/internal/products/repository"
)

type GetProductUseCase interface {
	Execute(ctx context.Context, productID string) (*products.Product, error)
}

type getProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewGetProductUseCase(
	productRepo repository.ProductRepository,
) GetProductUseCase {
	return &getProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *getProductUseCase) Execute(
	ctx context.Context,
	productID string,
) (*products.Product, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}
