package services

import (
	"github.com/djfemz/simbaCodingChallenge/test/services"
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
	//validate receiver using receiver's id
	//if receiver is a registered user that exists, perform transfer
	var transaction models.Transaction = models.Transaction{
		Amount:          transferRequest.Amount,
		Currency:        transferRequest.Currency,
		UserID:          transferRequest.UserID,
		ReceiversID:     transferRequest.RecipientsID,
		TransactionType: models.TRANSFER,
	}
	recipient, err := validateUser(transferRequest.RecipientsID)
	if err != nil {
		return dtos.TransactionResponse{Status: models.FAILED}
	}

	//retrieve recipients name
	log.Println("recipient-->", recipient.Name)
	//steps to perform transfer
	//find sender
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
	//if funds are sufficient, take funds from senders account and add it to receivers account
	//if funds are not sufficient, return an error
	//if receiver isn't a valid user, don't perform transfer
	transaction.Status = models.FAILED
	transactionRepo.Save(&transaction)
	return dtos.TransactionResponse{Status: models.FAILED}
} //if receiver is a valid user, perform transfer

//to transfer to target currency
//1. check senders balance for target currency
//2. if funds in senders target currency wallet is greater or equal to transfer amount
//3. perform transfer and exit
//4. if funds in senders target currency account are insufficient, check other currency balances
//5. if the balance in other currency accounts are sufficient,
//6. convert funds to target currency and transfer
//7. if balance in other currency accounts are insufficient, return insufficient balance

func isValidDeposit(recipient, sender *models.User, transferRequest dtos.TransactionRequest) bool {
	var sendersTargetCurrencyBalance float64

	for index, balance := range sender.Balance {
		//get senders total balance in target currency
		if balance.Currency == transferRequest.Currency {
			sendersTargetCurrencyBalance += balance.Amount

			if balance.Amount >= transferRequest.Amount {
				var sendersBalanceBeforeDeposit = sender.Balance[index].Amount
				var recipientsBalanceBeforeDeposit = recipient.Balance[index].Amount
				transfer(index, recipient, sender, transferRequest)
				userRepo.UpdateUserDetails(sender)
				userRepo.UpdateUserDetails(recipient)
				if sendersBalanceBeforeDeposit == sender.Balance[index].Amount ||
					recipientsBalanceBeforeDeposit == recipient.Balance[index].Amount {
					return false
				}
				return true
			}
		}
		if balance.Currency != transferRequest.Currency {
			//if balance is less than 1 don't pull any cash from particular balance
			if balance.Amount < 1.0 {
				continue
			}
			rate := services.GetCurrencyExchangeRate(string(balance.Currency), string(transferRequest.Currency))
			log.Println("rate---->", rate)
			convertedCash := convertCurrency(rate, balance.Amount)
			log.Println("amount to convert---->", balance.Amount)
			log.Println("converted cash---->", convertedCash)
			sendersTargetCurrencyBalance = sendersTargetCurrencyBalance + convertedCash
			if sendersTargetCurrencyBalance < transferRequest.Amount {
				continue
			} else {
				var sendersBalanceBeforeDeposit = sender.Balance[index].Amount
				var recipientsBalanceBeforeDeposit = recipient.Balance[index].Amount
				transfer(index, recipient, sender, transferRequest)
				userRepo.UpdateUserDetails(sender)
				userRepo.UpdateUserDetails(recipient)
				if sendersBalanceBeforeDeposit == sender.Balance[index].Amount ||
					recipientsBalanceBeforeDeposit == recipient.Balance[index].Amount {
					return false
				}
				return true
			}

		}
	}
	return false
}

func convertCurrency(rate, amount float64) float64 {
	return rate * amount
}

func transfer(index int, recipient, sender *models.User, transferRequest dtos.TransactionRequest) {
	sender.Balance[index].Amount -= transferRequest.Amount
	recipient.Balance[index].Amount += transferRequest.Amount
}
