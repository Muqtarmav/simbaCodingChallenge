package dtos

import (
	"github.com/djfemz/simbaCodingChallenge/data"
)

type LoginResponse struct {
	ID           uint
	Name         string
	Message      string
	Balance      []data.Money
	Transactions []data.Transaction
}
