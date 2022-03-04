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
	FindById(id int) *models.User
	FindAllUsers() []*models.User
	DeleteById(id int)
}

var (
	config util.Config
	Db *gorm.DB
	err error
)

func Connect(){
	config, err =util.LoadConfig("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge")
	if err!=nil{
		log.Fatal(err)
	}
	Db, err = gorm.Open(config.DBDriver, config.DBSource)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("connected " + "to db")
}

type UserRepositoryImpl struct{

}

func (userRepo *UserRepositoryImpl) Save(user *models.User) *models.User {
	// userRepo.l.Println("in save")
	savedUser := &models.User{}
	Connect()
	defer Db.Close()
	Db.Create(user)
	Db.Where("Id=?", user.Id).Find(&savedUser)
	log.Println(&savedUser)
	return savedUser
}


func (userRepo *UserRepositoryImpl) FindById(id int) *models.User{
	savedUser := &models.User{}
	Connect()
	defer Db.Close()
	Db.Where("Id=?", id).Find(&savedUser)
	return savedUser
}

func (userRepo *UserRepositoryImpl) FindAllUsers() []*models.User{
	var users []*models.User
	Connect()
	defer Db.Close()
	Db.Find(&users)
	return users
}

func (userRepo *UserRepositoryImpl) DeleteById(id int){
	Connect()
	defer Db.Close()
	Db.Where("Id=?", id).Delete(&models.User{})
}