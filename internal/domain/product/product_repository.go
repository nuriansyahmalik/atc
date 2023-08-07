package product

//go:generate go run github.com/golang/mock/mockgen -source product_repository.go -destination mock/product_repository_mock.go -package product_mock

import (
	"database/sql"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	productQueries = struct {
		selectProduct  string
		insertProduct  string
		insertCategory string
	}{
		selectProduct: `
			SELECT 
			    p.product_id,
			    p.category_id,
				p.name,
				p.description,
				p.price,
				p.stock,
				p.created_at,
				p.created_by,
				p.updated_at,
				p.updated_by,
				p.deleted_at,
				p.deleted_by
			FROM product p`,
		insertProduct: `
			INSERT INTO product (
				product_id,
			    category_id,
			    name,
			    description,
			    price,
			    stock,
			    created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			) VALUES (
				:product_id,
			          :category_id,
			    :name,
			    :description,
			    :price,
			    :stock,
			    :created_at,
				:created_by,
				:updated_at,
				:updated_by,
				:deleted_at,
				:deleted_by)`,
		insertCategory: `
			INSERT INTO product_categories (
				category_id,
			    name,
			    description,
			    created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			) VALUES (
				:category_id,
			    :name,
			    :description,
			    :created_at,
				:created_by,
				:updated_at,
				:updated_by,
				:deleted_at,
				:deleted_by)`,
	}
)

type ProductRepository interface {
	Create(product Product) (err error)
	CreateCategory(category ProductCategories) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(productID uuid.UUID) (product Product, err error)
	ResolveProduct(limit, page int) (product []Product, err error)
	ResolveProductByCategory(limit, page int, categoryName string) (product []Product, err error)
	UpdateProductStock(productID uuid.UUID, stock float64) (err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepository(db *infras.MySQLConn) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{DB: db}
}

func (p *ProductRepositoryMySQL) Create(product Product) (err error) {
	stmt, err := p.DB.Write.PrepareNamed(productQueries.insertProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product)
	return err
}
func (u *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = u.DB.Read.Get(
		&exists,
		"SELECT COUNT(product_id) FROM product p WHERE p.product_id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (p *ProductRepositoryMySQL) ResolveProduct(limit, page int) (product []Product, err error) {
	query, args, err := sqlx.In(productQueries.selectProduct+" LIMIT ? OFFSET ?", limit, page)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	query = p.DB.Read.Rebind(query)
	var products []Product
	err = p.DB.Read.Select(&products, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	return products, nil
}
func (p *ProductRepositoryMySQL) ResolveProductByCategory(limit, page int, categoryName string) (product []Product, err error) {
	query, args, err := sqlx.In(productQueries.selectProduct+" WHERE category_id = (SELECT category_id FROM product_categories WHERE name = ? ) LIMIT ? OFFSET ?", categoryName, limit, limit*page)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	query = p.DB.Read.Rebind(query)
	var products []Product
	err = p.DB.Read.Select(&products, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	return products, nil
}
func (p *ProductRepositoryMySQL) ResolveByID(productID uuid.UUID) (product Product, err error) {
	err = p.DB.Read.Get(&product, productQueries.selectProduct+" WHERE product_id = ?", productID.String())
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("product")
		logger.ErrorWithStack(err)
		return
	}
	return
}
func (p *ProductRepositoryMySQL) CreateCategory(category ProductCategories) (err error) {
	stmt, err := p.DB.Write.PrepareNamed(productQueries.insertCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(category)
	return err
}
func (p *ProductRepositoryMySQL) UpdateProductStock(productID uuid.UUID, stock float64) (err error) {
	_, err = p.DB.Write.Exec("UPDATE product SET stock = ? WHERE product_id = ?", stock, productID)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}
