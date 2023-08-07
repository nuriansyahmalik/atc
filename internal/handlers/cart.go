package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
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
		r.Post("/", h.Checkout)
		r.Get("/", h.GetCartByID)

	})
}
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

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok {
		http.Error(w, "Error Claims", http.StatusUnauthorized)
		return
	}
	checkout, err := h.CartService.CheckoutCarts(claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, checkout)
}

func (h *CartHandler) GetCartByID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat cart.GetCartRequestFormat
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
		http.Error(w, "Error claims", http.StatusUnauthorized)
		return
	}
	cart, err := h.CartService.ResolveCartByID(requestFormat, claims.ID)
	if err != nil {
		return
	}
	response.WithJSON(w, http.StatusOK, cart)
}
