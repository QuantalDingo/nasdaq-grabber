package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PrimaryData struct {
	LastSalePrice      string `json:"lastSalePrice"`
	NetChange          string `json:"netChange"`
	PercentageChange   string `json:"percentageChange"`
	DeltaIndicator     string `json:"deltaIndicator"`
	LastTradeTimestamp string `json:"lastTradeTimestamp"`
	IsRealTime         bool   `json:"isRealTime"`
}

type SecondaryData struct {
	LastSalePrice      string `json:"lastSalePrice"`
	NetChange          string `json:"netChange"`
	PercentageChange   string `json:"percentageChange"`
	DeltaIndicator     string `json:"deltaIndicator"`
	LastTradeTimestamp string `json:"lastTradeTimestamp"`
	IsRealTime         bool   `json:"isRealTime"`
}

type Volume struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type PreviousClose struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type OpenPrice struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type MarketCap struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type KeyStats struct {
	Volume        Volume        `json:"Volume"`
	PreviousClose PreviousClose `json:"PreviousClose"`
	OpenPrice     OpenPrice     `json:"OpenPrice"`
	MarketCap     MarketCap     `json:"MarketCap"`
}

type Data struct {
	Symbol         string        `json:"symbol"`
	CompanyName    string        `json:"companyName"`
	StockType      string        `json:"stockType"`
	Exchange       string        `json:"exchange"`
	IsNasdaqListed bool          `json:"isNasdaqListed"`
	IsNasdaq100    bool          `json:"isNasdaq100"`
	IsHeld         bool          `json:"isHeld"`
	PrimaryData    PrimaryData   `json:"primaryData"`
	SecondaryData  SecondaryData `json:"secondaryData"`
	KeyStats       KeyStats      `json:"keyStats"`
	MarketStatus   string        `json:"marketStatus"`
	AssetClass     string        `json:"assetClass"`
}

type ResponseData struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Status  Status `json:"status"`
}

type Status struct {
	RCode            int    `json:"rCode"`
	BCodeMessage     string `json:"bCodeMessage"`
	DeveloperMessage string `json:"developerMessage"`
}

func main() {
	parsePage("https://api.nasdaq.com/api/quote/AMZN/info?assetclass=stocks")
}

func parsePage(url string) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseData ResponseData

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", responseData)
}
