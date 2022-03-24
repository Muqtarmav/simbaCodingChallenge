package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	userService services.UserService = &services.UserServiceImpl{}
)

type User struct{}

func (user *User) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/user/register" && req.Method == http.MethodPost {
		signup(rw, req)
	} else if req.URL.Path == "/user/login" && req.Method == http.MethodPost {
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
		login(rw, req)
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
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/" +
		"simbaCodingChallenge/views/templates/overview.html"))
	if response.Message == "user loggedin successfully" {
		user, err := userService.GetUser(response.ID)
		log.Println("user retrieved---->", user)
		if err != nil {
			log.Fatal(err)
		}
		session, err := user.CreateSession()
		if err != nil {
			log.Fatal("error creating session---->", err)
		}
		log.Println("session after creation---->", session)
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    strconv.Itoa(int(session.UserID)),
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(rw, &cookie)
		log.Println("cookie sent to browser---->", cookie)
		log.Println("response---->", response)
		log.Println("header----->", rw.Header())
		err = tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("redirecting in login")
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}

}
