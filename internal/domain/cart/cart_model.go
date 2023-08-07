package cart

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type (
	Cart struct {
		CartID    uuid.UUID   `db:"cart_id"`
		UserID    uuid.UUID   `db:"user_id"`
		CreatedAt time.Time   `db:"created_at"`
		CreatedBy uuid.UUID   `db:"created_by"`
		UpdatedAt null.Time   `db:"updated_at"`
		UpdatedBy nuuid.NUUID `db:"updated_by"`
		DeletedAt null.Time   `db:"deleted_at"`
		DeletedBy nuuid.NUUID `db:"deleted_by"`
		Items     []CartItems `db:"-"`
	}

	CartItems struct {
		CartItemID uuid.UUID   `db:"cart_item_id"`
		CartID     uuid.UUID   `db:"cart_id"`
		ProductID  uuid.UUID   `db:"product_id"`
		Quantity   float64     `db:"quantity"`
		CreatedAt  time.Time   `db:"created_at"`
		CreatedBy  uuid.UUID   `db:"created_by"`
		UpdatedAt  null.Time   `db:"updated_at"`
		UpdatedBy  nuuid.NUUID `db:"updated_by"`
		DeletedAt  null.Time   `db:"deleted_at"`
		DeletedBy  nuuid.NUUID `db:"deleted_by"`
	}
	Order struct {
		OrderID     uuid.UUID   `db:"order_id"`
		UserID      uuid.UUID   `db:"user_id"`
		TotalAmount float64     `db:"total_amount"`
		CreatedAt   time.Time   `db:"created_at"`
		CreatedBy   uuid.UUID   `db:"created_by"`
		UpdatedAt   null.Time   `db:"updated_at"`
		UpdatedBy   nuuid.NUUID `db:"updated_by"`
		DeletedAt   null.Time   `db:"deleted_at"`
		DeletedBy   nuuid.NUUID `db:"deleted_by"`
		Items       []OrderItem `db:"-"`
	}

	OrderItem struct {
		OrderItemID uuid.UUID   `db:"order_item_id"`
		OrderID     uuid.UUID   `db:"order_id"`
		ProductID   uuid.UUID   `db:"product_id"`
		Quantity    float64     `db:"quantity"`
		CreatedAt   time.Time   `db:"created_at"`
		CreatedBy   uuid.UUID   `db:"created_by"`
		UpdatedAt   null.Time   `db:"updated_at"`
		UpdatedBy   nuuid.NUUID `db:"updated_by"`
		DeletedAt   null.Time   `db:"deleted_at"`
		DeletedBy   nuuid.NUUID `db:"deleted_by"`
	}
)

type (
	AddToCartRequestFormat struct {
		ProductID uuid.UUID `json:"productID"`
		Quantity  float64   `json:"quantity"`
	}
	GetCartRequestFormat struct {
		CartID uuid.UUID `json:"cartID"`
	}
	CheckoutRequest struct {
		Items []CartItems
	}
)

type (
	CartResponseFormat struct {
		CartID    uuid.UUID                 `json:"cartID"`
		UserID    uuid.UUID                 `json:"userID"`
		CreatedAt time.Time                 `json:"createdAt"`
		CreatedBy uuid.UUID                 `json:"createdBy"`
		UpdatedAt null.Time                 `json:"updatedAt,omitempty"`
		UpdatedBy *uuid.UUID                `json:"updatedBy,omitempty"`
		DeletedAt null.Time                 `json:"deletedAt,omitempty"`
		DeletedBy *uuid.UUID                `json:"deletedBy,omitempty"`
		Items     []CartItemsResponseFormat `json:"items"`
	}
	CartItemsResponseFormat struct {
		CartItemID uuid.UUID  `json:"cartItemID"`
		CartID     uuid.UUID  `json:"cartID"`
		ProductID  uuid.UUID  `json:"productID"`
		Quantity   float64    `json:"quantity"`
		CreatedAt  time.Time  `json:"createdAt"`
		CreatedBy  uuid.UUID  `json:"createdBy"`
		UpdatedAt  null.Time  `json:"updatedAt,omitempty"`
		UpdatedBy  *uuid.UUID `json:"updatedBy,omitempty"`
		DeletedAt  null.Time  `json:"deletedAt,omitempty"`
		DeletedBy  *uuid.UUID `json:"deletedBy,omitempty"`
	}

	OrderResponseFormat struct {
		OrderID     uuid.UUID
		UserID      uuid.UUID
		TotalAmount float64
		CreatedAt   time.Time                 `json:"createdAt"`
		CreatedBy   uuid.UUID                 `json:"createdBy"`
		UpdatedAt   null.Time                 `json:"updatedAt,omitempty"`
		UpdatedBy   *uuid.UUID                `json:"updatedBy,omitempty"`
		DeletedAt   null.Time                 `json:"deletedAt,omitempty"`
		DeletedBy   *uuid.UUID                `json:"deletedBy,omitempty"`
		Items       []OrderItemResponseFormat `json:"items"`
	}
	OrderItemResponseFormat struct {
		OrderItemID uuid.UUID
		OrderID     uuid.UUID
		ProductID   uuid.UUID
		Quantity    float64
		CreatedAt   time.Time  `json:"createdAt"`
		CreatedBy   uuid.UUID  `json:"createdBy"`
		UpdatedAt   null.Time  `json:"updatedAt,omitempty"`
		UpdatedBy   *uuid.UUID `json:"updatedBy,omitempty"`
		DeletedAt   null.Time  `json:"deletedAt,omitempty"`
		DeletedBy   *uuid.UUID `json:"deletedBy,omitempty"`
	}
)

func (c *Cart) AttachItems(items []CartItems) Cart {
	for _, item := range items {
		if item.CartID == c.CartID {
			c.Items = append(c.Items, item)
		}
	}
	return *c
}

func (c Cart) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToResponseFormat())
}

func (ci CartItems) MarshalJSON() ([]byte, error) {
	return json.Marshal(ci.ToResponseFormat())
}

func (c Cart) ToResponseFormat() CartResponseFormat {
	resp := CartResponseFormat{
		CartID:    c.CartID,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy.Ptr(),
		DeletedAt: c.DeletedAt,
		DeletedBy: c.DeletedBy.Ptr(),
		Items:     make([]CartItemsResponseFormat, 0),
	}

	for _, item := range c.Items {
		resp.Items = append(resp.Items, item.ToResponseFormat())
	}

	return resp
}

func (ci *CartItems) ToResponseFormat() CartItemsResponseFormat {
	return CartItemsResponseFormat{
		CartItemID: ci.CartItemID,
		CartID:     ci.CartID,
		ProductID:  ci.ProductID,
		Quantity:   ci.Quantity,
		CreatedAt:  ci.CreatedAt,
		CreatedBy:  ci.CreatedBy,
	}
}

//Order

func (o Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ToResponseFormat())
}

func (oi OrderItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(oi.ToResponseFormat())
}

func (o Order) ToResponseFormat() OrderResponseFormat {
	resp := OrderResponseFormat{
		OrderID:   o.OrderID,
		UserID:    o.UserID,
		CreatedAt: o.CreatedAt,
		CreatedBy: o.CreatedBy,
		UpdatedAt: o.UpdatedAt,
		UpdatedBy: o.UpdatedBy.Ptr(),
		DeletedAt: o.DeletedAt,
		DeletedBy: o.DeletedBy.Ptr(),
		Items:     make([]OrderItemResponseFormat, 0),
	}

	for _, item := range o.Items {
		resp.Items = append(resp.Items, item.ToResponseFormat())
	}

	return resp
}

func (oi *OrderItem) ToResponseFormat() OrderItemResponseFormat {
	return OrderItemResponseFormat{
		OrderItemID: oi.OrderItemID,
		OrderID:     oi.OrderID,
		ProductID:   oi.ProductID,
		Quantity:    oi.Quantity,
		CreatedAt:   oi.CreatedAt,
		CreatedBy:   oi.CreatedBy,
	}
}

//Makes OrderResponse

type OrderResponse struct {
	ID         uuid.UUID       `json:"id"`
	TotalPrice float64         `json:"totalPrice"`
	UserID     uuid.UUID       `json:"userId"`
	CreatedAt  time.Time       `json:"createdAt"`
	CreatedBy  uuid.UUID       `json:"createdBy"`
	Items      []OrderItemInfo `json:"items"`
}

type OrderItemInfo struct {
	ID        uuid.UUID      `json:"id"`
	Quantity  float64        `json:"quantity"`
	ProductID uuid.UUID      `json:"productId"`
	OrderID   uuid.UUID      `json:"orderId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Product   ProductDetails `json:"product"`
}

type ProductDetails struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       float64   `json:"stock"`
	CategoryID  uuid.UUID `json:"categoryId"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   uuid.UUID `json:"createdBy"`
}

func (o Order) BuildOrderResponse(order Order, items []OrderItemInfo) OrderResponse {
	return OrderResponse{
		ID:         order.OrderID,
		TotalPrice: order.TotalAmount,
		UserID:     order.UserID,
		CreatedAt:  order.CreatedAt,
		Items:      items,
	}
}
