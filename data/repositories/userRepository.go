package repositories

import (
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserRepository interface {
	Save(user *models.User) *models.User
	FindById(id uint) *models.User
	FindAllUsers() []*models.User
	FindByEmail(email string) *models.User
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
	Db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Money{})
	log.Println("connected " + "to db")
	return Db
}

type UserRepositoryImpl struct {
}

func (userRepo *UserRepositoryImpl) Save(user *models.User) *models.User {
	// userRepo.l.Println("in save")
	Db := Connect()
	defer Db.Close()
	savedUser := &models.User{}
	log.Println("user to be saved is-->", user)
	Db.Create(user)
	Db.Where("Id=?", user.ID).Find(&savedUser)
	log.Println("saved user is -->", savedUser)
	return savedUser
}

func (userRepo *UserRepositoryImpl) UpdateUserDetails(userToBeUpdated *models.User) {
	Db := Connect()
	defer Db.Close()
	var user models.User
	log.Println("user to be updated", userToBeUpdated)
	Db.Preload("Balance").First(&user, userToBeUpdated.ID)
	log.Println("user to be updated-->", userToBeUpdated)
	user.Balance = userToBeUpdated.Balance
	log.Println("user after updating balance, before save", user)
	Db.Save(user)
}

func (userRepo *UserRepositoryImpl) FindByEmail(email string) *models.User {
	Db := Connect()
	defer Db.Close()
	savedUser := &models.User{}
	Db.Omit("created_at", "updated_at", "deleted_at").Preload("Balance").Preload("Transactions").
		Where("Email=?", email).Find(&savedUser)
	if savedUser == nil {
		return nil
	}
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindById(id uint) *models.User {
	Db := Connect()
	defer Db.Close()
	savedUser := &models.User{}
	Db.Where("ID=?", id).Preload("Transactions").Preload("Balance").Find(&savedUser)
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindAllUsers() []*models.User {
	Db := Connect()
	defer Db.Close()
	var users []*models.User
	Db.Find(&users)
	return users
}

func (userRepo *UserRepositoryImpl) DeleteById(id uint) {
	Db := Connect()
	defer Db.Close()
	Db.Where("Id=?", id).Delete(&models.User{})
}
