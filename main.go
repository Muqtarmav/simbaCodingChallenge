package main

import (
	"github.com/djfemz/simbaCodingChallenge/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func main() {
	user := handlers.GetUser()
	router := mux.NewRouter()
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/views/static"))))
	router.HandleFunc("/", index)
	router.Handle("/transaction", &handlers.Transaction{})
	router.Handle("/user/register", user)
	router.Handle("/user/login", user)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func index(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/templates/login.html"))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		log.Fatal(err)
	}
}
