package services

import (

	"github.com/djfemz/simbaCodingChallenge/dtos"
)


type TransactionService interface{
	Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse
}


type TransactionServiceImpl struct{

}


func (transactionService TransactionServiceImpl)  Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse {
	setUp()
	defer Db.Close()
	//validate receiver using receiver's id
	if isValidUser(transferRequest.ReceiversID){
		//if receiver is a valid user that exists, perform transfer
		//steps to perform transfer
		//find sender
		//check senders account for funds
		//if funds are sufficient, take funds from senders account and add it to receivers account
		//if funds are not sufficient, return an error
	}
	
	//if receiver isn't a valid user, don't perform transfer
	return dtos.TransactionResponse{}
}