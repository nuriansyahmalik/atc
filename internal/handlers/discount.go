package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/discount"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"net/http"
)

type DiscountHandler struct {
	DiscountService discount.DiscountService
}

func ProvideDiscountHandler(discountService discount.DiscountService) DiscountHandler {
	return DiscountHandler{DiscountService: discountService}
}

func (h *DiscountHandler) Router(r chi.Router) {
	r.Route("/discount", func(r chi.Router) {
		r.Use(middleware.ValidateJWTMiddleware)
		//r.Use(jwt.AuthMiddleware)
		r.Post("/", h.CreateDiscount)
	})
}

// CreateDiscount create a new discount
// @Summary Create a new discount
// @Description this endpoint create a new discount
// @Tags discount/discount
// @Security JWTAuthentication
// @Param user body discount.DiscountRequestFormat true "The Product to be created."
// @Produce json
// @Success 201 {object} response.Base{data=discount.DiscountResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/ [post]
func (h *DiscountHandler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat discount.DiscountRequestFormat
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
	discount, err := h.DiscountService.CreateDiscount(requestFormat, claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, discount)
}