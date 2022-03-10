package models

import (
	"github.com/jinzhu/gorm"
)

type TransactionStatus string
type TransactionType string

const (
	SUCCESS  TransactionStatus = "success"
	FAILED   TransactionStatus = "failed"
	TRANSFER TransactionType   = "transfer"
	CONVERT  TransactionType   = "conversion"
)

type Transaction struct {
	gorm.Model
	Amount          float64
	SourceCurrency  Currency
	TargetCurrency  Currency
	ExchangeRate    float64
	UserID          uint
	ReceiversID     uint
	TransactionType TransactionType
	Status          TransactionStatus
}
