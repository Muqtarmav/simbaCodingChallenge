package services

import (
	"github.com/djfemz/simbaCodingChallenge/data"
	"github.com/jinzhu/gorm"
	"log"

	"github.com/djfemz/simbaCodingChallenge/dtos"
)

var (
	transactionRepo data.TransactionRepository = &data.TransactionRepositoryImpl{}
	rate            float64
)

type TransactionService interface {
	Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse
	ConvertMoney(request dtos.TransactionRequest) dtos.TransactionResponse
	GetUsersTransactions(userID uint) []*data.Transaction
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
	transaction := &data.Transaction{
		Amount:          request.Amount,
		SourceCurrency:  request.SourceCurrency,
		TargetCurrency:  request.TargetCurrency,
		TransactionType: data.CONVERT,
	}

	var convertedAmount float64
	if request.Amount < 1 {
		transaction.Status = data.FAILED
		transactionRepo.Save(transaction)
		return dtos.TransactionResponse{Status: data.FAILED}
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
			transaction.Status = data.SUCCESS
			transactionRepo.Save(transaction)
			userRepo.UpdateUserDetails(foundUser)
			userAfterCurrencyConversion := userRepo.FindById(foundUser.ID)
			return dtos.TransactionResponse{
				UserID:          foundUser.ID,
				TransactionType: data.CONVERT,
				Balance:         userAfterCurrencyConversion.Balance,
				Transactions:    userAfterCurrencyConversion.Transactions,
				Status:          data.SUCCESS,
			}
		}
	}
	transaction.Status = data.FAILED
	transactionRepo.Save(transaction)
	return dtos.TransactionResponse{Status: data.FAILED}
}

func (transactionService TransactionServiceImpl) Deposit(transferRequest dtos.TransactionRequest) dtos.TransactionResponse {
	setUp()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(Db)
	var sender *data.User

	var transaction = data.Transaction{
		Amount:          transferRequest.Amount,
		SourceCurrency:  transferRequest.SourceCurrency,
		TargetCurrency:  transferRequest.TargetCurrency,
		UserID:          transferRequest.UserID,
		ReceiversID:     transferRequest.RecipientsID,
		TransactionType: data.TRANSFER,
	}
	recipient, err := validateUser(transferRequest.RecipientsID)
	if err != nil {
		return dtos.TransactionResponse{Status: data.FAILED}
	}

	if transferRequest.UserID == 123456789 {
		log.Println("----> admin transfer")
		adminTransfer(recipient, transferRequest)
		transaction.Status = data.SUCCESS
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{
			UserID:          transferRequest.UserID,
			ReceiversID:     transferRequest.RecipientsID,
			Amount:          transferRequest.Amount,
			TransactionType: transferRequest.TransactionType,
			Status:          data.SUCCESS,
		}
	}

	sender, err = validateUser(transferRequest.UserID)
	if err != nil {
		transaction.Status = data.FAILED
		transactionRepo.Save(&transaction)
		return dtos.TransactionResponse{Status: data.FAILED}
	} else {
		if isValidDeposit(recipient, sender, transferRequest) {
			transaction.Status = data.SUCCESS
			transaction.ExchangeRate = rate
			transactionRepo.Save(&transaction)
			senderAfterSuccessfulTransfer := userRepo.FindById(sender.ID)
			return dtos.TransactionResponse{
				UserID:          sender.ID,
				ReceiversID:     recipient.ID,
				Amount:          transferRequest.Amount,
				Balance:         senderAfterSuccessfulTransfer.Balance,
				Transactions:    senderAfterSuccessfulTransfer.Transactions,
				TransactionType: transferRequest.TransactionType,
				Status:          data.SUCCESS,
			}
		}
	}
	transaction.Status = data.FAILED
	transactionRepo.Save(&transaction)
	return dtos.TransactionResponse{Status: data.FAILED}
}

func isValidDeposit(recipient, sender *data.User, transferRequest dtos.TransactionRequest) bool {
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

func transfer(index int, recipient, sender *data.User, transferRequest dtos.TransactionRequest) {
	sender.Balance[index].Amount -= transferRequest.Amount
	for index, balance := range recipient.Balance {
		if balance.Currency == transferRequest.TargetCurrency {
			rate = GetCurrencyExchangeRate(string(transferRequest.SourceCurrency), string(transferRequest.TargetCurrency))
			amount := convertCurrency(rate, transferRequest.Amount)
			recipient.Balance[index].Amount += amount
		}
	}
}

func adminTransfer(recipient *data.User, transferRequest dtos.TransactionRequest) {
	for index, balance := range recipient.Balance {
		if balance.Currency == transferRequest.TargetCurrency {
			recipient.Balance[index].Amount += transferRequest.Amount
			userRepo.UpdateUserDetails(recipient)
		}
	}
}

func (transactionService TransactionServiceImpl) GetUsersTransactions(userID uint) []*data.Transaction {
	var usersTransactions []*data.Transaction
	transactions := transactionRepo.FindAllTransactions()
	for _, transaction := range transactions {
		if transaction.UserID == userID {
			log.Println("each transaction---->", transaction)
			usersTransactions = append(usersTransactions, transaction)
		}
	}
	log.Println(usersTransactions)
	return usersTransactions
}
