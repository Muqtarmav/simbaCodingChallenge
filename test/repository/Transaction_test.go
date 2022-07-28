package test

import (
	"github.com/djfemz/simbaCodingChallenge/data"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/stretchr/testify/assert"
)

var (
	Db                                         = data.Connect()
	transactionRepo data.TransactionRepository = &data.TransactionRepositoryImpl{}
)

func transactionSetUp() []*data.Transaction {
	return []*data.Transaction{{
		Amount:         200.00,
		SourceCurrency: data.DOLLAR,
	},
		{
			Amount:         100.00,
			SourceCurrency: data.DOLLAR,
		},
	}
}

func TestThatTransactionCanBeSaved(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction := transactionRepo.Save(transactions[0])
	assert.NotEmpty(t, savedTransaction)
	assert.Greater(t, savedTransaction.ID, uint(0))
}

func TestThatTransactionCanBeFoundById(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction := transactionRepo.Save(transactions[0])
	returnedTransaction := transactionRepo.FindById(savedTransaction.ID)
	assert.Equal(t, returnedTransaction.ID, savedTransaction.ID)
	assert.Equal(t, returnedTransaction.Amount, 200.00)
}

func TestThatAllTransactionsCanBeRetrieved(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}
	assert.Equal(t, 2, len(transactionRepo.FindAllTransactions()))
}

func TestThatTransactionCanBeDeletedById(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}
	transactionRepo.DeleteById(transactions[0].ID)
	assert.Equal(t, 1, len(transactionRepo.FindAllTransactions()))
}

func TestThatATransactionCanBeDeleted(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}

	transactionRepo.Delete(transactions[1])
	assert.Equal(t, 1, len(transactionRepo.FindAllTransactions()))
}
