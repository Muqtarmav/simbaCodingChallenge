package services_test

import (
	"log"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"github.com/stretchr/testify/assert"
)


var userService services.UserService = &services.UserServiceImpl{}

func TestThatUserCanBeRegistered(t *testing.T)  {
	var addUserRequest = dtos.AddUserRequest{
		Name:"John Doe",
		Email:"john@gmail.com",
		Password:"12345678",
	}

	addUserResponse:= userService.Register(addUserRequest)
	assert.NotEmpty(t, addUserResponse)
	log.Println(addUserResponse)
	assert.Equal(t, addUserResponse.Name, addUserRequest.Name)
}

func TestThatUserIsntRegisteredWhenPasswordLessThan_8_Characters(t *testing.T){
	var addUserRequest = dtos.AddUserRequest{
		Name:"John Doe",
		Email:"john@gmail.com",
		Password:"1234567",
	}
	addUserResponse:= userService.Register(addUserRequest)
	assert.Empty(t, addUserResponse)
}