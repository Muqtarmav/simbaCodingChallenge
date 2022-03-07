package repositories

import "github.com/djfemz/simbaCodingChallenge/data/models"


type TransactionRepository interface{
	Save(transaction *models.Transaction) *models.Transaction
	FindById(id uint) *models.Transaction
	FindAllTransactions() []*models.Transaction
	DeleteById(id uint)
	Delete(transaction *models.Transaction)
}

type TransactionRepositoryImpl struct{

}

func (transactionRepo *TransactionRepositoryImpl) Save(transaction *models.Transaction) (*models.Transaction) {
	Db:=Connect()
	defer Db.Close()
	savedTransaction := &models.Transaction{}
	Db.Create(transaction)
	Db.Where("ID=?", transaction.ID).Find(savedTransaction)
 	return savedTransaction
}


func (transactionRepo *TransactionRepositoryImpl) FindById(id uint) *models.Transaction{
	Db:=Connect()
	defer Db.Close()
	foundTransaction:= &models.Transaction{}
	Db.Where("ID=?", id).Find(foundTransaction)
	return foundTransaction
}

func (transactionRepo *TransactionRepositoryImpl) FindAllTransactions() []*models.Transaction{
	Db:=Connect()
	defer Db.Close()
	var transactions []*models.Transaction
	Db.Find(&transactions)
	return transactions
}

func (transactionRepo *TransactionRepositoryImpl) DeleteById(id uint){
	Db:=Connect()
	defer Db.Close()
	Db.Where("ID=?", id).Delete(&models.Transaction{})
}

func (transactionRepo *TransactionRepositoryImpl) Delete(transaction *models.Transaction){
	transactionRepo.DeleteById(transaction.ID)
}