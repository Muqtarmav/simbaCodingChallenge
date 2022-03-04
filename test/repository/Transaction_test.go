package test

import (
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
)

var transactionRepo repositories.TransactionRepository = &repositories.TransactionRepositoryImpl{}

func transactionSetUp() *models.Transaction{
	return &models.Transaction{
		Id: 1,
	}
}

func TestThatTransactionCanBeSaved(t *testing.T) {
	transactionRepo.Save()
}
