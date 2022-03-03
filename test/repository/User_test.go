package test

import (
	"context"
	"log"
	"testing"

	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var source string = "user=postgres password=debbie200 dbname=simba_db sslmode=disable"

var (
	Db  *gorm.DB
	err  error
)

func init()  {
	Db, err = gorm.Open("postgres", source)
	if err!=nil{
		log.Fatal(err.Error())
	}
	conn, err :=Db.DB().Conn(context.TODO())
	if err!=nil{
		log.Fatal(err.Error())
	}
	if conn!=nil{
		log.Println("connected successfully")
	}
	Db.AutoMigrate(&models.User{})
}

func setUp() *models.User{
	Db = Db.Delete(&models.User{}, []int{1,2,3})
	return &models.User{
		Id: 1,
		Name: "John Doe",
		Email: "johndoe@email.com",
		Password: "1234",
	}
}

func TestThatUserCanBeSaved(t *testing.T){ 
	defer Db.Close()
	user1:= setUp()
	Db.Create(&user1)
	returnedUser := &models.User{}
	assert.Empty(t, returnedUser)
	Db.Where("Id=?", 1).Omit("created_at", "updated_at").First(&returnedUser)
	assert.NotEmpty(t, returnedUser)
}