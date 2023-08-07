package order

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/jmoiron/sqlx"
)

var (
	orderQueries = struct {
		insertOrder      string
		insertOrderItems string
	}{
		insertOrder: `
			INSERT INTO orders (
                order_id,
			    user_id,
                total_amount,
				created_at,
				created_by
			)VALUES(
			    :order_id,
			    :user_id,
                :total_amount,
			    :created_at,
				:created_by)`,
		insertOrderItems: `
			INSERT INTO order_items(
                order_item_id,
			    order_id,
			    product_id,
			    quantity,
			    created_at,
				created_by             
			) VALUES (
				:order_item_id,
			    :order_id,
			    :product_id,
			    :quantity,
			    :created_at,
				:created_by)`,
	}
)

type OrderRepository interface {
	CreateOrder(order Order) (err error)
	CreateOrderItem(orderItem OrderItem) (err error)
}

type OrderRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideOrderRepositoryMySQL(db *infras.MySQLConn) *OrderRepositoryMySQL {
	return &OrderRepositoryMySQL{DB: db}
}

func (o *OrderRepositoryMySQL) CreateOrder(order Order) (err error) {
	return o.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := o.txCreateOrder(tx, order); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}
func (o *OrderRepositoryMySQL) CreateOrderItem(orderItem OrderItem) (err error) {
	return o.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := o.txCreateOrderItems(tx, orderItem); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (o *OrderRepositoryMySQL) txCreateOrder(tx *sqlx.Tx, order Order) (err error) {
	stmt, err := tx.PrepareNamed(orderQueries.insertOrder)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(order)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
func (o *OrderRepositoryMySQL) txCreateOrderItems(tx *sqlx.Tx, orderItems OrderItem) (err error) {
	stmt, err := tx.PrepareNamed(orderQueries.insertOrderItems)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(orderItems)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
