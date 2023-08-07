package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductService product.ProductService
}

func ProvideProductHandler(productService product.ProductService) ProductHandler {
	return ProductHandler{ProductService: productService}
}

func (h *ProductHandler) Router(r chi.Router) {
	r.Route("/product", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Post("/", h.CreateProduct)
		r.Post("/category", h.CreateCategory)
		r.Get("/", h.GetAllProduct)
	})
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat product.ProductRequestFormat
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
	product, err := h.ProductService.Create(requestFormat, claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	categoryName := r.URL.Query().Get("category")

	var products []product.Product
	if categoryName == "" {
		products, err = h.ProductService.ResolveProduct(limit, page-1)
		if err != nil {
			response.WithError(w, err)
			return
		}
	} else {
		products, err = h.ProductService.ResolveProductByCategory(limit, page-1, categoryName)
		if err != nil {
			response.WithError(w, err)
			return
		}
	}
	response.WithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat product.CategoriesRequestFormat
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
	prodCategory, err := h.ProductService.CreateCategory(requestFormat, claims.ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, prodCategory)
}
