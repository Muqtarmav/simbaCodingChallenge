package dtos

import (
	"github.com/djfemz/simbaCodingChallenge/data"
)

type TransactionResponse struct {
	UserID          uint
	ReceiversID     uint
	Amount          float64
	Balance         []data.Money
	Transactions    []data.Transaction
	Image           string
	TransactionType data.TransactionType
	Status          data.TransactionStatus
}
