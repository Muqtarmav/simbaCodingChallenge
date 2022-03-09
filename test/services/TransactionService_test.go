package services_test

import (
	"log"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"github.com/stretchr/testify/assert"
)

var transactionService services.TransactionService = services.TransactionServiceImpl{}

func TestThatAUserCanTransferVirtual_CashToOtherUsers(t *testing.T) {
	var transferRequest = dtos.TransactionRequest{
		Amount:          50,
		SourceCurrency:  models.DOLLAR,
		TargetCurrency:  models.DOLLAR,
		UserID:          2,
		RecipientsID:    uint(3),
		TransactionType: models.TRANSFER,
	}

	transferResponse := transactionService.Deposit(transferRequest)
	assert.NotEmpty(t, transferResponse)
	log.Println("transfer response -->", transferResponse)
	assert.Equal(t, models.SUCCESS, transferResponse.Status)
}

func TestThatTransferFailsWhenUserHasInsufficientFunds(t *testing.T) {
	//get a users account balance
	var transferRequest = dtos.TransactionRequest{
		Amount:          40000,
		SourceCurrency:  models.DOLLAR,
		TargetCurrency:  models.DOLLAR,
		UserID:          uint(2),
		RecipientsID:    uint(3),
		TransactionType: models.TRANSFER,
	}
	sender := userRepo.FindById(transferRequest.UserID)
	assert.NotEmpty(t, sender)
	log.Println(sender)
	assert.Equal(t, 300.00, sender.Balance[2].Amount)

	//fetch recipients account balance
	recipient := userRepo.FindById(transferRequest.RecipientsID)
	assert.Equal(t, 950.00, recipient.Balance[2].Amount)

	response := transactionService.Deposit(transferRequest)
	log.Println("response-->", response)
	assert.Equal(t, 300.00, sender.Balance[2].Amount)
	assert.Equal(t, 950.00, recipient.Balance[2].Amount)

}

func TestThatUserCanSelectTargetCurrencyDuringTransfer(t *testing.T) {
	var transferRequest = dtos.TransactionRequest{
		Amount:          50,
		SourceCurrency:  models.DOLLAR,
		TargetCurrency:  models.EURO,
		UserID:          uint(3),
		RecipientsID:    uint(2),
		TransactionType: models.TRANSFER,
	}
	response := transactionService.Deposit(transferRequest)
	log.Println(response)
	assert.NotEmpty(t, response)
}
