package test

import (
	"github.com/djfemz/simbaCodingChallenge/data"
	"github.com/jinzhu/gorm"
	"log"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/stretchr/testify/assert"
)

var (
	userRepo repositories.UserRepository = &repositories.UserRepositoryImpl{}
)

func setUp() []*data.data {
	return []*data.User{
		{
			Name:     "Janey Doe",
			Email:    "janeydoe@email.com",
			Password: "1234",
		},
		{
			Name:     "John Doe",
			Email:    "john@gmail.com",
			Password: "1234",
		},
	}
}

func TestThatUserCanBeSaved(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	users := setUp()
	returnedUser := userRepo.Save(users[0])
	assert.NotEmpty(t, returnedUser)
	assert.Equal(t, users[0].Email, returnedUser.Email)
}

func TestThatUserCanBeFoundById(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {

		}
	}(Db)
	defer cleaner()
	//users := setUp()
	//savedUser := userRepo.Save(users[0])
	returnedUser := userRepo.FindById(3)
	log.Println(returnedUser)
	assert.Equal(t, returnedUser.ID, uint(3))
	assert.Equal(t, returnedUser.Name, "Janey Doe")
}

func TestThatUserCanBeFoundByEmail(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	//users := setUp()
	//for _, user := range users {
	//	userRepo.Save(user)
	//}
	foundUser := userRepo.FindByEmail("janeydoe@email.com")
	assert.NotEmpty(t, foundUser)
	log.Println("found user---->", foundUser)
}

func TestThatAllUsersCanBeFound(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	users := setUp()
	for _, user := range users {
		userRepo.Save(user)
	}
	allUsers := userRepo.FindAllUsers()
	assert.Equal(t, 2, len(allUsers))
}

func TestDeleteById(t *testing.T) {
	cleaner := util.DeleteCreatedModels(Db)
	defer Db.Close()
	defer cleaner()
	users := setUp()
	for _, user := range users {
		userRepo.Save(user)
	}
	assert.Equal(t, 2, len(userRepo.FindAllUsers()))
	userRepo.DeleteById(users[0].ID)
	assert.Equal(t, 1, len(userRepo.FindAllUsers()))
}
