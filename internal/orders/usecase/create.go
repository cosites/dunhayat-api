package usecase

import (
	"context"
	"errors"
	"time"

	"dunhayat-api/internal/orders"
	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/orders/repository"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, req *orders.CreateOrderRequest) (*orders.OrderResponse, error)
}

type createOrderUseCase struct {
	saleRepo            repository.SaleRepository
	saleItemRepo        repository.SaleItemRepository
	cartReservationRepo repository.CartReservationRepository
	productService      port.ProductService
	userService         port.UserService
}

func NewCreateOrderUseCase(
	saleRepo repository.SaleRepository,
	saleItemRepo repository.SaleItemRepository,
	cartReservationRepo repository.CartReservationRepository,
	productService port.ProductService,
	userService port.UserService,
) CreateOrderUseCase {
	return &createOrderUseCase{
		saleRepo:            saleRepo,
		saleItemRepo:        saleItemRepo,
		cartReservationRepo: cartReservationRepo,
		productService:      productService,
		userService:         userService,
	}
}

func (uc *createOrderUseCase) Execute(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (*orders.OrderResponse, error) {
	user, err := uc.userService.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	var totalPrice int
	var saleItems []orders.SaleItem

	for _, item := range req.Items {
		product, err := uc.productService.GetProductByID(
			ctx, item.ProductID,
		)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New(
				"product not found: " + item.ProductID,
			)
		}
		if product.InStock < item.Quantity {
			return nil, errors.New(
				"insufficient stock for product: " + item.ProductID,
			)
		}

		itemPrice := product.Price * item.Quantity
		totalPrice += itemPrice

		reservation := &orders.CartReservation{
			UserID:    req.UserID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}
		if err := uc.cartReservationRepo.Create(
			ctx, reservation,
		); err != nil {
			return nil, err
		}

		if err := uc.productService.UpdateStock(
			ctx, item.ProductID, -item.Quantity,
		); err != nil {
			return nil, err
		}
	}

	sale := &orders.Sale{
		UserID:     req.UserID,
		Status:     orders.OrderStatusPending,
		TotalPrice: totalPrice,
	}
	if err := uc.saleRepo.Create(ctx, sale); err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		product, _ := uc.productService.GetProductByID(
			ctx, item.ProductID,
		)
		saleItem := &orders.SaleItem{
			SaleID:    sale.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		if err := uc.saleItemRepo.Create(ctx, saleItem); err != nil {
			return nil, err
		}
		saleItems = append(saleItems, *saleItem)
	}

	return &orders.OrderResponse{
		ID:           sale.ID,
		UserID:       sale.UserID,
		Status:       sale.Status,
		TrackingCode: sale.TrackingCode,
		TotalPrice:   sale.TotalPrice,
		Items:        saleItems,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		CreatedAt:    sale.CreatedAt,
		UpdatedAt:    sale.UpdatedAt,
	}, nil
}
