package services

import (
	"github.com/djfemz/simbaCodingChallenge/dtos"
)

type UserService interface{
	Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse
}

type UserServiceImpl struct{

}

func (userserviceImpl *UserServiceImpl) Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse{

	return dtos.AddUserResponse{}
}