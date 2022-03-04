package test

import (
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/stretchr/testify/assert"
)


var (
	userRepo repositories.UserRepository = &repositories.UserRepositoryImpl{} 
)



func setUp() *models.User{
	return &models.User{
		Id: 3,
		Name: "Janey Doe",
		Email: "janeydoe@email.com",
		Password: "1234",
	}
}

func TestThatUserCanBeSaved(t *testing.T){ 
	user := setUp()
	returnedUser := userRepo.Save(user)
	assert.NotEmpty(t, returnedUser)
	assert.Equal(t, user.Email, returnedUser.Email)
}

func TestThatUserCanBeFoundById(t *testing.T){
	returnedUser := userRepo.FindById(1)
	assert.Equal(t, returnedUser.Id, 1)
	assert.Equal(t, returnedUser.Name, "John Doe")
}

func TestThatAllUsersCanBeFound(t *testing.T){
	allUsers:=userRepo.FindAllUsers()
	assert.Equal(t, 3, len(allUsers))
}

func TestDeleteById(t *testing.T) {
	assert.Equal(t, 3, len(userRepo.FindAllUsers()))
	userRepo.DeleteById(3)
	assert.Equal(t, 2, len(userRepo.FindAllUsers()))
}