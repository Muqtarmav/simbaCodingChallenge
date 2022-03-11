package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"html/template"
	"log"
	"net/http"
)

var (
	userService services.UserService = &services.UserServiceImpl{}
)

type User struct{}

func (user *User) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("here----->")
	if req.URL.Path == "/user/register" && req.Method == http.MethodPost {
		log.Println("signup ----> in servehttp")
		signup(rw, req)
	} else if req.URL.Path == "/user/login" && req.Method == http.MethodPost {
		log.Println("login ----> in servehttp")
		login(rw, req)
	}
}

func GetUser() *User {
	return &User{}
}

func signup(rw http.ResponseWriter, req *http.Request) {
	log.Println("signing up")
	registerForm := dtos.AddUserRequest{
		Name:     req.FormValue("name"),
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}
	log.Println("sent to user service----<")
	response := userService.Register(registerForm)
	if response.ID > 0 {
		tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/views/templates/overview.html"))
		err := tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("redirecting in signup---->", response)
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}

}

func login(rw http.ResponseWriter, req *http.Request) {
	loginRequest := dtos.LoginRequest{
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}
	response := userService.Login(loginRequest)
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/views/templates/overview.html"))
	if response.Message == "user loggedin successfully" {
		err := tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("redirecting in login")
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}

}
