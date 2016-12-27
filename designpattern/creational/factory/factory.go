package creational

import (
	"errors"
	"fmt"
)

type PaymentMethod interface {
	Pay(amount float32) string
}

type CashPM struct{}

type DebitCardPM struct{}

type NewDebitCardPM struct{}

const (
	Cash         = 1
	DebitCard    = 2
	NewDebitCard = 3
)

func GetPaymentMethod(m int) (PaymentMethod, error) {
	switch m {
	case Cash:
		return new(CashPM), nil
	case DebitCard:
		return new(DebitCardPM), nil
	case NewDebitCard:
		return new(NewDebitCardPM), nil
	default:
		return nil, errors.New(fmt.Sprintf("Payment method %d not recognized\n", m))
	}

}

func (c *CashPM) Pay(amount float32) string {
	return fmt.Sprintf("%0.2f paid using cash\n", amount)
}

func (c *DebitCardPM) Pay(amount float32) string {
	return fmt.Sprintf("%#0.2f paid using debit card\n", amount)
}

func (d *NewDebitCardPM) Pay(amount float32) string {
	return fmt.Sprintf("%#0.2f paid using new debit card implementation\n", amount)
}
