package dtos

import (
	"github.com/djfemz/simbaCodingChallenge/data"
)

type TransactionRequest struct {
	Amount          float64
	SourceCurrency  data.Currency
	TargetCurrency  data.Currency
	UserID          uint
	RecipientsID    uint
	TransactionType data.TransactionType
}
