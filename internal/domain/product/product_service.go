package product

//go:generate go run github.com/golang/mock/mockgen -source product_service.go -destination mock/product_service_mock.go -package product_mock

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error)
	CreateCategory(requestFormat CategoriesRequestFormat, userID uuid.UUID) (prodCategory ProductCategories, err error)
	ResolveProduct(limit, page int) (product []Product, err error)
	ResolveProductByCategory(limit, page int, categoryName string) (product []Product, err error)
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	Config            *configs.Config
}

func ProvideProductServiceImpl(productRepository ProductRepository, config *configs.Config) *ProductServiceImpl {
	return &ProductServiceImpl{ProductRepository: productRepository, Config: config}
}

func (p *ProductServiceImpl) Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error) {
	product, err = product.ProductRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}
	if err != nil {
		return product, failure.BadRequest(err)
	}
	err = p.ProductRepository.Create(product)
	if err != nil {
		return
	}
	return
}

func (p *ProductServiceImpl) ResolveProduct(limit, page int) (product []Product, err error) {
	product, err = p.ProductRepository.ResolveProduct(limit, page)
	if err != nil {
		return nil, err
	}
	return
}

func (p *ProductServiceImpl) ResolveProductByCategory(limit, page int, categoryName string) (product []Product, err error) {
	product, err = p.ProductRepository.ResolveProductByCategory(limit, page, categoryName)
	if err != nil {
		return nil, err
	}
	return
}

func (p *ProductServiceImpl) CreateCategory(requestFormat CategoriesRequestFormat, userID uuid.UUID) (prodCategory ProductCategories, err error) {
	prodCategory, err = prodCategory.CategoryRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}
	if err != nil {
		return prodCategory, failure.BadRequest(err)
	}
	err = p.ProductRepository.CreateCategory(prodCategory)
	if err != nil {
		return
	}
	return
}
