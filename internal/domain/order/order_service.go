package order

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/gofrs/uuid"
)

type OrderService interface {
	CheckoutCart(userID uuid.UUID) (order Order, err error)
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

func (o *OrderServiceImpl) CheckoutCart(userID uuid.UUID) (order Order, err error) {
	return Order{}, err
}

//func (o *OrderServiceImpl) CheckoutCarts(userID uuid.UUID) (OrderResponse, error) {
//	cart, err := o.CartRepository.ResolveCartByID(userID)
//	if err != nil {
//		logger.ErrorWithStack(err)
//		return OrderResponse{}, err
//	}
//
//	if cart.CartID == uuid.Nil {
//		return OrderResponse{}, failure.InternalError(err)
//	}
//
//	cartItems, err := o.CartRepository.ResolveCartItemsByCartID(cart.CartID)
//	if err != nil {
//		logger.ErrorWithStack(err)
//		return OrderResponse{}, err
//	}
//
//	if len(cartItems) == 0 {
//		return OrderResponse{}, fmt.Errorf("No items in the cart")
//	}
//
//	totalAmount, items, err := o.calculateTotalAndItems(cartItems)
//	if err != nil {
//		return OrderResponse{}, err
//	}
//
//	order, err := o.createOrder(userID, totalAmount)
//	if err != nil {
//		return OrderResponse{}, err
//	}
//
//	if err := o.processOrderItemsAndStock(order, cartItems); err != nil {
//		return OrderResponse{}, err
//	}
//
//	if err := o.CartRepository.ClearCart(cart.CartID); err != nil {
//		return OrderResponse{}, err
//	}
//
//	return order.BuildOrderResponse(order, items), nil
//}
//
//func (o *OrderServiceImpl) createOrder(userID uuid.UUID, totalAmount float64) (Order, error) {
//	orderID, err := uuid.NewV4()
//	if err != nil {
//		return Order{}, err
//	}
//	order := Order{
//		OrderID:     orderID,
//		UserID:      userID,
//		TotalAmount: totalAmount,
//		CreatedAt:   time.Now(),
//		CreatedBy:   userID,
//	}
//	if err := o.CartRepository.CreateOrder(order); err != nil {
//		return Order{}, err
//	}
//	return order, nil
//}
//
//func (o *OrderServiceImpl) calculateTotalAmount(cartItems []cart.CartItems) (float64, error) {
//	var totalAmount float64
//	for _, cartItem := range cartItems {
//		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
//		if err != nil {
//			return 0, err
//		}
//		totalAmount += cartItem.Quantity * product.Price
//	}
//	return totalAmount, nil
//}
//
//func (o *OrderServiceImpl) checkProductStock(cartItems []cart.CartItems) error {
//	for _, cartItem := range cartItems {
//		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
//		if err != nil {
//			return err
//		}
//		if product.Stock < cartItem.Quantity {
//			return fmt.Errorf("insufficient stock for product: %s", product.Name)
//		}
//	}
//	return nil
//}
//
//func (o *OrderServiceImpl) processOrderItemsAndStock(order Order, cartItems []cart.CartItems) error {
//	for _, cartItem := range cartItems {
//		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
//		if err != nil {
//			logger.ErrorWithStack(err)
//			return err
//		}
//		if product.Stock < cartItem.Quantity {
//			return fmt.Errorf("insufficient stock for product: %s", product.Name)
//		}
//
//		stock := product.Stock - cartItem.Quantity
//		if err := c.ProductRepository.UpdateProductStock(cartItem.ProductID, stock); err != nil {
//			logger.ErrorWithStack(err)
//			return err
//		}
//
//		orderItemID, err := uuid.NewV4()
//		if err != nil {
//			return err
//		}
//		if err := c.CartRepository.CreateOrderItem(OrderItem{
//			OrderItemID: orderItemID,
//			OrderID:     order.OrderID,
//			ProductID:   cartItem.ProductID,
//			Quantity:    cartItem.Quantity,
//			CreatedAt:   time.Now(),
//			CreatedBy:   order.CreatedBy,
//		}); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (o *OrderServiceImpl) calculateTotalAndItems(cartItems []cart.CartItems) (float64, []OrderItemInfo, error) {
//	var totalAmount float64
//	orderItemsInfo := make([]OrderItemInfo, 0)
//
//	for _, cartItem := range cartItems {
//		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
//		if err != nil {
//			return 0, nil, err
//		}
//
//		itemTotal := cartItem.Quantity * product.Price
//		totalAmount += itemTotal
//
//		orderItemInfo := OrderItemInfo{
//			ID:       cartItem.CartItemID,
//			Quantity: cartItem.Quantity,
//			Product: ProductDetails{
//				ID:          product.ProductID,
//				Name:        product.Name,
//				Description: product.Description,
//				Price:       product.Price,
//				Stock:       product.Stock,
//				CategoryID:  product.CategoryID,
//				CreatedAt:   product.CreatedAt,
//			},
//		}
//		orderItemsInfo = append(orderItemsInfo, orderItemInfo)
//	}
//
//	return totalAmount, orderItemsInfo, nil
//}
