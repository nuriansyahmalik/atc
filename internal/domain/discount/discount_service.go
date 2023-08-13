package discount

import (
	"fmt"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type DiscountService interface {
	ResolveByID(discountID string) (Discount, error)
	CreateDiscount(requestFormat DiscountRequestFormat, userID uuid.UUID) (discount Discount, err error)
	ApplyDiscount(totalAmount float64, discountCode string) (float64, float64, error)
}

type DiscountServiceImpl struct {
	DiscountRepository DiscountRepository
	Config             *configs.Config
}

func ProvideDiscountServiceImpl(discountRepository DiscountRepository, config *configs.Config) *DiscountServiceImpl {
	return &DiscountServiceImpl{DiscountRepository: discountRepository, Config: config}
}

func (d *DiscountServiceImpl) ResolveByID(discountID string) (Discount, error) {
	discount, err := d.DiscountRepository.ResolveByID(discountID)
	if err != nil {
		return Discount{}, fmt.Errorf("failed to get discount: %w", err)
	}
	return discount, nil
}
func (d *DiscountServiceImpl) ApplyDiscount(totalAmount float64, discountCode string) (float64, float64, error) {
	if discountCode == "" {
		return totalAmount, 0, nil
	}

	discount, err := d.DiscountRepository.ResolveByCode(discountCode)
	if err != nil {
		return 0, 0, err
	}

	discountAmount := calculateDiscountAmount(totalAmount, discount)
	totalPriceAfterDiscount := totalAmount - discountAmount

	return totalPriceAfterDiscount, discountAmount, nil
}

func (d *DiscountServiceImpl) CreateDiscount(requestFormat DiscountRequestFormat, userID uuid.UUID) (discount Discount, err error) {
	discount, err = discount.DiscountRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}
	if err != nil {
		return discount, failure.BadRequest(err)
	}
	err = d.DiscountRepository.CreateDiscount(discount)
	if err != nil {
		return
	}
	return
}

func calculateDiscountAmount(totalAmount float64, discount Discount) float64 {
	if discount.Type == "percentage" {
		return totalAmount * discount.Price / 100
	} else if discount.Type == "fixed_amount" {
		return discount.Price
	}
	return 0
}
