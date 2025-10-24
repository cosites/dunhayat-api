package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dunhayat-api/internal/orders"
	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/orders/repository"

	"github.com/google/uuid"
)

type CreateOrderUseCase interface {
	Execute(
		ctx context.Context,
		userID uuid.UUID,
		req *orders.CreateOrderRequest,
	) (*orders.OrderResponse, error)
}

type createOrderUseCase struct {
	saleRepo            repository.SaleRepository
	saleItemRepo        repository.SaleItemRepository
	cartReservationRepo repository.CartReservationRepository
	productPort         port.ProductPort
	userPort            port.UserPort
	paymentPort         port.PaymentPort
}

func NewCreateOrderUseCase(
	saleRepo repository.SaleRepository,
	saleItemRepo repository.SaleItemRepository,
	cartReservationRepo repository.CartReservationRepository,
	productPort port.ProductPort,
	userPort port.UserPort,
	paymentPort port.PaymentPort,
) CreateOrderUseCase {
	return &createOrderUseCase{
		saleRepo:            saleRepo,
		saleItemRepo:        saleItemRepo,
		cartReservationRepo: cartReservationRepo,
		productPort:         productPort,
		userPort:            userPort,
		paymentPort:         paymentPort,
	}
}

func (uc *createOrderUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	req *orders.CreateOrderRequest,
) (*orders.OrderResponse, error) {
	user, err := uc.userPort.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	var totalPrice int
	var saleItems []orders.SaleItem

	for _, item := range req.Items {
		product, err := uc.productPort.GetProductByID(
			ctx, item.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get product %s: %w",
				item.ProductID, err,
			)
		}
		if product == nil {
			return nil, fmt.Errorf(
				"product not found: %s",
				item.ProductID,
			)
		}
		if product.InStock < item.Quantity {
			return nil, fmt.Errorf(
				"insufficient stock for product %s: requested %d, available %d",
				item.ProductID, item.Quantity, product.InStock,
			)
		}

		itemPrice := product.Price * item.Quantity
		totalPrice += itemPrice

		reservation := &orders.CartReservation{
			UserID:    userID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}
		if err := uc.cartReservationRepo.Create(
			ctx, reservation,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to create cart reservation for product %s: %w",
				item.ProductID, err,
			)
		}

		if err := uc.productPort.UpdateStock(
			ctx, item.ProductID, -item.Quantity,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to update product stock for %s: %w",
				item.ProductID, err,
			)
		}
	}

	sale := &orders.Sale{
		UserID:     userID,
		Status:     orders.OrderStatusPending,
		TotalPrice: totalPrice,
	}
	if err := uc.saleRepo.Create(ctx, sale); err != nil {
		return nil, fmt.Errorf("failed to create sale: %w", err)
	}

	for _, item := range req.Items {
		product, err := uc.productPort.GetProductByID(
			ctx, item.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get product for sale item %s: %w",
				item.ProductID, err,
			)
		}
		saleItem := &orders.SaleItem{
			SaleID:    sale.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		if err := uc.saleItemRepo.Create(ctx, saleItem); err != nil {
			return nil, fmt.Errorf(
				"failed to create sale item for product %s: %w",
				item.ProductID, err,
			)
		}
		saleItems = append(saleItems, *saleItem)
	}

	paymentReq := &port.InitiatePaymentRequest{
		OrderID:     sale.ID,
		UserID:      userID,
		Amount:      totalPrice,
		CallbackURL: req.CallbackURL,
		ReturnURL:   req.ReturnURL,
		Description: fmt.Sprintf(
			"Payment for order %s", sale.ID.String(),
		),
		Metadata: map[string]any{
			"order_id":    sale.ID.String(),
			"address":     req.Address,
			"postal_code": req.PostalCode,
		},
	}

	paymentResp, err := uc.paymentPort.InitiatePayment(ctx, paymentReq)
	if err != nil {
		if err := uc.saleRepo.UpdateStatus(
			ctx, sale.ID, orders.OrderStatusCancelled,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to update sale status: %w",
				err,
			)
		}
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	return &orders.OrderResponse{
		ID:           sale.ID,
		UserID:       userID,
		Status:       sale.Status,
		TrackingCode: sale.TrackingCode,
		TotalPrice:   sale.TotalPrice,
		Items:        saleItems,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Payment: &orders.PaymentInfo{
			PaymentID:    paymentResp.PaymentID,
			GatewayURL:   paymentResp.GatewayURL,
			GatewayRefID: paymentResp.GatewayRefID,
			Status:       paymentResp.Status,
			Amount:       paymentResp.Amount,
			ExpiresAt:    paymentResp.ExpiresAt,
		},
		CreatedAt: sale.CreatedAt,
		UpdatedAt: sale.UpdatedAt,
	}, nil
}
