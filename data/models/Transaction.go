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
)

type Transaction struct {
	gorm.Model
	Amount          float64
	SourceCurrency  Currency
	TargetCurrency  Currency
	UserID          uint
	ReceiversID     uint
	TransactionType TransactionType
	Status          TransactionStatus
}
