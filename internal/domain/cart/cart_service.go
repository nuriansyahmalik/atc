package cart

//go:generate go run github.com/golang/mock/mockgen -source cart_service.go -destination mock/cart_service_mock.go -package cart_mock

import (
	"database/sql"
	"fmt"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type CartService interface {
	CheckoutCarts(userID uuid.UUID, cartItems []CartItems) (orders OrderResponse, err error)
	AddItemToCart(requestFormat AddToCartRequestFormat, userID uuid.UUID) (cart Cart, err error)
	ResolveCartByID(requestFormat GetCartRequestFormat, userID uuid.UUID) (cart Cart, err error)
}

type CartServiceImpl struct {
	CartRepository    CartRepository
	ProductRepository product.ProductRepository
	Config            *configs.Config
}

func ProvideCarServiceImpl(cartRepository CartRepository, productRepository product.ProductRepository, config *configs.Config) *CartServiceImpl {
	return &CartServiceImpl{CartRepository: cartRepository, ProductRepository: productRepository, Config: config}
}

func (c *CartServiceImpl) AddItemToCart(req AddToCartRequestFormat, userID uuid.UUID) (cart Cart, err error) {
	productID := req.ProductID
	quantity := req.Quantity

	product, err := c.ProductRepository.ResolveByID(productID)
	if err != nil {
		return cart, err
	}

	if product.Stock < quantity {
		log.Info().Msg("insufficient stock")
		return
	}

	cart, err = c.CartRepository.ResolveCartByID(userID)
	if err != sql.ErrNoRows && err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if err == sql.ErrNoRows {
		cartID, err := uuid.NewV4()
		if err != nil {
			return cart, nil
		}
		cart = Cart{
			CartID:    cartID,
			UserID:    userID,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		}
		if err := c.CartRepository.CreateCart(cart); err != nil {
			return cart, nil
		}
	}

	existingItem, err := c.CartRepository.ResolveCartItemByProduct(cart.CartID, productID)
	if err != nil {
		return cart, nil
	}

	if existingItem == nil {
		cartItemID, err := uuid.NewV4()
		if err != nil {
			return cart, nil
		}
		cartItem := CartItems{
			CartItemID: cartItemID,
			CartID:     cart.CartID,
			ProductID:  productID,
			Quantity:   quantity,
			CreatedAt:  time.Now(),
			CreatedBy:  userID,
		}
		if err := c.CartRepository.CreateCartItems(cartItem); err != nil {
			return cart, nil
		}
	} else {
		existingItem[0].Quantity += quantity
		existingItem[0].CreatedAt = time.Now()
		if err := c.CartRepository.UpdateCartItem(existingItem[0]); err != nil {
			return cart, err
		}
	}

	items, err := c.CartRepository.ResolveCartItemsByCartID(cart.CartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	cart.AttachItems(items)
	return
}
func (c CartServiceImpl) ResolveCartByID(requestFormat GetCartRequestFormat, userID uuid.UUID) (cart Cart, err error) {
	cart, err = c.CartRepository.ResolveCartByID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return cart, err
	}

	cartItems, err := c.CartRepository.ResolveCartItemsByCartID(requestFormat.CartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return cart, err
	}

	cart.AttachItems(cartItems)

	return
}
func (c *CartServiceImpl) CheckoutCarts(userID uuid.UUID) (OrderResponse, error) {
	cart, err := c.CartRepository.ResolveCartByID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	if cart.CartID == uuid.Nil {
		return OrderResponse{}, failure.InternalError(err)
	}

	cartItems, err := c.CartRepository.ResolveCartItemsByCartID(cart.CartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	if len(cartItems) == 0 {
		return OrderResponse{}, fmt.Errorf("No items in the cart")
	}

	totalAmount, items, err := c.calculateTotalAndItems(cartItems)
	if err != nil {
		return OrderResponse{}, err
	}

	order, err := c.createOrder(userID, totalAmount)
	if err != nil {
		return OrderResponse{}, err
	}

	if err := c.processOrderItemsAndStock(order, cartItems); err != nil {
		return OrderResponse{}, err
	}

	if err := c.CartRepository.ClearCart(cart.CartID); err != nil {
		return OrderResponse{}, err
	}

	return order.BuildOrderResponse(order, items), nil
}

func (c *CartServiceImpl) createOrder(userID uuid.UUID, totalAmount float64) (Order, error) {
	orderID, err := uuid.NewV4()
	if err != nil {
		return Order{}, err
	}
	order := Order{
		OrderID:     orderID,
		UserID:      userID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}
	if err := c.CartRepository.CreateOrder(order); err != nil {
		return Order{}, err
	}
	return order, nil
}

func (c *CartServiceImpl) calculateTotalAmount(cartItems []CartItems) (float64, error) {
	var totalAmount float64
	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			return 0, err
		}
		totalAmount += cartItem.Quantity * product.Price
	}
	return totalAmount, nil
}

func (c *CartServiceImpl) checkProductStock(cartItems []CartItems) error {
	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < cartItem.Quantity {
			return fmt.Errorf("insufficient stock for product: %s", product.Name)
		}
	}
	return nil
}

func (c *CartServiceImpl) processOrderItemsAndStock(order Order, cartItems []CartItems) error {
	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			logger.ErrorWithStack(err)
			return err
		}
		if product.Stock < cartItem.Quantity {
			return fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		stock := product.Stock - cartItem.Quantity
		if err := c.ProductRepository.UpdateProductStock(cartItem.ProductID, stock); err != nil {
			logger.ErrorWithStack(err)
			return err
		}

		orderItemID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		if err := c.CartRepository.CreateOrderItem(OrderItem{
			OrderItemID: orderItemID,
			OrderID:     order.OrderID,
			ProductID:   cartItem.ProductID,
			Quantity:    cartItem.Quantity,
			CreatedAt:   time.Now(),
			CreatedBy:   order.CreatedBy,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *CartServiceImpl) calculateTotalAndItems(cartItems []CartItems) (float64, []OrderItemInfo, error) {
	var totalAmount float64
	orderItemsInfo := make([]OrderItemInfo, 0)

	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			return 0, nil, err
		}

		itemTotal := cartItem.Quantity * product.Price
		totalAmount += itemTotal

		orderItemInfo := OrderItemInfo{
			ID:       cartItem.CartItemID,
			Quantity: cartItem.Quantity,
			Product: ProductDetails{
				ID:          product.ProductID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Stock:       product.Stock,
				CategoryID:  product.CategoryID,
				CreatedAt:   product.CreatedAt,
			},
		}
		orderItemsInfo = append(orderItemsInfo, orderItemInfo)
	}

	return totalAmount, orderItemsInfo, nil
}
