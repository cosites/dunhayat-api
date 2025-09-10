package repository

import (
	"context"
	"errors"
	"time"

	"dunhayat-api/internal/orders"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleRepository interface {
	Create(ctx context.Context, sale *orders.Sale) error
	GetByID(ctx context.Context, id uuid.UUID) (*orders.Sale, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]orders.Sale, error)
	Update(ctx context.Context, sale *orders.Sale) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SaleItemRepository interface {
	Create(ctx context.Context, item *orders.SaleItem) error
	GetBySaleID(ctx context.Context, saleID uuid.UUID) ([]orders.SaleItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CartReservationRepository interface {
	Create(ctx context.Context, reservation *orders.CartReservation) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]orders.CartReservation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CleanExpired(ctx context.Context) error
}

type postgresSaleRepository struct {
	db *gorm.DB
}

type postgresSaleItemRepository struct {
	db *gorm.DB
}

type postgresCartReservationRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &postgresSaleRepository{db: db}
}

func NewSaleItemRepository(db *gorm.DB) SaleItemRepository {
	return &postgresSaleItemRepository{db: db}
}

func NewCartReservationRepository(db *gorm.DB) CartReservationRepository {
	return &postgresCartReservationRepository{db: db}
}

func (r *postgresSaleRepository) Create(
	ctx context.Context,
	sale *orders.Sale,
) error {
	return r.db.WithContext(ctx).Create(sale).Error
}

func (r *postgresSaleRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*orders.Sale, error) {
	var sale orders.Sale
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sale, nil
}

func (r *postgresSaleRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]orders.Sale, error) {
	var sales []orders.Sale
	err := r.db.WithContext(ctx).Where(
		"user_id = ?",
		userID,
	).Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *postgresSaleRepository) Update(
	ctx context.Context,
	sale *orders.Sale,
) error {
	return r.db.WithContext(ctx).Save(sale).Error
}

func (r *postgresSaleRepository) Delete(
	ctx context.Context, id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&orders.Sale{}, id).Error
}

func (r *postgresSaleItemRepository) Create(
	ctx context.Context, item *orders.SaleItem,
) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *postgresSaleItemRepository) GetBySaleID(
	ctx context.Context, saleID uuid.UUID,
) ([]orders.SaleItem, error) {
	var items []orders.SaleItem
	err := r.db.WithContext(ctx).Where(
		"sale_id = ?", saleID,
	).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresSaleItemRepository) Delete(
	ctx context.Context, id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&orders.SaleItem{}, id).Error
}

func (r *postgresCartReservationRepository) Create(
	ctx context.Context, reservation *orders.CartReservation,
) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

func (r *postgresCartReservationRepository) GetByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]orders.CartReservation, error) {
	var reservations []orders.CartReservation
	err := r.db.WithContext(ctx).Where(
		"user_id = ?", userID,
	).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *postgresCartReservationRepository) Delete(
	ctx context.Context, id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(
		&orders.CartReservation{}, id,
	).Error
}

func (r *postgresCartReservationRepository) CleanExpired(
	ctx context.Context,
) error {
	return r.db.WithContext(ctx).Where(
		"expires_at < ?", time.Now(),
	).Delete(&orders.CartReservation{}).Error
}
