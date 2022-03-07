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
		Amount: 200,
		Currency: models.DOLLAR,
		UserID: 63,
		RecipientsID: uint(62),
		TransactionType: models.TRANSFER,
	}

	transferResponse:= transactionService.Deposit(transferRequest)
	assert.NotEmpty(t, transferResponse)
	log.Println("transfer respone -->", transferResponse)
	assert.Equal(t, models.SUCCESS, transferResponse.Status)
}


