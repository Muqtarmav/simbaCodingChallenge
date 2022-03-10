package services

import (
	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/data/repositories"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	Db                 *gorm.DB
	userRepo           = repositories.UserRepositoryImpl{}
	transactionService = TransactionServiceImpl{}
)

type UserService interface {
	Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse
	Login(loginRequest dtos.LoginRequest) dtos.LoginResponse
	GetAccountBalance(userID uint) (float64, error)
}

type UserServiceImpl struct {
}

func setUp() {
	Db = repositories.Connect()
}

func (userServiceImpl *UserServiceImpl) Register(addUserDto dtos.AddUserRequest) dtos.AddUserResponse {
	setUp()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(Db)
	var addUserResponse = dtos.AddUserResponse{}
	if len(addUserDto.Name) < 1 || len(addUserDto.Email) < 1 ||
		len(addUserDto.Password) < 8 {
		return dtos.AddUserResponse{}
	} else {
		password, err := hashPassword(addUserDto.Password)
		if err != nil {
			log.Println("couldn't hash password")
			return dtos.AddUserResponse{}
		}
		var user = &models.User{
			Name:     addUserDto.Name,
			Email:    addUserDto.Email,
			Password: password,
			Balance: []models.Money{
				{Amount: 0, Currency: models.EURO},
				{Amount: 0, Currency: models.NAIRA},
				{Amount: 0, Currency: models.DOLLAR},
			},
			Transactions: []models.Transaction{},
		}
		savedUser := userRepo.Save(user)
		var transferRequest = dtos.TransactionRequest{
			Amount:          1000,
			SourceCurrency:  models.DOLLAR,
			TargetCurrency:  models.DOLLAR,
			UserID:          123456789,
			RecipientsID:    savedUser.ID,
			TransactionType: models.TRANSFER,
		}

		transactionService.Deposit(transferRequest)

		addUserResponse.Name = savedUser.Name
		addUserResponse.ID = savedUser.ID
		return addUserResponse
	}
}

func (userServiceImpl *UserServiceImpl) GetAccountBalance(userID uint) (float64, error) {
	savedUser := userRepo.FindById(userID)
	log.Println("user-->", savedUser)
	transactions := savedUser.Transactions

	log.Println(transactions)
	return 0, nil
}

func (userServiceImpl *UserServiceImpl) Login(loginRequest dtos.LoginRequest) dtos.LoginResponse {
	foundUser := userRepo.FindByEmail(loginRequest.Email)
	if foundUser == nil {
		return dtos.LoginResponse{Message: "user not found"}
	}

	if decryptPassword([]byte(foundUser.Password), []byte(loginRequest.Password)) {
		return dtos.LoginResponse{ID: foundUser.ID, Message: "user loggedin successfully"}
	} else {
		return dtos.LoginResponse{Message: "bad login credentials"}
	}
}

func validateUser(userID uint) (*models.User, error) {
	user := userRepo.FindById(userID)
	return user, nil
}

func hashPassword(password string) (hash string, err error) {
	byteSlice, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Println("Error hashing password")
		return err.Error(), err
	}
	hash = string(byteSlice)
	return hash, nil
}

func decryptPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
