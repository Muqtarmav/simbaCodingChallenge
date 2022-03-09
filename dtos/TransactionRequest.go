package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"

type TransactionRequest struct {
	Amount          float64
	SourceCurrency  models.Currency
	TargetCurrency  models.Currency
	UserID          uint
	RecipientsID    uint
	TransactionType models.TransactionType
}
