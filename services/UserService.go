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
	Login(loginRequest dtos.LoginRequest) dtos.LoginResponse
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
		var user = &models.User{
			Name: addUserDto.Name, 
			Email: addUserDto.Email, 
			Password: password,
			Balance: []models.Money{
				{Amount: 0, Currency:models.NAIRA }, 
				{Amount: 0, Currency: models.POUNDS}, 
				{Amount: 1000, Currency: models.DOLLAR},
			},
		}
		savedUser:=userRepo.Save(user)
		addUserResponse.Name = savedUser.Name
		addUserResponse.ID=savedUser.ID
		return addUserResponse
	}
}



func(userserviceImpl *UserServiceImpl) Login(loginRequest dtos.LoginRequest) dtos.LoginResponse{
	return dtos.LoginResponse{}
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