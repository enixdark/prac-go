package creational

import (
	"strings"
	"testing"
)

func TestCreatePaymentMethodCash(t *testing.T) {
	payment, err := GetPaymentMethod(Cash)
	// fmt.Println(err)
	if err != nil {
		t.Fatal("A payment method of type 'Cash' must exist")
	}

	msg := payment.Pay(10.30)

	if !strings.Contains(msg, "paid using cash") {
		t.Error("The cash payment method message wasn't correct")
	}

	t.Log("LOG:", msg)

	payment, err = GetPaymentMethod(DebitCard)

	if err != nil {
		t.Error("A payment method of type 'DebitCard' must exist")
	}

	msg = payment.Pay(20.30)

	if err == nil {
		t.Error("A payment method of type 'DebitCard' must exist")
	}

	t.Log("LOG:", err)
}
