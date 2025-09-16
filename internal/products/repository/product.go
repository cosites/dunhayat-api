package repository

import (
	"context"
	"errors"

	"dunhayat-api/internal/products"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *products.Product) error
	GetByID(ctx context.Context, id string) (*products.Product, error)
	GetAll(ctx context.Context) ([]products.Product, error)
	GetByCategory(ctx context.Context, category products.Category) ([]products.Product, error)
	Update(ctx context.Context, product *products.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, quantity int) error
}

type postgresProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(
	ctx context.Context,
	product *products.Product,
) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *postgresProductRepository) GetByID(
	ctx context.Context,
	id string,
) (*products.Product, error) {
	var product products.Product
	err := r.db.WithContext(ctx).Where(
		"id = ?", id,
	).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *postgresProductRepository) GetAll(
	ctx context.Context,
) ([]products.Product, error) {
	var productList []products.Product
	err := r.db.WithContext(ctx).Find(&productList).Error
	if err != nil {
		return nil, err
	}
	return productList, nil
}

func (r *postgresProductRepository) GetByCategory(
	ctx context.Context,
	category products.Category,
) ([]products.Product, error) {
	var productList []products.Product
	err := r.db.WithContext(ctx).Where(
		"category = ?", category,
	).Find(&productList).Error
	if err != nil {
		return nil, err
	}
	return productList, nil
}

func (r *postgresProductRepository) Update(
	ctx context.Context,
	product *products.Product,
) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *postgresProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	return r.db.WithContext(ctx).Delete(&products.Product{}, id).Error
}

func (r *postgresProductRepository) UpdateStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	return r.db.WithContext(ctx).Model(&products.Product{}).
		Where("id = ?", id).Update(
			"in_stock",
			gorm.Expr(
				"in_stock + ?", quantity,
			),
		).Error
}
