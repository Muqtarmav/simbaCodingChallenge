package models

import (

	"github.com/jinzhu/gorm"
)



type Transaction struct{
	gorm.Model
	Amount float64
	Currency Currency
	UserID uint `gorm:"many"`
	Sender User
	Receiver User
}