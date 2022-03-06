package models

import (

	"github.com/jinzhu/gorm"
)


type TransactionStatus string

const(
	Success TransactionStatus = "success"
	Failed TransactionStatus = "failed"
)


type Transaction struct{
	gorm.Model
	Amount float64
	Currency Currency
	UserID uint `gorm:"many"`
	SendersID uint
	ReceiversID uint 
	Status TransactionStatus
}