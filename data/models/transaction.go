package models

import (

	"github.com/jinzhu/gorm"
)

const (
	DOLLAR = "USD"
	POUNDS = "GBP"
	NAIRA = "NGN"
) 

type Transaction struct{
	gorm.Model
	Amount float64
	Currency string
	UserID uint
	Sender User
	Receiver User
}