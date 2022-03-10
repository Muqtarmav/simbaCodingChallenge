package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"log"
	"net/http"
	"text/template"
)

var (
	userService services.UserService = &services.UserServiceImpl{}
)

type User struct{}

func (user *User) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/register" {
		signup(rw, req)
	} else if req.URL.Path == "/login" {
		login(rw, req)
	}
}

func signup(rw http.ResponseWriter, req *http.Request) {
	registerForm := dtos.AddUserRequest{
		Name:     req.FormValue("name"),
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}
	response := userService.Register(registerForm)
	tmpl := template.Must(template.ParseFiles("overview.html"))
	err := tmpl.Execute(rw, response)
	if err != nil {
		log.Fatal(err)
	}
}

func login(rw http.ResponseWriter, req *http.Request) {
	loginRequest := dtos.LoginRequest{
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}
	response := userService.Login(loginRequest)
	tmpl := template.Must(template.ParseFiles("overview.html"))
	err := tmpl.Execute(rw, response)
	if err != nil {
		log.Fatal(err)
	}
}
