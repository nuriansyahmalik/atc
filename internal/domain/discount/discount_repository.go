package discount

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
)

var discountQueries = struct {
	insertDiscount string
	selectDiscount string
}{
	insertDiscount: `
		INSERT INTO discounts (discount_id, code, type, price, start_date, end_date, created_at, created_by)
		VALUES (:discount_id, :code, :type, :price, :start_date, :end_date, :created_at, :created_by)`,
	selectDiscount: `
		SELECT
			d.discount_id,
			d.code,
			d.type,
			d.price,
			d.start_date,
			d.end_date,
			d.created_at,
			d.created_by,
			d.updated_at,
			d.updated_by,
			d.deleted_at,
			d.deleted_by
		FROM discounts d`,
}

type DiscountRepository interface {
	CreateDiscount(discount Discount) error
	Delete(discount Discount) error
	ResolveByID(discountID string) (Discount, error)
	ResolveByCode(code string) (Discount, error)
}

type DiscountRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideDiscountRepositoryMySQL(db *infras.MySQLConn) *DiscountRepositoryMySQL {
	return &DiscountRepositoryMySQL{DB: db}
}

func (d *DiscountRepositoryMySQL) ResolveByID(discountID string) (Discount, error) {
	query := discountQueries.selectDiscount + " WHERE d.discount_id = ?"
	var discount Discount
	err := d.DB.Read.Get(&discount, query, discountID)
	if err != nil {
		logger.ErrorWithStack(err)
		return Discount{}, err
	}
	return discount, nil
}

func (d *DiscountRepositoryMySQL) ResolveByCode(code string) (Discount, error) {
	query := discountQueries.selectDiscount + " WHERE d.code = ?"
	var discount Discount
	err := d.DB.Read.Get(&discount, query, code)
	if err != nil {
		return Discount{}, err
	}
	return discount, nil
}

func (d *DiscountRepositoryMySQL) CreateDiscount(discount Discount) error {
	_, err := d.DB.Write.NamedExec(discountQueries.insertDiscount, discount)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiscountRepositoryMySQL) Delete(discount Discount) error {

	return nil
}
