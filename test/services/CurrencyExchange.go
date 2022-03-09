package services

import (
	"encoding/json"
	"github.com/djfemz/simbaCodingChallenge/util"
	"io/ioutil"
	"log"
	"net/http"
)

type Result struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func GetCurrencyExchangeRate(currency, targetCurrency string) float64 {
	config, err := util.LoadConfig("/home/djfemz/Documents/goworkspace/github.com/simbaCodingChallenge/app.env")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", string("https://v6.exchangerate-api.com/v6/latest/"+currency), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+config.ApiKey)
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Result
	err = json.Unmarshal(bs, &result)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range result.ConversionRates {
		if key == targetCurrency {
			return value
		}
	}
	return 0.0
}
