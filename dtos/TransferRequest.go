package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"


type TransferRequest struct{
	SendersID uint
	ReceiversID uint
	Amount float64
	Currency models.Currency

}