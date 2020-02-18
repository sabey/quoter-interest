package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

// valid request
/*
clear ; curl -X POST -H "Content-Type: application/json" http://localhost
--data '{"property-price":1,"down-payment":1,"annual-interest":1,"amortization":5,"payment-schedule":4}' ; echo
*/
func main() {
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":80", nil))
}
func index(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method == "POST" {
		cr := &CalculateRequest{}
		err := json.NewDecoder(r.Body).Decode(cr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid JSON POST")
			return
		}
		pps, err := calculate(cr)
		if err != nil {
			bs, err := json.Marshal(&CaclulateResponse{
				Error: fmt.Sprintf("Invalid Calculation: \"%s\"", err.Error()),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				// this shouldn't be possible
				fmt.Fprintf(w, "failed to marshal json?!: %s\n", err)
				return
			}
			// marshaled error
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", bs)
			return
		}
		bs, err := json.Marshal(&CaclulateResponse{
			Result: pps,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// this shouldn't be possible
			fmt.Fprintf(w, "failed to marshal json?!: %s\n", err)
			return
		}
		// marshaled error
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", bs)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "unknown api endpoint, use POST /")
}
func calculate(
	cr *CalculateRequest,
) (
	float64,
	error,
) {
	// validate request
	if err := cr.Validate(); err != nil {
		// invalid request
		return 0.0, err
	}
	// r = per payment schedule interest rate, calculated by dividing the annual interest rate by the number of payments per annum
	scheduledInterestRate := cr.AnnualInterest / float64(cr.PaymentSchedule)
	// n = total number of payments over the amortization period
	numberOfPayments := float64(cr.PaymentSchedule) * float64(cr.Amortization)
	//
	// Mortgage Payment Formula
	// M = P ( (r(1+r)^n) / (((1+r)^n)-1) )
	//
	// r(1+r)^n
	a := scheduledInterestRate * math.Pow(1+scheduledInterestRate, numberOfPayments)
	// ((1+r)^n)-1
	b := math.Pow(1+scheduledInterestRate, numberOfPayments) - 1
	// m = p * a/b
	m := cr.PropertyPrice * (a / b)
	return m, nil
}

const (
	AmortizationYearsMin    = 5
	AmortizationYearsMod    = 5
	AmortizationYearsMax    = 30
	PaymentScheduleBiWeekly = 2
	PaymentScheduleMonthly  = 4
)

type CalculateRequest struct {
	PropertyPrice   float64 `json:"property-price"`
	DownPayment     float64 `json:"down-payment"`
	AnnualInterest  float64 `json:"annual-interest"`
	Amortization    int     `json:"amortization"`
	PaymentSchedule int     `json:"payment-schedule"`
}

func (self *CalculateRequest) Validate() error {
	if self.PropertyPrice <= 0 {
		return fmt.Errorf("Invalid property price")
	}
	if self.DownPayment <= 0 {
		return fmt.Errorf("Invalid down payment")
	}
	if self.AnnualInterest <= 0 {
		return fmt.Errorf("Invalid annual interest rate")
	}
	if self.Amortization < AmortizationYearsMin ||
		self.Amortization > AmortizationYearsMax {
		return fmt.Errorf("Invalid amortization, must be > 5 and < 30")
	}
	if self.Amortization%AmortizationYearsMod != 0 {
		// must be a mod of 5
		return fmt.Errorf("Invalid amortization, must be mod of %d", AmortizationYearsMod)
	}
	if self.PaymentSchedule != PaymentScheduleBiWeekly && // bi-weekly
		self.PaymentSchedule != PaymentScheduleMonthly { // monthly
		// invalid request
		return fmt.Errorf("Invalid payment schedule")
	}
	return nil
}

type CaclulateResponse struct {
	Error  string  `json:"error,omitempty"`
	Result float64 `json:"result,omitempty"`
}
