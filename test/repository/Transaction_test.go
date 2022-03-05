package test

import (
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/stretchr/testify/assert"
)



var (
	Db = repositories.Connect()
	transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}
)


func transactionSetUp() []*models.Transaction{
	return []*models.Transaction{ {
		Amount: 200.00,
		Currency: models.DOLLAR,
		Sender: models.User{Name: "Femi", Email: "femi@gmail.com"},
		Receiver: models.User{Name: "Adeola", Email: "adeola@gmail.com"},
	},
	{
		Amount: 100.00,
		Currency: models.DOLLAR,
		Sender: models.User{Name: "Jack", Email: "jack@gmail.com"},
		Receiver: models.User{Name: "Jill", Email: "jill@gmail.com"},
	},	
	}
}

func TestThatTransactionCanBeSaved(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction:=transactionRepo.Save(transactions[0])
	assert.NotEmpty(t,savedTransaction)
	assert.Greater(t, savedTransaction.ID, uint(0))
}


func TestThatTransactionCanBeFoundById(t *testing.T){
	cleaner:=util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction:=transactionRepo.Save(transactions[0])
	returnedTransaction := transactionRepo.FindById(savedTransaction.ID)
	assert.Equal(t, returnedTransaction.ID, savedTransaction.ID)
	assert.Equal(t, returnedTransaction.Amount, 200.00)
}

func TestThatAllTransactionsCanBeRetrieved(t *testing.T){
	cleaner:=util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}
	assert.Equal(t, 2, len(transactionRepo.FindAllTransactions()))
}

func TestThatTransactionCanBeDeletedById(t *testing.T)  {
	cleaner:=util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}
	transactionRepo.DeleteById(transactions[0].ID)
	assert.Equal(t, 1, len(transactionRepo.FindAllTransactions()))
}

func TestThatATransactionCanBeDeleted(t *testing.T){
	cleaner:=util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}

	transactionRepo.Delete(transactions[1])
	assert.Equal(t, 1, len(transactionRepo.FindAllTransactions()))
}