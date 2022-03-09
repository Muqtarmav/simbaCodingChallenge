package models

import "github.com/jinzhu/gorm"

type Currency string

const (
	DOLLAR Currency = "USD"
	EURO   Currency = "EUR"
	NAIRA  Currency = "NGN"
)

type Money struct {
	gorm.Model
	Amount   float64
	Currency Currency
	UserID   uint
}
