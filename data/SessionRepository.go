package data

import (
	"github.com/jinzhu/gorm"
	"log"
)

type SessionRepository interface {
	Save(session *Session) *Session
	FindByUserID(id uint) *Session
	Delete(session *Session)
}

type SessionRepositoryImpl struct {
}

func (sessionRepo *SessionRepositoryImpl) Save(session *Session) *Session {
	var savedSession = &Session{}
	Db := Connect()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)
	Db.Create(&session)
	Db.Where("user_id=?", session.UserID).First(savedSession)
	return savedSession
}

func (sessionRepo *SessionRepositoryImpl) FindByUserID(id uint) *Session {
	var foundSession = &Session{}
	Db := Connect()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)
	Db.Where("user_id=?", id).Last(&foundSession)
	return foundSession
}

func (sessionRepo *SessionRepositoryImpl) Delete(session *Session) {
	Db := Connect()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)
	Db.Delete(session)
}
