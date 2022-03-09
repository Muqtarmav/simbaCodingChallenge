package services

import (
	"github.com/jinzhu/gorm"
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/dtos"
)

var transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}

type TransactionService interface {
	Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse
}

//TODOS
// 1. whenever a transaction fails, it is visible to the receiving person as a message in his transaction history,
// which also resembles his bank account.
//2. The balance of the bank account is not hardcoded but derived from the entire transaction history.
//3. The exchange rates can be pulled from a public API
//4. When a user successfully logs in, he sees a page with all of his transactions,
// including the initial transaction from the signup (1000 USD).
//5. The page with all the transactions shows also the current balance for each currency. (e.g. start: 1000 USD, 0 EUR, 0 NGN)
//6. A transaction consists of the sender and the receiver, the source currency, target currency, exchange rate, and the amount.
//7. A user can select the target currency.
//8. Check if a transaction is possible, so if a user has enough funds.

type TransactionServiceImpl struct {
}

func (transactionService TransactionServiceImpl) Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse {
	// var amountToDeposit float64 = transferRequest.Amount
	setUp()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(Db)
	var sender *models.User

	var transaction models.Transaction = models.Transaction{
		Amount:          transferRequest.Amount,
		SourceCurrency:  transferRequest.SourceCurrency,
		TargetCurrency:  transferRequest.TargetCurrency,
		UserID:          transferRequest.UserID,
		ReceiversID:     transferRequest.RecipientsID,
		TransactionType: models.TRANSFER,
	}
	recipient, err := validateUser(transferRequest.RecipientsID)
	if err != nil {
		return dtos.TransactionResponse{Status: models.FAILED}
	}

	sender, err = validateUser(transferRequest.UserID)
	if err != nil {
		transaction.Status = models.FAILED
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{Status: models.FAILED}
	} else {
		//check senders account for funds
		if isValidDeposit(recipient, sender, transferRequest) {
			transaction.Status = models.SUCCESS
			transactionRepo.Save(&transaction)
			return dtos.TransactionResponse{
				UserID:          sender.ID,
				ReceiversID:     recipient.ID,
				Amount:          transferRequest.Amount,
				TransactionType: transferRequest.TransactionType,
				Status:          models.SUCCESS,
			}
		}
	}
	transaction.Status = models.FAILED
	transactionRepo.Save(&transaction)
	return dtos.TransactionResponse{Status: models.FAILED}
}

func isValidDeposit(recipient, sender *models.User, transferRequest dtos.TransactionRequest) bool {
	log.Println(transferRequest.SourceCurrency)
	for index, balance := range sender.Balance {
		if balance.Currency == transferRequest.SourceCurrency && balance.Amount >= transferRequest.Amount {
			transfer(index, recipient, sender, transferRequest)
			log.Println("sender after transfer--->", sender)
			log.Println("recipient after transfer--->", recipient)

			userRepo.UpdateUserDetails(sender)
			userRepo.UpdateUserDetails(recipient)
			return true
		}
	}
	return false
}

func convertCurrency(rate, amount float64) float64 {
	return rate * amount
}

func transfer(index int, recipient, sender *models.User, transferRequest dtos.TransactionRequest) {
	sender.Balance[index].Amount -= transferRequest.Amount
	for index, balance := range recipient.Balance {
		if balance.Currency == transferRequest.TargetCurrency {
			rate := GetCurrencyExchangeRate(string(transferRequest.SourceCurrency), string(transferRequest.TargetCurrency))
			log.Println("rate", rate)
			amount := convertCurrency(rate, transferRequest.Amount)
			recipient.Balance[index].Amount += amount
			log.Println("after adding---->", recipient.Balance[index].Amount)
		}
	}
}
