package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"nasdaq-grabber/controllers"
	"nasdaq-grabber/middleware"
	_ "nasdaq-grabber/models"
	"nasdaq-grabber/sse"
	"nasdaq-grabber/types"

	"github.com/gorilla/mux"
)

var companies = [...]types.PopularQuote{
	{Symbol: "FB", CompanyName: "Facebook, Inc."},
	{Symbol: "AMZN", CompanyName: "Amazon.com, Inc."},
	{Symbol: "MSFT", CompanyName: "Microsoft Corporation"},
	{Symbol: "AAPL", CompanyName: "Apple Inc."},
	{Symbol: "GOOGL", CompanyName: "Alphabet Inc."},
}

func main() {

	router := mux.NewRouter()

	broker := sse.NewServer()

	router.Handle("/stocks", broker)

	router.HandleFunc("/quotes", middleware.AuthMiddleware(QuotesHandler)).Methods("GET")

	// router.HandleFunc("/quotes", QuotesPostHandler).Methods("POST")

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 5)
	// 		eventString := fmt.Sprintf("the time is %v", time.Now())
	// 		log.Println("Receiving event")
	// 		broker.Notifier <- []byte(eventString)
	// 	}
	// }()

	router.HandleFunc("/signin", controllers.SignIn).Methods("POST")
	router.HandleFunc("/login", controllers.LogIn).Methods("POST")
	router.HandleFunc("/refresh", controllers.Refresh).Methods("POST")

	go func() {
		done := make(chan bool)

		for n := range sendData(getData("https://api.nasdaq.com/api/quote/MSFT/info?assetclass=stocks", done)) {
			json, _ := json.Marshal(n)
			log.Println("Receiving event")
			broker.Notifier <- json
		}

		defer func() {
			done <- true
			close(done)
		}()
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", router))
}

var QuotesHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(companies)
	if err != nil {
		log.Fatal(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Write(json)
})

var QuotesPostHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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
})

// https://api.nasdaq.com/api/quote/MSFT/info?assetclass=stocks
func getQuotes(url string) types.Data {
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

	return responseData.Data
}

func getData(url string, done <-chan bool) <-chan types.Data {
	out := make(chan types.Data)
	go func() {
		defer close(out)
		for {
			time.Sleep(time.Second * 10)
			select {
			case <-done:
				return
			default:
				out <- getQuotes(url)
			}
		}
	}()
	return out
}

func sendData(in <-chan types.Data) <-chan types.Data {
	out := make(chan types.Data)
	go func() {
		defer close(out)
		for n := range in {
			out <- n
		}
		return
	}()
	return out
}
