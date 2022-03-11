package handlers

import (
	"github.com/djfemz/simbaCodingChallenge/data/models"
	"github.com/djfemz/simbaCodingChallenge/dtos"
	"github.com/djfemz/simbaCodingChallenge/services"
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
	transactionType := req.FormValue("transaction")
	if req.URL.Path == "/transaction" &&
		transactionType == "transfer" && req.Method == http.MethodPost {
		transfer(rw, req)
	} else if req.URL.Path == "/transaction" &&
		transactionType == "convert" && req.Method == http.MethodPost {
		convert(rw, req)
	}
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func transfer(rw http.ResponseWriter, req *http.Request) {
	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = models.TRANSFER
	response := transactionService.Deposit(transactionRequest)
	if response.Status == models.SUCCESS {
		tmpl := template.Must(template.ParseFiles("overview.html"))
		err := tmpl.Execute(rw, response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(rw, req, "/", 303)
	}

}

func convert(rw http.ResponseWriter, req *http.Request) {
	transactionRequest := getRequestParameters(req)
	transactionRequest.TransactionType = models.CONVERT
	response := transactionService.ConvertMoney(transactionRequest)
	tmpl := template.Must(template.ParseFiles("overview.html"))
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
	sourceCurrency := req.FormValue("source")
	currency := checkCurrency(sourceCurrency)
	targetCurrency := req.FormValue("target")
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
