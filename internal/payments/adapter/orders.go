package adapter

import (
	"context"

	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/usecase"

	"github.com/google/uuid"
)

type OrdersPaymentService struct {
	initiatePaymentUseCase usecase.InitiatePaymentUseCase
	verifyPaymentUseCase   usecase.VerifyPaymentUseCase
}

func NewOrdersPaymentService(
	initiatePaymentUseCase usecase.InitiatePaymentUseCase,
	verifyPaymentUseCase usecase.VerifyPaymentUseCase,
) port.PaymentService {
	return &OrdersPaymentService{
		initiatePaymentUseCase: initiatePaymentUseCase,
		verifyPaymentUseCase:   verifyPaymentUseCase,
	}
}

func (s *OrdersPaymentService) InitiatePayment(
	ctx context.Context,
	req *port.InitiatePaymentRequest,
) (*port.InitiatePaymentResponse, error) {
	paymentReq := &payments.InitiatePaymentRequest{
		OrderID:     req.OrderID,
		UserID:      req.UserID,
		Amount:      req.Amount,
		Method:      payments.PaymentMethodZibal, // TODO: Shall be configurable
		CallbackURL: req.CallbackURL,
		ReturnURL:   req.ReturnURL,
		Description: req.Description,
		Metadata:    req.Metadata,
	}

	paymentResp, err := s.initiatePaymentUseCase.Execute(ctx, paymentReq)
	if err != nil {
		return nil, err
	}

	return &port.InitiatePaymentResponse{
		PaymentID:    paymentResp.PaymentID,
		GatewayURL:   paymentResp.GatewayURL,
		GatewayRefID: paymentResp.GatewayRefID,
		Status:       paymentResp.Status.String(),
		Amount:       paymentResp.Amount,
		ExpiresAt:    paymentResp.ExpiresAt,
	}, nil
}

func (s *OrdersPaymentService) VerifyPayment(
	ctx context.Context,
	paymentID uuid.UUID,
) (*port.VerifyPaymentResponse, error) {
	paymentReq := &payments.VerifyPaymentRequest{
		PaymentID: paymentID,
	}

	paymentResp, err := s.verifyPaymentUseCase.Execute(ctx, paymentReq)
	if err != nil {
		return nil, err
	}

	return &port.VerifyPaymentResponse{
		PaymentID:    paymentResp.PaymentID,
		Status:       paymentResp.Status.String(),
		Amount:       paymentResp.Amount,
		GatewayRefID: paymentResp.GatewayRefID,
		PaidAt:       paymentResp.PaidAt,
		FailedAt:     paymentResp.FailedAt,
	}, nil
}
