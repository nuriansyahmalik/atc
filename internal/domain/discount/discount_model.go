package discount

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Discount struct {
	ID        uuid.UUID   `db:"discount_id"`
	Code      string      `db:"code"`
	Type      string      `db:"type"`
	Price     float64     `db:"price"`
	StartDate time.Time   `db:"start_date"`
	EndDate   time.Time   `db:"end_date"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy uuid.UUID   `db:"created_by"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

type (
	DiscountRequestFormat struct {
		ID        uuid.UUID `json:"ID"`
		Code      string    `json:"code"`
		Type      string    `json:"type"`
		Price     float64   `json:"value"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}
	DiscountResponseFormat struct {
		ID        uuid.UUID `json:"ID"`
		Code      string    `json:"code"`
		Type      string    `json:"type"`
		Price     float64   `json:"value"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}
)

func (d Discount) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.ToResponseFormat())
}

func (d *Discount) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(d)
}

func (d Discount) DiscountRequestFormat(req DiscountRequestFormat, userID uuid.UUID) (discount Discount, err error) {
	discountID, err := uuid.NewV4()
	if err != nil {
		return
	}
	discount = Discount{
		ID:        discountID,
		Code:      req.Code,
		Type:      req.Type,
		Price:     req.Price,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}
	return
}

func (d Discount) ToResponseFormat() DiscountResponseFormat {
	return DiscountResponseFormat{
		ID:        d.ID,
		Code:      d.Code,
		Type:      d.Type,
		Price:     d.Price,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
	}
}
