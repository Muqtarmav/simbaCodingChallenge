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
		Amount:          200,
		Currency:        models.DOLLAR,
		UserID:          65,
		RecipientsID:    uint(62),
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
		Currency:        models.EURO,
		UserID:          uint(62),
		RecipientsID:    uint(65),
		TransactionType: models.TRANSFER,
	}
	log.Println(transferRequest.UserID)
	sender := userRepo.FindById(transferRequest.UserID)
	log.Println(sender)
	assert.Equal(t, 4000.00, sender.Balance[2].Amount)

	//fetch recipients account balance
	recipient := userRepo.FindById(transferRequest.RecipientsID)
	assert.Equal(t, 400.00, recipient.Balance[2].Amount)

	response := transactionService.Deposit(transferRequest)
	log.Println("response-->", response)
	assert.Equal(t, 4000.00, sender.Balance[2].Amount)
	assert.Equal(t, 400.00, recipient.Balance[2].Amount)

}

func TestThatUserCanSelectTargetCurrencyDuringTransfer(t *testing.T) {

}
