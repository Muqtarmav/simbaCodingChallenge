package repositories

import "github.com/djfemz/simbaCodingChallenge/data/models"

type UserRepository interface{
	Save(user *models.User) *models.User
	FindById(id int) *models.User
	FindAllUsers() []*models.User
	DeleteById(id int)
}

type UserRepositoryImpl struct{
	
}