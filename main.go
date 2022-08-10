package main

import (
	"github.com/djfemz/simbaCodingChallenge/handlers"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func main() {
	user := handlers.GetUser()
	transaction := handlers.NewTransaction()
	router := mux.NewRouter()
	csrfMiddleware := csrf.Protect([]byte("secret"), csrf.Path("/"), csrf.Secure(false))
	router.Use(csrfMiddleware)
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("/home/djfemz/Documents/goworkspace/"+
				"github.com/simbaCodingChallenge/views/static"))))
	router.PathPrefix("/img/").
		Handler(http.StripPrefix("/img/",
			http.FileServer(http.Dir("/home/djfemz/Documents/goworkspace/"+
				"github.com/simbaCodingChallenge/views/img"))))
	router.HandleFunc("/", index)
	router.Handle("/transaction", transaction)
	router.Handle("/transaction/new", transaction)
	router.Handle("/transaction/overview", transaction)
	router.Handle("/user/register", user)
	router.Handle("/user/login", user)

	server := http.Server{
		Addr:    ":8082",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func index(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/templates/login.html"))
	rw.Header().Set("X-CSRF-Token", csrf.Token(req))
	err := tmpl.Execute(rw, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
	})
	if err != nil {
		log.Fatal(err)
	}
}
