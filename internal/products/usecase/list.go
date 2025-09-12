package usecase

import (
	"context"

	"dunhayat-api/internal/products"
	"dunhayat-api/internal/products/repository"
)

type ListProductsUseCase interface {
	Execute(ctx context.Context, category *products.Category) ([]products.Product, error)
}

type listProductsUseCase struct {
	productRepo repository.ProductRepository
}

func NewListProductsUseCase(
	productRepo repository.ProductRepository,
) ListProductsUseCase {
	return &listProductsUseCase{
		productRepo: productRepo,
	}
}

func (uc *listProductsUseCase) Execute(
	ctx context.Context, category *products.Category,
) ([]products.Product, error) {
	if category != nil {
		return uc.productRepo.GetByCategory(ctx, *category)
	}
	return uc.productRepo.GetAll(ctx)
}
