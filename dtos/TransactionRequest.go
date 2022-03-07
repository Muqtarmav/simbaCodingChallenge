package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"

type TransactionRequest struct{
	Amount float64
	Currency models.Currency
	UserID uint 
	RecipientsID uint 
	TransactionType models.TransactionType
}