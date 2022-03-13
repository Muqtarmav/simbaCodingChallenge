package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Transaction struct{}

var (
	transactionService services.TransactionService = services.TransactionServiceImpl{}
)

func (transaction *Transaction) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transactionType := req.FormValue("transaction")
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
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/" +
		"simbaCodingChallenge/views/templates/transaction.html"))
	//userID, err := strconv.Atoi(req.FormValue("userID"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	err := tmpl.Execute(rw, 13)
	if err != nil {
		log.Fatal(err)
	}
}

func transactionOverview(rw http.ResponseWriter, req *http.Request) {
	userID, err := strconv.Atoi(req.FormValue("userID"))
	if err != nil {
		log.Fatal(err)
	}
	usersTransactions := transactionService.GetUsersTransactions(uint(userID))
	tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/" +
		"views/templates/overview.html"))
	err = tmpl.Execute(rw, usersTransactions)
	if err != nil {
		log.Fatal(err)
	}
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func transfer(rw http.ResponseWriter, req *http.Request) {
	log.Println("here---->")

	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = models.TRANSFER
	response := transactionService.Deposit(transactionRequest)
	if response.Status == models.SUCCESS {
		response.Image = "/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/views/img/successful.png"
		tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/" +
			"views/templates/succesful-transaction.html"))
		err := tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		response.Image = "/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/views/img/icons8-multiply-90.png"
		tmpl := template.Must(template.ParseFiles("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/" +
			"views/templates/succesful-transaction.html"))
		err := tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 2)
		http.Redirect(rw, req, "/transaction/overview", http.StatusTemporaryRedirect)
	}
}

func convert(rw http.ResponseWriter, req *http.Request) {
	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = models.CONVERT
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

func checkCurrency(currency string) models.Currency {
	if currency == "USD" {
		return models.DOLLAR
	} else if currency == "EUR" {
		return models.EURO
	} else {
		return models.NAIRA
	}
}
