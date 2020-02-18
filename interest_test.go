package main

import (
	"testing"
)

func TestValidate(t *testing.T) {

	cr := &CalculateRequest{}
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}

	cr = &CalculateRequest{
		PropertyPrice:   1,
		DownPayment:     1,
		AnnualInterest:  1,
		Amortization:    5,
		PaymentSchedule: 2,
	}
	if err := cr.Validate(); err != nil {
		t.Fatal("should be valid")
	}
	// invalid PaymentSchedule
	cr.PaymentSchedule = 1
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	cr.PaymentSchedule = 3
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	cr.PaymentSchedule = 5
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	// fix back
	cr.PaymentSchedule = 2

	// invalid Amortization
	cr.Amortization = 1
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	cr.Amortization = 4
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	cr.Amortization = 6
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
	cr.Amortization = 31
	if err := cr.Validate(); err == nil {
		t.Fatal("should be invalid")
	}
}
