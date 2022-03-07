package services

import (
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/dtos"
)


var transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}

type TransactionService interface{
	Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse
}




type TransactionServiceImpl struct{

}


func (transactionService TransactionServiceImpl)  Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse {
	// var amountToDeposit float64 = transferRequest.Amount
	setUp()
	defer Db.Close()
	var sender *models.User
	//validate receiver using receiver's id
	//if receiver is a registered user that exists, perform transfer
	var transaction models.Transaction= models.Transaction{
		Amount: transferRequest.Amount, 
		Currency: transferRequest.Currency,
		UserID: transferRequest.UserID,
		ReceiversID: transferRequest.RecipientsID,
		TransactionType: models.TRANSFER,
	}
	recipient, err:=validateUser(transferRequest.RecipientsID)
	if err!=nil{
		return dtos.TransactionResponse{Status: models.FAILED}
	}
	
		//retrieve recipients name 
		log.Println("recipient-->", recipient.Name)
		//steps to perform transfer
			//find sender
	sender, err = validateUser(transferRequest.UserID)

	if err!=nil{
		transaction.Status = models.FAILED
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{Status: models.FAILED}
	}else {
		//check senders account for funds
		if isDeposit(recipient, sender, transferRequest){
			transaction.Status=models.SUCCESS
			transactionRepo.Save(&transaction)
			return dtos.TransactionResponse{
				UserID: sender.ID, 
				ReceiversID: recipient.ID, 
				Amount: transferRequest.Amount,
				TransactionType: transferRequest.TransactionType,
				Status: models.SUCCESS,
			}
		}
	}	
		//if funds are sufficient, take funds from senders account and add it to receivers account
		//if funds are not sufficient, return an error
	//if receiver isn't a valid user, don't perform transfer
	transaction.Status = models.FAILED
	transactionRepo.Save(&transaction)
	return dtos.TransactionResponse{Status: models.FAILED}
}//if receiver is a valid user, perform transfer


func isDeposit(recipient, sender *models.User, transferRequest dtos.TransactionRequest)  bool{
	for index, balance := range sender.Balance  {
		if balance.Currency==transferRequest.Currency&&
			balance.Amount>=transferRequest.Amount{
			var sendersBalanceBeforeDeposit float64 = sender.Balance[index].Amount
			var recipientsBalanceBeforeDeposit float64 = recipient.Balance[index].Amount
			log.Println("senders balance before-->", sender.Balance)
			log.Println("recipients balance before -->", recipient.Balance)
			sender.Balance[index].Amount -=transferRequest.Amount
			recipient.Balance[index].Amount+=transferRequest.Amount
			log.Println("senders balance after-->", sender.Balance)
			log.Println("recipients balance after", recipient.Balance)
			userRepo.UpdateUserDetails(sender)
			userRepo.UpdateUserDetails(recipient)
			if sendersBalanceBeforeDeposit==sender.Balance[index].Amount || 
						recipientsBalanceBeforeDeposit==recipient.Balance[index].Amount {
				return false
			}
				return true
			}
	}
	return false
}