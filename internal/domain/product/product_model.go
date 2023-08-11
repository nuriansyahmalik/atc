package product

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Product struct {
	ProductID   uuid.UUID   `db:"product_id"`
	CategoryID  uuid.UUID   `db:"category_id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	Price       float64     `db:"price"`
	Stock       float64     `db:"stock"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   uuid.UUID   `db:"created_by"`
	UpdatedAt   null.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
}

type ProductCategories struct {
	ID          uuid.UUID   `db:"category_id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   uuid.UUID   `db:"created_by"`
	UpdatedAt   null.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
}

type (
	ProductRequestFormat struct {
		ID          uuid.UUID `json:"ID"`
		CategoryID  uuid.UUID `json:"categoryID" validate:"required"`
		ProductName string    `json:"productName" validate:"required"`
		Description string    `json:"description" validate:"required"`
		Price       float64   `json:"price" validate:"required"`
		Stock       float64   `json:"stock" validate:"required"`
	}
	ProductResponseFormat struct {
		ID          uuid.UUID `json:"ID,omitempty"`
		CategoryID  uuid.UUID `json:"categoryID,omitempty"`
		ProductName string    `json:"productName,omitempty"`
		Description string    `json:"description,omitempty"`
		Price       float64   `json:"price,omitempty"`
		Stock       float64   `json:"stock,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		CreatedBy   uuid.UUID `json:"createdBy"`
	}
	CategoriesRequestFormat struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	CategoryResponseFormat struct {
		ID           uuid.UUID `json:"categoryID,omitempty"`
		CategoryName string    `json:"categoryName,omitempty"`
		Description  string    `json:"description,omitempty"`
		CreatedAt    time.Time `json:"createdAt"`
		CreatedBy    uuid.UUID `json:"createdBy"`
	}
)

func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}

func (pc ProductCategories) MarshalJSON() ([]byte, error) {
	return json.Marshal(pc.ToResponseFormat())
}

// Validate validates the entity.
func (p *Product) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(p)
}

func (p Product) ProductRequestFormat(req ProductRequestFormat, userID uuid.UUID) (product Product, err error) {
	productID, err := uuid.NewV4()
	if err != nil {
		return
	}
	product = Product{
		ProductID:   productID,
		CategoryID:  req.CategoryID,
		Name:        req.ProductName,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}
	return
}

func (pc ProductCategories) CategoryRequestFormat(req CategoriesRequestFormat, userID uuid.UUID) (prodCategory ProductCategories, err error) {
	categoryID, err := uuid.NewV4()
	if err != nil {
		return
	}
	prodCategory = ProductCategories{
		ID:          categoryID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}
	return
}

func (p Product) ToResponseFormat() ProductResponseFormat {
	return ProductResponseFormat{
		ID:          p.ProductID,
		CategoryID:  p.CategoryID,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		CreatedBy:   p.CreatedBy,
	}
}
func (pc ProductCategories) ToResponseFormat() CategoryResponseFormat {
	return CategoryResponseFormat{
		ID:           pc.ID,
		CategoryName: pc.Name,
		Description:  pc.Description,
		CreatedAt:    pc.CreatedAt,
		CreatedBy:    pc.CreatedBy,
	}
}
