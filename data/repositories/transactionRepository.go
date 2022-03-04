package repositories

import "github.com/djfemz/simbaCodingChallenge/data/models"


type TransactionRepository interface{
	Save(transaction *models.Transaction) *models.Transaction
	FindById(id int) *models.Transaction
	FindAllUsers() []*models.Transaction
	DeleteById(id int)
	Delete(transaction *models.Transaction)
}

type TransactionRepositoryImpl struct{

}

func (transactionRepo *TransactionRepositoryImpl) Save(transaction *models.Transaction) *models.Transaction {
 return nil
}


func (transactionRepo *TransactionRepositoryImpl) FindById(id int) *models.Transaction{
	return nil
}

func (transactionRepo *TransactionRepositoryImpl) FindAllUsers() []*models.Transaction{
	return nil
}

func (transactionRepo *TransactionRepositoryImpl) DeleteById(id int){
	Connect()
	defer Db.Close()
}

func (transactionRepo *TransactionRepositoryImpl) Delete(transaction *models.Transaction){
	Connect()
	defer Db.Close()
	transactionRepo.DeleteById(transaction.Id)
}