package cart

//go:generate go run github.com/golang/mock/mockgen -source cart_service.go -destination mock/cart_service_mock.go -package cart_mock

import (
	"database/sql"
	"fmt"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/discount"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type CartService interface {
	CheckoutCarts(requestFormat CheckoutRequestFormat, userID uuid.UUID) (orders OrderResponse, err error)
	AddItemToCart(requestFormat AddToCartRequestFormat, userID uuid.UUID) (cart Cart, err error)
	ResolveCartByID(cartID uuid.UUID, userID uuid.UUID) (cart Cart, err error)
}

type CartServiceImpl struct {
	CartRepository     CartRepository
	ProductRepository  product.ProductRepository
	DiscountRepository discount.DiscountRepository
	Config             *configs.Config
}

func ProvideCarServiceImpl(cartRepository CartRepository, productRepository product.ProductRepository, discountRepository discount.DiscountRepository, config *configs.Config) *CartServiceImpl {
	return &CartServiceImpl{CartRepository: cartRepository, ProductRepository: productRepository, DiscountRepository: discountRepository, Config: config}
}

func (c *CartServiceImpl) AddItemToCart(req AddToCartRequestFormat, userID uuid.UUID) (cart Cart, err error) {
	product, err := c.ProductRepository.ResolveByID(req.ProductID)
	if err != nil {
		return
	}

	if product.Stock < req.Quantity {
		log.Info().Msg("Insufficient stock")
		return cart, failure.InternalError(err)
	}

	cart, err = c.getOrCreateCart(userID)
	if err != nil {
		fmt.Println(cart.UserID)
		log.Info().Msg("error disini 1")
		return
	}

	existingItem, err := c.CartRepository.ResolveCartItemByProduct(cart.CartID, req.ProductID)
	if err != nil {
		return
	}

	if existingItem == nil {
		err = c.createCartItem(cart.CartID, userID, req.ProductID, req.Quantity)
	} else {
		existingItem[0].Quantity += req.Quantity
		existingItem[0].CreatedAt = time.Now()
		err = c.CartRepository.UpdateCartItem(existingItem[0])
	}

	if err != nil {
		return
	}

	items, err := c.CartRepository.ResolveCartItemsByCartID(cart.CartID)
	if err != nil {
		return
	}
	cart.AttachItems(items)
	return
}
func (c *CartServiceImpl) ResolveCartByID(cartID uuid.UUID, userID uuid.UUID) (cart Cart, err error) {
	cart, err = c.CartRepository.ResolveCartByID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return cart, err
	}
	cartItems, err := c.CartRepository.ResolveCartItemsByCartID(cartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return cart, err
	}

	cart.AttachItems(cartItems)
	return
}

func (c *CartServiceImpl) CheckoutCarts(requestFormat CheckoutRequestFormat, userID uuid.UUID) (OrderResponse, error) {
	cart, err := c.CartRepository.ResolveCartByID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	if cart.CartID == uuid.Nil {
		return OrderResponse{}, fmt.Errorf("CartID not found for user %s", userID)
	}

	cartItems, err := c.CartRepository.ResolveCartItemsByCartID(cart.CartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return OrderResponse{}, err
	}

	if len(cartItems) == 0 {
		return OrderResponse{}, fmt.Errorf("No items in the cartItems list")
	}

	totalAmount, items, err := c.calculateTotalAndItems(cartItems)
	if err != nil {
		return OrderResponse{}, err
	}

	var discountAmount float64
	var discountID uuid.UUID
	if requestFormat.DiscountCode != "" {
		discount, err := c.DiscountRepository.ResolveByCode(requestFormat.DiscountCode)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Info().Msgf("Discount with code '%s' not found", requestFormat.DiscountCode)
			} else {
				logger.ErrorWithStack(err)
				return OrderResponse{}, fmt.Errorf("error retrieving discount with code '%s': %w", requestFormat.DiscountCode, err)
			}
		} else {
			discountAmount = calculateDiscountAmount(totalAmount, discount)
			discountID = discount.ID
		}
	}

	totalPriceAfterDiscount := totalAmount - discountAmount

	order, err := c.createOrder(userID, totalAmount)
	if err != nil {
		return OrderResponse{}, err
	}

	for _, cartItem := range cartItems {
		product, err := c.ProductRepository.ResolveByID(cartItem.ProductID)
		if err != nil {
			return OrderResponse{}, err
		}

		if cartItem.Quantity > product.Stock {
			return OrderResponse{}, fmt.Errorf("product '%s' is out of stock", product.Name)
		}

		totalAmount += float64(cartItem.Quantity) * product.Price

		if err := c.processOrderItemsAndStock(order, []CartItems{cartItem}); err != nil {
			return OrderResponse{}, err
		}
	}

	if err := c.CartRepository.ClearCart(cart.CartID); err != nil {
		return OrderResponse{}, err
	}

	orderResponse := order.BuildOrderResponse(order, items, discountAmount)
	orderResponse.TotalPrice = totalPriceAfterDiscount
	orderResponse.DiscountID = discountID
	return orderResponse, nil
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
func (c *CartServiceImpl) getOrCreateCart(userID uuid.UUID) (cart Cart, err error) {
	cart, err = c.CartRepository.ResolveCartByID(userID)
	if err == sql.ErrNoRows {
		cartID, err := uuid.NewV4()
		if err != nil {
			log.Info().Msg("error disni 2")
			return cart, err
		}
		if err := c.CartRepository.CreateCart(Cart{
			CartID:    cartID,
			UserID:    userID,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		}); err != nil {
			log.Info().Msg("error disini 3")
			return cart, err
		}
	} else if err != nil {
		log.Info().Msg("error disini 4")
		return
	}
	return
}
func (c *CartServiceImpl) createCartItem(cartID, userID, productID uuid.UUID, quantity float64) (err error) {
	cartItemID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return c.CartRepository.CreateCartItems(CartItems{
		CartItemID: cartItemID,
		CartID:     cartID,
		ProductID:  productID,
		Quantity:   quantity,
		CreatedAt:  time.Now(),
		CreatedBy:  userID,
	})
}
