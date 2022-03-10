package services

import (
	"github.com/jinzhu/gorm"
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/dtos"
)

var (
	transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}
	rate            float64
)

type TransactionService interface {
	Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse
	ConvertMoney(request dtos.TransactionRequest) dtos.TransactionResponse
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

func (transactionService TransactionServiceImpl) ConvertMoney(request dtos.TransactionRequest) dtos.TransactionResponse {
	transaction := &models.Transaction{
		Amount:          request.Amount,
		SourceCurrency:  request.SourceCurrency,
		TargetCurrency:  request.TargetCurrency,
		TransactionType: models.CONVERT,
	}

	var convertedAmount float64
	if request.Amount < 1 {
		transaction.Status = models.FAILED
		transactionRepo.Save(transaction)
		return dtos.TransactionResponse{Status: models.FAILED}
	}
	foundUser := userRepo.FindById(request.UserID)
	for index, balance := range foundUser.Balance {
		if balance.Currency == request.SourceCurrency {
			rate := GetCurrencyExchangeRate(string(request.SourceCurrency), string(request.TargetCurrency))
			log.Println("rate---->", rate)
			convertedAmount = convertCurrency(rate, request.Amount)

			if convertedAmount > 0.0 {
				foundUser.Balance[index].Amount -= request.Amount
			}
		}
	}
	for index, balance := range foundUser.Balance {
		if balance.Currency == request.TargetCurrency {
			foundUser.Balance[index].Amount += convertedAmount
			log.Println("converted amount---->", convertedAmount)
			log.Println("new balance---->", balance.Amount)
			log.Println("found user---->", foundUser)
			transaction.Status = models.SUCCESS
			transactionRepo.Save(transaction)
			userRepo.UpdateUserDetails(foundUser)
			return dtos.TransactionResponse{UserID: foundUser.ID, TransactionType: models.CONVERT, Status: models.SUCCESS}
		}
	}
	transaction.Status = models.FAILED
	transactionRepo.Save(transaction)
	return dtos.TransactionResponse{Status: models.FAILED}
}

func (transactionService TransactionServiceImpl) Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse {
	setUp()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(Db)
	var sender *models.User

	var transaction = models.Transaction{
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

	if transferRequest.UserID == 123456789 {
		log.Println("----> admin transfer")
		adminTransfer(recipient, transferRequest)
		transaction.Status = models.SUCCESS
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{
			UserID:          transferRequest.UserID,
			ReceiversID:     transferRequest.RecipientsID,
			Amount:          transferRequest.Amount,
			TransactionType: transferRequest.TransactionType,
			Status:          models.SUCCESS,
		}
	}

	sender, err = validateUser(transferRequest.UserID)
	if err != nil {
		transaction.Status = models.FAILED
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{Status: models.FAILED}
	} else {
		if isValidDeposit(recipient, sender, transferRequest) {
			transaction.Status = models.SUCCESS
			transaction.ExchangeRate = rate
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
	for index, balance := range sender.Balance {
		if balance.Currency == transferRequest.SourceCurrency &&
			balance.Amount >= transferRequest.Amount {
			transfer(index, recipient, sender, transferRequest)

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
			rate = GetCurrencyExchangeRate(string(transferRequest.SourceCurrency), string(transferRequest.TargetCurrency))
			amount := convertCurrency(rate, transferRequest.Amount)
			recipient.Balance[index].Amount += amount
		}
	}
}

func adminTransfer(recipient *models.User, transferRequest dtos.TransactionRequest) {
	for index, balance := range recipient.Balance {
		if balance.Currency == transferRequest.TargetCurrency {
			recipient.Balance[index].Amount += transferRequest.Amount
			userRepo.UpdateUserDetails(recipient)
		}
	}
}
