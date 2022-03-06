package services

import "github.com/djfemz/simbaCodingChallenge/dtos"


type TransactionService interface{
	Transfer(transferRequest dtos.TransferRequest) dtos.TransferResponse
}

type TransactionServiceImpl struct{

}


func (transactionService TransactionServiceImpl)  Transfer(transferRequest dtos.TransferRequest) dtos.TransferResponse {
	return dtos.TransferResponse{}
}