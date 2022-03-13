package dtos

import "github.com/djfemz/simbaCodingChallenge/data/models"

type AddUserResponse struct {
	ID      uint
	Name    string
	Balance []models.Money
}
