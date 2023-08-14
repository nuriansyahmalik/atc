package cart_test

import (
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"testing"

	cart_mock "github.com/evermos/boilerplate-go/internal/domain/cart/mock"
	product_mock "github.com/evermos/boilerplate-go/internal/domain/product/mock"
)

func getRandomUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func uuidFromString(s string) uuid.UUID {
	id, _ := uuid.FromString(s)
	return id
}

func TestCartService(t *testing.T) {
	t.Run("ResolveCartByID", func(t *testing.T) {
		tests := []struct {
			name            string
			cartID          uuid.UUID
			userID          uuid.UUID
			setupMock       func(*cart_mock.MockCartRepository, *product_mock.MockProductRepository, uuid.UUID, cart.Cart, []cart.CartItems, error)
			returns         cart.Cart
			returnCartItems []cart.CartItems
			err             error
			returnProduct   product.Product
		}{
			{
				name:   "Default",
				cartID: getRandomUUID(),
				setupMock: func(mockCartRepo *cart_mock.MockCartRepository, mockProductRepo *product_mock.MockProductRepository, id uuid.UUID, cart cart.Cart, cartItems []cart.CartItems, err error) {
					mockCartRepo.EXPECT().ResolveCartByID(id).Return(cart, err)
					mockCartRepo.EXPECT().ResolveCartItemsByCartID(id).Return(cartItems, err)

					mockProductRepo.EXPECT().ResolveByID(gomock.Any()).Return(product.Product{}, nil)
				},
				returns: cart.Cart{
					CartID: getRandomUUID(),
					UserID: getRandomUUID(),
				},
				returnCartItems: []cart.CartItems{
					{
						CartItemID: getRandomUUID(),
						CartID:     getRandomUUID(),
						ProductID:  getRandomUUID(),
					},
				},
				returnProduct: product.Product{
					ProductID:   getRandomUUID(),
					CategoryID:  getRandomUUID(),
					Name:        "Iphone XX",
					Description: "The New Iphone",
					Price:       float64(65000),
					Stock:       float64(650),
				},
				err: nil,
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		//for _, test := range tests {
		//	t.Run(test.name, func(t *testing.T) {
		//		mockCartRepo := cart_mock.NewMockCartRepository(ctrl)
		//		mockProductRepo := product_mock.NewMockProductRepository(ctrl)
		//		//service := cart.ProvideCarServiceImpl(mockCartRepo, mockProductRepo, nil)
		//		test.setupMock(mockCartRepo, mockProductRepo, test.cartID, test.returns, test.returnCartItems, test.err)
		//		//_, err := service.ResolveCartByID(test.cartID, test.userID)
		//
		//		assert.Equal(t, test.err, err)
		//		assert.Equal(t, test.returnProduct.Price, nil)
		//		// assert.Equal(t, len(test.returnCartItems), len(gotCartItems))
		//	})
		//}
	})
}
