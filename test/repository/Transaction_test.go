package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var Db = repositories.Connect()


//clears models that are created during tests
func DeleteCreatedModels(db *gorm.DB) func(){
	type entity struct{
		table string
		keyname string 
		key interface{}
	}

	var entries []entity
	hookName:="cleanupHook"

	db.Callback().Create().After("gorm:create").Register(hookName,
		 func(scope *gorm.Scope) {
			fmt.Printf("Inserted entities of %s with %s = %v\n",
			scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
			entries = append(entries, entity{table: scope.TableName(), 
			keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
		 })
	return func() {
		defer db.Callback().Create().Remove(hookName)
		_, inTransaction:= db.CommonDB().(*sql.Tx)
		tx:=db
		if!inTransaction{
			tx = db.Begin()
		}
		for i:= len(entries) -1;i>=0;i--{
			entry := entries[i]
			fmt.Printf("Deleting entities from %s table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+"=?", entry.key).Delete("")
		}
		if !inTransaction{
			tx.Commit()
		}
	}	 
}


var transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}

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
	cleaner := DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction:=transactionRepo.Save(transactions[0])
	assert.NotEmpty(t,savedTransaction)
	assert.Greater(t, savedTransaction.ID, uint(0))
}


func TestThatTransactionCanBeFoundById(t *testing.T){
	cleaner:=DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()
	savedTransaction:=transactionRepo.Save(transactions[0])
	returnedTransaction := transactionRepo.FindById(savedTransaction.ID)
	assert.Equal(t, returnedTransaction.ID, savedTransaction.ID)
	assert.Equal(t, returnedTransaction.Amount, 200.00)
}

func TestThatAllTransactionsCanBeRetrieved(t *testing.T){
	cleaner:=DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}
	assert.Equal(t, 2, len(transactionRepo.FindAllTransactions()))
}

func TestThatTransactionCanBeDeletedById(t *testing.T)  {
	cleaner:=DeleteCreatedModels(Db)
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
	cleaner:=DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	transactions := transactionSetUp()

	for _, transaction := range transactions {
		transactionRepo.Save(transaction)
	}

	transactionRepo.Delete(transactions[1])
	assert.Equal(t, 1, len(transactionRepo.FindAllTransactions()))
}