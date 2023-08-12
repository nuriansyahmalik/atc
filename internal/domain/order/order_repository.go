package order

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	orderQueries = struct {
		insertOrder      string
		insertOrderItems string
		selectOrder      string
		selectOrderItems string
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
		selectOrder: `
			SELECT
			    o.order_id,
			    o.user_id,
			    o.total_amount,
			    o.created_at, 
				o.created_by, 
				o.updated_at, 
				o.updated_by, 
				o.deleted_at, 
				o.deleted_by 
			FROM orders o`,
		selectOrderItems: `
			SELECT 
			oi.order_item_id,
			oi.order_id,
			oi.product_id,
			oi.quantity,
			oi.created_at, 
			oi.created_by, 
			oi.updated_at, 
			oi.updated_by, 
			oi.deleted_at, 
			oi.deleted_by 
		FROM order_items oi`,
	}
)

type OrderRepository interface {
	CreateOrder(order Order) (err error)
	CreateOrderItem(orderItem OrderItem) (err error)
	ResolveAllOrderByUserID(userID uuid.UUID, limit, page int) ([]Order, error)
	ResolveOrderItemsByOrderID(orderID uuid.UUID) ([]OrderItemInfo, error)
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
func (o *OrderRepositoryMySQL) ResolveAllOrderByUserID(userID uuid.UUID, limit, page int) ([]Order, error) {
	query := o.DB.Read.Rebind(orderQueries.selectOrder + " WHERE o.user_id = ? LIMIT ? OFFSET ?")
	var orders []Order
	err := o.DB.Read.Select(&orders, query, userID, limit, page)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepositoryMySQL) ResolveOrderItemsByOrderID(orderID uuid.UUID) ([]OrderItemInfo, error) {
	query := o.DB.Read.Rebind(orderQueries.selectOrderItems + " WHERE oi.order_id = ?")
	var orderItems []OrderItemInfo
	err := o.DB.Read.Select(&orderItems, query, orderID)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	return orderItems, nil
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
