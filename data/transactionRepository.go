package data

import (
	"log"
)

type TransactionRepository interface {
	Save(transaction *Transaction) *Transaction
	FindById(id uint) *Transaction
	FindTransactionsByUserID(userID uint) []Transaction
	FindAllTransactions() []*Transaction
	DeleteById(id uint)
	Delete(transaction *Transaction)
}

type TransactionRepositoryImpl struct {
}

func (transactionRepo *TransactionRepositoryImpl) Save(transaction *Transaction) *Transaction {
	Db := Connect()
	defer Db.Close()
	savedTransaction := &Transaction{}
	Db.Create(transaction)
	Db.Where("ID=?", transaction.ID).Find(savedTransaction)
	return savedTransaction
}

func (transactionRepo *TransactionRepositoryImpl) FindTransactionsByUserID(userID uint) []Transaction {
	log.Println("user id---->", userID)
	Db := Connect()
	defer Db.Close()
	transactions := []Transaction{}
	Db.Where("user_id=?", userID).Or("receivers_id=?", userID).Find(&transactions)
	return transactions
}

func (transactionRepo *TransactionRepositoryImpl) FindById(id uint) *Transaction {
	Db := Connect()
	defer Db.Close()
	foundTransaction := &Transaction{}
	Db.Where("ID=?", id).Find(foundTransaction)
	return foundTransaction
}

func (transactionRepo *TransactionRepositoryImpl) FindAllTransactions() []*Transaction {
	Db := Connect()
	defer Db.Close()
	var transactions []*Transaction
	Db.Preload("Balance").Find(&transactions)
	return transactions
}

func (transactionRepo *TransactionRepositoryImpl) DeleteById(id uint) {
	Db := Connect()
	defer Db.Close()
	Db.Where("ID=?", id).Delete(&Transaction{})
}

func (transactionRepo *TransactionRepositoryImpl) Delete(transaction *Transaction) {
	transactionRepo.DeleteById(transaction.ID)
}
