package data

import (
	"log"

	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserRepository interface {
	Save(user *User) *User
	FindById(id uint) *User
	FindAllUsers() []*User
	FindByEmail(email string) *User
	DeleteById(id uint)
}

var (
	config util.Config
	Db     *gorm.DB
	err    error
)

func Connect() *gorm.DB {
	config, err = util.LoadConfig("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge")
	if err != nil {
		log.Fatal(err)
	}
	Db, err = gorm.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	Db.AutoMigrate(&User{}, &Transaction{}, &Money{}, &Session{})
	log.Println("connected " + "to db")
	return Db
}

type UserRepositoryImpl struct {
}

func (userRepo *UserRepositoryImpl) Save(user *User) *User {
	// userRepo.l.Println("in save")
	Db := Connect()
	defer Db.Close()
	savedUser := &User{}
	Db.Create(user)
	Db.Where("ID=?", user.ID).Find(&savedUser)
	log.Println("saved user is -->", savedUser)
	return savedUser
}

func (userRepo *UserRepositoryImpl) UpdateUserDetails(userToBeUpdated *User) {
	Db := Connect()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)
	var user User
	log.Println("user to be updated---->", userToBeUpdated)
	Db.Preload("Balance").First(&user, userToBeUpdated.ID)
	log.Println("user to be updated-->", userToBeUpdated)
	user.Balance = userToBeUpdated.Balance
	log.Println("user after updating balance, before save", user)
	Db.Save(user)
}

func (userRepo *UserRepositoryImpl) FindByEmail(email string) *User {
	Db := Connect()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)
	savedUser := &User{}
	Db.Where("email=?", email).Preload("Balance").Preload("Transactions").Find(&savedUser)

	if savedUser == nil {
		return nil
	}
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindById(id uint) *User {
	Db := Connect()
	defer Db.Close()
	savedUser := &User{}
	Db.Omit("CreatedAt", "UpdatedAt", "DeletedAt").Where("id=?", id).
		Omit("CreatedAt", "UpdatedAt", "DeletedAt").Preload("Balance").
		Preload("Transactions").Find(&savedUser).Omit("CreatedAt", "UpdatedAt", "DeletedAt")
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindAllUsers() []*User {
	Db := Connect()
	defer Db.Close()
	var users []*User
	Db.Find(&users)
	return users
}

func (userRepo *UserRepositoryImpl) DeleteById(id uint) {
	Db := Connect()
	defer Db.Close()
	Db.Where("Id=?", id).Delete(&User{})
}
