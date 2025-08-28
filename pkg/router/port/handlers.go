package port

import (
	"net/http"
)

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

type ProductHandler interface {
	RouteRegistrar
	ListProducts(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
}

type OrderHandler interface {
	RouteRegistrar
	CreateOrder(w http.ResponseWriter, r *http.Request)
	GetOrder(w http.ResponseWriter, r *http.Request)
	CancelOrder(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	RouteRegistrar
	RequestOTP(w http.ResponseWriter, r *http.Request)
	VerifyOTP(w http.ResponseWriter, r *http.Request)
	GetOTPStatus(w http.ResponseWriter, r *http.Request)
}

type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}
