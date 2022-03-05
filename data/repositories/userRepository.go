package repositories

import (
	"log"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserRepository interface{
	Save(user *models.User) *models.User
	FindById(id uint) *models.User
	FindAllUsers() []*models.User
	DeleteById(id uint)
}

var (
	config util.Config
	Db *gorm.DB
	err error
)

func Connect() *gorm.DB{
	config, err =util.LoadConfig("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge")
	if err!=nil{
		log.Fatal(err)
	}
	Db, err = gorm.Open(config.DBDriver, config.DBSource)
	if err!=nil{
		log.Fatal(err)
	}
	Db.AutoMigrate(&models.User{}, &models.Transaction{})
	log.Println("connected " + "to db")
	return Db
}

type UserRepositoryImpl struct{

}

func (userRepo *UserRepositoryImpl) Save(user *models.User) *models.User {
	// userRepo.l.Println("in save")
	savedUser := &models.User{}
	Db.Create(user)
	Db.Where("Id=?", user.ID).Find(&savedUser)
	return savedUser
}


func (userRepo *UserRepositoryImpl) FindById(id uint) *models.User{
	savedUser := &models.User{}
	Db.Where("Id=?", id).Find(&savedUser)
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindAllUsers() []*models.User{
	var users []*models.User
	Db.Find(&users)
	return users
}

func (userRepo *UserRepositoryImpl) DeleteById(id uint){
	Db.Where("Id=?", id).Delete(&models.User{})
}