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
	l *log.Logger
}

func (userRepo *UserRepositoryImpl) Save(user *models.User) *models.User {
	userRepo.l.Println("in save")
	return nil
}


func (userRepo *UserRepositoryImpl) FindById(id int) *models.User{
	return nil
}

func (userRepo *UserRepositoryImpl) FindAllUsers() []*models.User{
	return nil
}

func (userRepo *UserRepositoryImpl) DeleteById(id int){
	
}