package handlers

import (
	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
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
		r.Get("/", h.GetAllOrder)
	})
}

// GetAllOrder resolves all Orders
// @Summary Retrieve list of orders
// @Description This endpoint retrieves a list of orders with optional filters.
// @Tags order/order
// @Security JWTAuthentication
// @Param limit query int false "The number of products per page."
// @Param page query int false "The page number."
// @Produce json
// @Success 200 {object} response.Base{data=order.OrderResponse}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/order [get]
func (h *OrderHandler) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Missing Param Query Limit")
		response.WithError(w, err)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Missing Param Query Page")
		response.WithError(w, err)
		return
	}

	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusBadRequest, "Unauthorize")
		return
	}

	orders, err := h.OrderService.ResolveAllCart(claims.ID, limit, page-1)
	if err != nil {
		response.WithMessage(w, http.StatusInternalServerError, "Failed to fetch orders")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, orders)
}
