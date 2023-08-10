package cart

//go:generate go run github.com/golang/mock/mockgen -source cart_service.go -destination mock/cart_service_mock.go -package cart_mock

import (
	"database/sql"
	"fmt"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type CartService interface {
	CheckoutCarts(requestFormat CheckoutRequestFormat, userID uuid.UUID) (orders OrderResponse, err error)
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

	product, err := c.ProductRepository.ResolveByID(req.ProductID)
	if err != nil {
		return cart, err
	}

	if product.Stock < req.Quantity {
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

		if err := c.CartRepository.CreateCart(Cart{
			CartID:    cartID,
			UserID:    userID,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		}); err != nil {
			return cart, nil
		}
	}

	existingItem, err := c.CartRepository.ResolveCartItemByProduct(cart.CartID, req.ProductID)
	if err != nil {
		return cart, nil
	}

	if existingItem == nil {
		cartItemID, err := uuid.NewV4()
		if err != nil {
			return cart, nil
		}
		if err := c.CartRepository.CreateCartItems(CartItems{
			CartItemID: cartItemID,
			CartID:     cart.CartID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			CreatedAt:  time.Now(),
			CreatedBy:  userID,
		}); err != nil {
			return cart, nil
		}
	} else {
		existingItem[0].Quantity += req.Quantity
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

func (c *CartServiceImpl) CheckoutCarts(requestFormat CheckoutRequestFormat, userID uuid.UUID) (OrderResponse, error) {
	// Resolve cart for the given userID
	cart, err := c.CartRepository.ResolveCartByID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	// Check if the cart exists
	if cart.CartID == uuid.Nil {
		return OrderResponse{}, fmt.Errorf("CartID not found for user %s", userID)
	}

	// Fetch cartItems based on product IDs in the request
	cartItems, err := c.CartRepository.ResolveCartItemsByCartID(cart.CartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	// Check if the cartItems list is empty
	if len(cartItems) == 0 {
		return OrderResponse{}, fmt.Errorf("No items in the cartItems list")
	}

	// Calculate totalAmount and items
	totalAmount, items, err := c.calculateTotalAndItems(cartItems)
	if err != nil {
		return OrderResponse{}, err
	}

	// Create order for the user
	order, err := c.createOrder(userID, totalAmount)
	if err != nil {
		return OrderResponse{}, err
	}

	// Check and update product stock and create order items
	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			return OrderResponse{}, err
		}

		if cartItem.Quantity > product.Stock {
			return OrderResponse{}, fmt.Errorf("product '%s' is out of stock", product.Name)
		}

		totalAmount += float64(cartItem.Quantity) * product.Price

		// Update product stock and create order item
		if err := c.processOrderItemsAndStock(order, []CartItems{cartItem}); err != nil {
			return OrderResponse{}, err
		}
	}

	// Clear the cart after successful checkout
	if err := c.CartRepository.ClearCart(cart.CartID); err != nil {
		return OrderResponse{}, err
	}

	// Build and return the order response
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
			ID:        cartItem.CartItemID,
			Quantity:  cartItem.Quantity,
			ProductID: cartItem.ProductID,
			CreatedAt: cartItem.CreatedAt,
			CreatedBy: cartItem.CreatedBy,
			Product: ProductDetails{
				ID:          product.ProductID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Stock:       product.Stock,
				CategoryID:  product.CategoryID,
				CreatedAt:   product.CreatedAt,
				CreatedBy:   product.CreatedBy,
			},
		}
		orderItemsInfo = append(orderItemsInfo, orderItemInfo)
	}

	return totalAmount, orderItemsInfo, nil
}
