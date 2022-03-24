package services_test

import (
	"github.com/djfemz/simbaCodingChallenge/data"
	"log"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"github.com/stretchr/testify/assert"
)

var (
	userService services.UserService = &services.UserServiceImpl{}
	userRepo    data.UserRepository  = &data.UserRepositoryImpl{}
)

func TestThatUserCanBeRegistered(t *testing.T) {
	var addUserRequest = dtos.AddUserRequest{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Password: "12345678",
	}

	addUserResponse := userService.Register(addUserRequest)
	assert.NotEmpty(t, addUserResponse)
	log.Println(addUserResponse)
	assert.Equal(t, addUserResponse.Name, addUserRequest.Name)
}

func TestThatUserIsntRegisteredWhenPasswordLessThan_8_Characters(t *testing.T) {
	var addUserRequest = dtos.AddUserRequest{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Password: "1234567",
	}
	addUserResponse := userService.Register(addUserRequest)
	assert.Empty(t, addUserResponse)
}

func TestThatEveryRegisteredUserGets_1000_USD_Upon_Registration(t *testing.T) {
	var addUserRequest = dtos.AddUserRequest{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Password: "12345678",
	}

	//register user
	addUserResponse := userService.Register(addUserRequest)
	//assert that user is registered
	assert.NotEmpty(t, addUserResponse)
	assert.Greater(t, addUserResponse.ID, uint(0))
	//find user
	log.Println("user response-->", addUserResponse)
	savedUser := userRepo.FindById(addUserResponse.ID)
	log.Println("found from db-->", savedUser)
	//assert that balance is not empty
	assert.NotEmpty(t, savedUser.Balance)
	for _, balance := range savedUser.Balance {
		if balance.Currency == data.data.DOLLAR {
			assert.Equal(t, 1000.00, balance.Amount)
		}
	}

}

func TestThatUserCanLoginWithEmailAndPassword(t *testing.T) {
	var loginRequest = dtos.LoginRequest{
		Email:    "john@gmail.com",
		Password: "12345678",
	}

	loginResponse := userService.Login(loginRequest)
	assert.NotEmpty(t, loginResponse)
	log.Println("logged in user-->", loginResponse)
	assert.Equal(t, "user loggedin successfully", loginResponse.Message)
}

func TestThatUserCannotLoginWithWrongCredentials(t *testing.T) {
	var loginRequest = dtos.LoginRequest{
		Email:    "john@gmail.com",
		Password: "1234567",
	}

	loginResponse := userService.Login(loginRequest)
	assert.NotEmpty(t, loginResponse)
	log.Println("logged in user-->", loginResponse)
	assert.Equal(t, "bad login credentials", loginResponse.Message)
}

func TestThatUserBalanceCanBeCalculatedFromTransactionHistory(t *testing.T) {
	balance, err := userService.GetAccountBalance(68)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(balance)
	//assert.Equal(t, 3800.00, balance)
}
