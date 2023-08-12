package order

import (
	"fmt"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

type OrderService interface {
	ResolveAllCart(userID uuid.UUID, limit, page int) ([]OrderResponse, error)
}

type OrderServiceImpl struct {
	OrderRepository   OrderRepository
	CartRepository    cart.CartRepository
	ProductRepository product.ProductRepository
	Config            *configs.Config
}

func ProvideOrderServiceImpl(orderRepository OrderRepository, config *configs.Config) *OrderServiceImpl {
	return &OrderServiceImpl{OrderRepository: orderRepository, Config: config}
}

func (o *OrderServiceImpl) ResolveAllCart(userID uuid.UUID, limit, page int) ([]OrderResponse, error) {
	orders, err := o.OrderRepository.ResolveAllOrderByUserID(userID, limit, page)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	var resp []OrderResponse
	for _, order := range orders {
		orderItems, err := o.OrderRepository.ResolveOrderItemsByOrderID(order.OrderID)
		if err != nil {
			logger.ErrorWithStack(err)
			return nil, fmt.Errorf("failed to fetch order items: %w", err)
		}

		orderResponse := order.BuildOrderResponse(order, orderItems)
		resp = append(resp, orderResponse)
	}

	return resp, nil
}
