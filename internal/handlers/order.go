package handlers

import (
	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"net/http"
)

type OrderHandler struct {
	OrderService order.OrderService
}

func ProvideOrderHandler(orderService order.OrderService) OrderHandler {
	return OrderHandler{OrderService: orderService}
}

func (h *OrderHandler) Router(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Post("/check-out", h.Checkout)
	})
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	checkout, err := h.OrderService.CheckoutCart(claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, checkout)
}
