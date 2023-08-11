package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type CartHandler struct {
	CartService cart.CartService
}

func ProvideCartHandler(cartService cart.CartService) CartHandler {
	return CartHandler{CartService: cartService}
}

func (h *CartHandler) Router(r chi.Router) {
	r.Route("/cart", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Post("/add", h.AddToCart)
		r.Post("/checkout", h.Checkout)
		r.Get("/{id}", h.GetCartByID)

	})
}

// AddToCart create a new cart
// @Summary Create a new cart
// @Description this endpoint create a new cart
// @Tags cart/cart
// @Security JWTAuthentication
// @Param user body cart.AddToCartRequestFormat true "The Cart to be created."
// @Produce json
// @Success 201 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/add [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat cart.AddToCartRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	cart, err := h.CartService.AddItemToCart(requestFormat, claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, cart)
}

// Checkout from cart
// @Summary Create a new order from cart
// @Description this endpoint create a new order from cart
// @Tags cart/cart
// @Security JWTAuthentication
// @Param user body cart.CheckoutRequestFormat true "The Order to be created."
// @Produce json
// @Success 201 {object} response.Base{data=cart.OrderResponse}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/checkout [post]
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat cart.CheckoutRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	checkout, err := h.CartService.CheckoutCarts(requestFormat, claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, checkout)
}

// GetCartByID resolves a Cart by its ID.
// @Summary Resolve Cart by ID
// @Description This endpoint resolves a Cart by its ID.
// @Tags cart/cart
// @Security JWTAuthentication
// @Param id path string true "The cart's identifier."
// @Produce json
// @Success 200 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/{id} [get]
func (h *CartHandler) GetCartByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error claims", http.StatusUnauthorized)
		return
	}
	cart, err := h.CartService.ResolveCartByID(id, claims.ID)
	if err != nil {
		return
	}

	response.WithJSON(w, http.StatusOK, cart)
}
