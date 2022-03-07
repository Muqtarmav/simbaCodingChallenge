package services_test

import (
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"github.com/stretchr/testify/assert"
)


var transactionService services.TransactionService = services.TransactionServiceImpl{}

func TestThatAUserCanTransferVirtual_CashToOtherUsers(t *testing.T) {
	var addUserDto = dtos.AddUserRequest{
		Name: "Jane Doe",
		Email: "jane@gmail.com",
		Password: "12345678",
	}

	addUserResponse:= userService.Register(addUserDto)
	assert.NotEmpty(t,addUserResponse)
	assert.Greater(t, addUserResponse.ID, 0)


	var transferRequest = dtos.TransactionRequest{
		Amount: 200,
		Currency: models.DOLLAR,
		UserID: addUserResponse.ID,
		ReceiversID: uint(62),
		TransactionType: models.TRANSFER,
	}

	transferResponse:= transactionService.Deposit(transferRequest)
	assert.NotEmpty(t, transferResponse)
	assert.Equal(t, models.SUCCESS, transferResponse.Status)

	
}


