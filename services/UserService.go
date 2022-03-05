package services

import (
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)


var (
	Db *gorm.DB
	userRepo = repositories.UserRepositoryImpl{}
)

type UserService interface{
	Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse
}

type UserServiceImpl struct{

}

func setUp(){
	Db=repositories.Connect()
}

func (userserviceImpl *UserServiceImpl) Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse{
	setUp()
	defer Db.Close()
	var addUserResponse = dtos.AddUserResponse{}
	if len(addUserDto.Name)<1 || len(addUserDto.Email)<1 || len(addUserDto.Password)<8{
		return dtos.AddUserResponse{}
	}else{
		password, err:=hashPassword(addUserDto.Password)
		if err!=nil{
			log.Println("couldn't hash password")
			return dtos.AddUserResponse{}
		}
		var user = &models.User{Name: addUserDto.Name, Email: addUserDto.Email, Password: password}
		savedUser:=userRepo.Save(user)
		addUserResponse.Name = savedUser.Name
		return addUserResponse
	}
}

func hashPassword(password string) (hash string, err error){
	byteSlice, err:=bcrypt.GenerateFromPassword([]byte(password), 15)
	if err!=nil{
		log.Println("Error hashing password")
		return err.Error(), err
	}
	hash=string(byteSlice)
	return hash, nil
}