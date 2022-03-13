package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"

type LoginResponse struct {
	ID           uint
	Name         string
	Message      string
	Balance      []models.Money
	Transactions []models.Transaction
}
