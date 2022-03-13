package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"

type TransactionResponse struct {
	UserID          uint
	ReceiversID     uint
	Amount          float64
	Balance         []models.Money
	Transactions    []models.Transaction
	Image           string
	TransactionType models.TransactionType
	Status          models.TransactionStatus
}
