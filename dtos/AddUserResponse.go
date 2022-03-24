package dtos

import (
	"github.com/djfemz/simbaCodingChallenge/data"
)

type AddUserResponse struct {
	ID      uint
	Name    string
	Balance []data.Money
}
