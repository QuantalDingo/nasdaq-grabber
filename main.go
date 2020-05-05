package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"nasdaq-grabber/sse"
	"nasdaq-grabber/types"

	"github.com/gorilla/mux"
)

func main() {
	companies := [...]types.PopularQuote{{Symbol: "FB", CompanyName: "Facebook, Inc."},
		{Symbol: "AMZN", CompanyName: "Amazon.com, Inc."},
		{Symbol: "MSFT", CompanyName: "Microsoft Corporation"},
		{Symbol: "AAPL", CompanyName: "Apple Inc."},
		{Symbol: "GOOGL", CompanyName: "Alphabet Inc."}}

	router := mux.NewRouter()

	broker := sse.NewServer()

	router.Handle("/stocks", broker)

	router.HandleFunc("/quotes", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var quotes []string
		err = json.Unmarshal(body, &quotes)
		if err != nil {
			log.Fatal(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(quotes)
	}).Methods("POST")

	router.HandleFunc("/quotes", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		json, err := json.Marshal(companies)
		if err != nil {
			log.Fatal(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Write(json)
	}).Methods("GET")

	go func() {
		for {
			time.Sleep(time.Second * 5)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", router))
}

func getQuotes(url string) (response types.ResponseData) {
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

	var responseData types.ResponseData

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%#v\n", responseData)

	return responseData
}
