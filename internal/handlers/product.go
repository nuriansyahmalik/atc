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
		//r.Use(middleware.ValidateJWTMiddleware)
		r.Post("/", h.CreateProduct)
		r.Post("/category", h.CreateCategory)
		r.Get("/", h.GetAllProduct)
	})
}

// CreateProduct create a new product
// @Summary Create a new product
// @Description this endpoint create a new product
// @Tags product/product
// @Security JWTAuthentication
// @Param user body product.ProductRequestFormat true "The Product to be created."
// @Produce json
// @Success 201 {object} response.Base{data=product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/ [post]
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

// GetAllProduct resolves a Product by its ID.
// @Summary Retrieve list of products
// @Description This endpoint retrieves a list of products with optional filters.
// @Tags product/product
// @Security JWTAuthentication
// @Param limit query int false "The number of products per page."
// @Param page query int false "The page number."
// @Param category query string false "Filter products by category."
// @Produce json
// @Success 200 {object} response.Base{data=product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product [get]
func (h *ProductHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
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
	categoryName := r.URL.Query().Get("category")

	var products []product.Product
	if categoryName == "" {
		products, err = h.ProductService.ResolveProduct(limit, page-1)
		if err != nil {
			response.WithMessage(w, http.StatusBadRequest, "Missing Query Param")
			response.WithError(w, err)
			return
		}
	} else {
		products, err = h.ProductService.ResolveProductByCategory(limit, page-1, categoryName)
		if err != nil {
			response.WithMessage(w, http.StatusBadRequest, "Missing Param Query / Wrong Category Name")
			response.WithError(w, err)
			return
		}
	}
	response.WithJSON(w, http.StatusOK, products)
}

// CreateCategory create a new product categories
// @Summary Create a new product categories
// @Description this endpoint create a new product categories
// @Tags product/product
// @Security JWTAuthentication
// @Param user body product.CategoriesRequestFormat true "The Product Categories to be created."
// @Produce json
// @Success 201 {object} response.Base{data=product.CategoryResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/ [post]
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
