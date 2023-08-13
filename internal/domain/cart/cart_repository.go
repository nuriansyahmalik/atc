package cart

//go:generate go run github.com/golang/mock/mockgen -source cart_repository.go -destination mock/cart_repository_mock.go -package cart_mock

import (
	"database/sql"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	cartQueries = struct {
		insertCart       string
		insertCartItems  string
		insertOrder      string
		insertOrderItems string
		selectCarts      string
		selectCartItems  string
		updateCartItems  string
		deleteCartItems  string
	}{
		insertCart: `
			INSERT INTO carts (
			    cart_id,
			    user_id,
			    created_at,
				created_by,
			    updated_at,
			    updated_by               
			)VALUES(
			    :cart_id,
			    :user_id,
			    :created_at,
				:created_by,
			    :updated_at,
			    :updated_by) `,
		insertCartItems: `
			INSERT INTO cart_items (
			    cart_item_id, 
			    cart_id,
			    product_id,
			    quantity,
			    created_at,
				created_by
			)VALUES(
			    :cart_item_id, 
			    :cart_id,
			    :product_id,
			    :quantity,
			    :created_at,
				:created_by)`,
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
			INSERT INTO order_items (
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
		selectCarts: `
			SELECT 
			    c.cart_id,
			    c.user_id,
			    c.created_at,
				c.created_by,
				c.updated_at,
				c.updated_by,
				c.deleted_at,
				c.deleted_by
			FROM carts c`,
		selectCartItems: `
			SELECT 
			    ci.cart_item_id, 
			    ci.cart_id,
			    ci.product_id,
			    ci.quantity,
			    ci.created_at,
				ci.created_by,
				ci.updated_at,
				ci.updated_by,
				ci.deleted_at,
				ci.updated_by
			FROM cart_items ci`,
		updateCartItems: `
			UPDATE 
			    cart_items 
			SET quantity = :quantity
			WHERE cart_id = :cart_id AND product_id = :product_id`,
	}
)

type CartRepository interface {
	CreateCart(cart Cart) (err error)
	CreateCartItems(cartItems CartItems) (err error)
	CreateOrder(order Order) (err error)
	CreateOrderItem(orderItem OrderItem) (err error)
	ClearCart(cartID uuid.UUID) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveCartByID(userID uuid.UUID) (cart Cart, err error)
	ResolveCartItemsByCartID(cartID uuid.UUID) (cartItem []CartItems, err error)
	ResolveCartItemByProduct(cartID uuid.UUID, productID uuid.UUID) (cartItems []CartItems, err error)
	UpdateCartItem(cartItems CartItems) (err error)
	RemoveItemFromCart(cartItems CartItems) (err error)
}
type CartRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCartRepositoryMySQL(db *infras.MySQLConn) *CartRepositoryMySQL {
	return &CartRepositoryMySQL{DB: db}
}

func (c *CartRepositoryMySQL) CreateCart(cart Cart) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txCreateCart(tx, cart); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (c *CartRepositoryMySQL) CreateCartItems(cartItems CartItems) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txCreateCartItems(tx, cartItems); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}
func (c *CartRepositoryMySQL) CreateOrder(order Order) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txCreateOrder(tx, order); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (c *CartRepositoryMySQL) CreateOrderItem(orderItem OrderItem) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txCreateOrderItems(tx, orderItem); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}
func (c *CartRepositoryMySQL) ClearCart(cartID uuid.UUID) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txDeleteCart(tx, cartID); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}
func (c *CartRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = c.DB.Read.Get(
		&exists,
		"SELECT COUNT(cart_id) FROM carts WHERE carts.cart_id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
func (c *CartRepositoryMySQL) ResolveCartByID(userID uuid.UUID) (cart Cart, err error) {
	err = c.DB.Read.Get(&cart, cartQueries.selectCarts+" WHERE c.user_id = ?", userID)
	if err != nil && err == sql.ErrNoRows {
		// err = failure.NotFound("cart")
		return
	}
	return
}
func (c *CartRepositoryMySQL) ResolveCartItemsByCartID(cartID uuid.UUID) (cartItems []CartItems, err error) {
	query, args, err := sqlx.In(cartQueries.selectCartItems+" WHERE ci.cart_id = ?", cartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	err = c.DB.Read.Select(&cartItems, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
func (c *CartRepositoryMySQL) ResolveCartItemByProduct(cartID uuid.UUID, productID uuid.UUID) (cartItems []CartItems, err error) {
	query, args, err := sqlx.In(cartQueries.selectCartItems+" WHERE ci.cart_id = ? AND ci.product_id = ? ", cartID, productID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	err = c.DB.Read.Select(&cartItems, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
func (c *CartRepositoryMySQL) UpdateCartItem(cart CartItems) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txUpdateCartItems(tx, cart); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}
func (c *CartRepositoryMySQL) RemoveItemFromCart(cart CartItems) (err error) {
	return c.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := c.txDeleteCartItems(tx, cart); err != nil {
			e <- err
			return
		}
		e <- err
	})
}

func (c *CartRepositoryMySQL) txCreateCart(tx *sqlx.Tx, cart Cart) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertCart)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cart)
	if err != nil {
		return err
	}

	return nil
}
func (c *CartRepositoryMySQL) txCreateCartItems(tx *sqlx.Tx, cartItems CartItems) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertCartItems)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(cartItems)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
func (c *CartRepositoryMySQL) txUpdateCartItems(tx *sqlx.Tx, cartItems CartItems) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.updateCartItems)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cartItems)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return
}
func (c *CartRepositoryMySQL) txCreateOrder(tx *sqlx.Tx, order Order) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertOrder)
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
func (c *CartRepositoryMySQL) txCreateOrderItems(tx *sqlx.Tx, orderItems OrderItem) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertOrderItems)
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
func (c *CartRepositoryMySQL) txDeleteCart(tx *sqlx.Tx, cartID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID.String())
	return
}
func (c *CartRepositoryMySQL) txDeleteCartItems(tx *sqlx.Tx, cartItems CartItems) (err error) {
	_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?", cartItems.CartID, cartItems.ProductID)
	if err != nil {
		return err
	}
	return nil
}
