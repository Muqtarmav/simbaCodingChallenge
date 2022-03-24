package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/data"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Transaction struct{}

var (
	transactionService services.TransactionService = services.TransactionServiceImpl{}
)

func (transaction *Transaction) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transactionType := req.FormValue("transaction-type")
	if req.URL.Path == "/transaction" &&
		transactionType == "transfer" && req.Method == http.MethodPost {
		transfer(rw, req)
	} else if req.URL.Path == "/transaction" &&
		transactionType == "convert" && req.Method == http.MethodPost {
		convert(rw, req)
	} else if req.URL.Path == "/transaction/new" {
		newTransaction(rw, req)
	} else if req.URL.Path == "/transaction/overview" {
		transactionOverview(rw, req)
	}
}

func newTransaction(rw http.ResponseWriter, req *http.Request) {
	log.Println("request for new transaction---->", req)
	session, err := data.ReturnSession(rw, req)
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/" +
		"simbaCodingChallenge/views/templates/transaction.html"))

	err = tmpl.Execute(rw, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
		"UserID":         session.UserID,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func transactionOverview(rw http.ResponseWriter, req *http.Request) {
	request := getRequestParameters(req)
	usersTransactions := transactionService.GetUsersTransactions(request.UserID)
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/" +
		"views/templates/overview.html"))
	err := tmpl.Execute(rw, usersTransactions)
	if err != nil {
		log.Fatal(err)
	}
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func transfer(rw http.ResponseWriter, req *http.Request) {
	log.Println("here---->in transfer")

	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = data.TRANSFER
	response := transactionService.Deposit(transactionRequest)
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/"+
		"views/templates/"+
		"failed-transaction.html", "/home/djfemz/Documents/goworkspace/github.com/"+
		"simbaCodingChallenge/views/templates/overview.html", "/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/"+
		"views/templates/succesful-transaction.html"))
	if response.Status == data.SUCCESS {
		err := tmpl.ExecuteTemplate(rw, "succesful-transaction.html", response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := tmpl.ExecuteTemplate(rw, "failed-transaction.html", response)
		if err != nil {
			log.Fatal(err)
		}
	}
	transactionOverview(rw, req)
}

func convert(rw http.ResponseWriter, req *http.Request) {
	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = data.CONVERT
	response := transactionService.ConvertMoney(transactionRequest)
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/" +
		"simbaCodingChallenge/views/templates/overview.html"))
	err := tmpl.Execute(rw, response)
	if err != nil {
		log.Fatal(err)
	}
}

func getRequestParameters(req *http.Request) dtos.TransactionRequest {
	amount, err := strconv.ParseFloat(req.FormValue("amount"), 0)
	if err != nil {
		log.Fatal(err)
	}
	sourceCurrency := req.FormValue("source-currency")
	currency := checkCurrency(sourceCurrency)
	targetCurrency := req.FormValue("target-currency")
	currencyTarget := checkCurrency(targetCurrency)
	userID, err := strconv.Atoi(req.FormValue("user-id"))
	if err != nil {
		log.Fatal(err)
	}
	recipientID, err := strconv.Atoi(req.FormValue("recipient"))
	if err != nil {
		log.Fatal(err)
	}

	transactionRequest := dtos.TransactionRequest{
		Amount:         amount,
		SourceCurrency: currency,
		TargetCurrency: currencyTarget,
		UserID:         uint(userID),
		RecipientsID:   uint(recipientID),
	}
	return transactionRequest
}

func checkCurrency(currency string) data.Currency {
	if currency == "USD" {
		return data.DOLLAR
	} else if currency == "EUR" {
		return data.EURO
	} else {
		return data.NAIRA
	}
}
