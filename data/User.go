package data

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique"`
	Password     string
	Balance      []Money
	Transactions []Transaction
}

type Session struct {
	gorm.Model
	Email  string
	UserID uint
}

var sessionRepo SessionRepository = &SessionRepositoryImpl{}

func (user *User) CreateSession() (*Session, error) {
	log.Println("in create session---->", user)
	var session = &Session{}
	session.Email = user.Email
	session.UserID = user.ID
	savedSession := sessionRepo.Save(session)
	return savedSession, nil
}

func (user *User) CheckSession() (session *Session, err error) {
	foundSession := sessionRepo.FindByUserID(user.ID)
	return foundSession, nil
}

func (session *Session) Check() (valid bool, err error) {
	foundSession := sessionRepo.FindByUserID(session.UserID)
	if foundSession.ID != 0 {
		valid = true
		return
	}
	if foundSession == nil {
		valid = false
		return
	}
	return
}

func ReturnSession(w http.ResponseWriter, req *http.Request) (Session, error) {
	cookie, err := req.Cookie("_cookie")
	log.Println("cookie retrieved---->", cookie)
	if err != nil {
		log.Fatal(err)
	}
	val, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	sess := Session{
		UserID: uint(val),
	}
	if ok, _ := sess.Check(); !ok {
		log.Fatal("session not valid")
	}
	return sess, nil
}
