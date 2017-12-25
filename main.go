package main

import (
	"crypto/data/postgres"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	gdax *GdaxClient
	pg   *postgres.Postgres
	btc  = make([]*postgres.Alert, 0)
	eth  = make([]*postgres.Alert, 0)
	ltc  = make([]*postgres.Alert, 0)
)
var alerts = map[string]float64{
	"BTC-USD": 14000,
	"ETH-USD": 700,
	"LTC-USD": 275,
}

func init() {
	gdax = NewGdaxClient(nil)
	pg = postgres.NewPostgres("postgres", "", "crypto", "10.0.1.193", "5432", "", "")
}

func status(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `{"status":"ok"}`)
}

func main() {
	defer pg.Close() // always close postgres connection before exiting program
	alerts, err := pg.GetAlerts()
	if err != nil {
		log.Fatalf("Failed to fetch alerts: %s", err)
	}
	for _, alert := range alerts {
		switch alert.CurrencyID {
		case 1:

		}
	}
	http.HandleFunc("/status", status) // set router
	log.Println("main serving")

	go checkBitcoin()
	go checkLitecoin()
	go checkEtherium()
	err = http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type GdaxClient struct {
	Base   string //https://api.gdax.com
	Client *http.Client
}

func checkEtherium() {
	ticker := time.Tick(10 * time.Second)
	alerting := false
	for {
		select {
		case <-ticker:
			resp, err := gdax.Fetch("ETH-USD")
			if err != nil {
				break
			}
			if len(resp.Bids) == 0 {
				// we don't have any info. wait until next check
				break
			}
			var current float64
			// use reflection to see if the first element of bids is a string
			if currPrice, ok := resp.Bids[0][0].(string); ok {
				current, _ = strconv.ParseFloat(currPrice, 64)
			}

			if !alerting && current > 0 && current <= alerts["ETH-USD"] {
				log.Printf("Etherium Alert! Etherium is less than %f (%v)", alerts["ETH-USD"], current)
				alerting = true
				break
			} else if current > alerts["ETH-USD"] {
				alerting = false
			}
		}
	}
}

func checkLitecoin() {
	ticker := time.Tick(10 * time.Second)
	alerting := false
	for {
		select {
		case <-ticker:
			resp, err := gdax.Fetch("LTC-USD")
			if err != nil {
				break
			}
			if len(resp.Bids) == 0 {
				// we don't have any info. wait until next check
				break
			}
			var current float64
			// use reflection to see if the first element of bids is a string
			if currPrice, ok := resp.Bids[0][0].(string); ok {
				current, _ = strconv.ParseFloat(currPrice, 64)
			}

			if !alerting && current > 0 && current <= alerts["LTC-USD"] {
				log.Printf("Litecoin Alert! Litecoin is less than %f (%v)", alerts["LTC-USD"], current)
				alerting = true
				break
			} else if current > alerts["LTC-USD"] {
				alerting = false
			}
		}
	}
}

func checkBitcoin() {
	ticker := time.Tick(10 * time.Second)
	alerting := false
	for {
		select {
		case <-ticker:
			resp, err := gdax.Fetch("BTC-USD")
			if err != nil {
				break
			}
			if len(resp.Bids) == 0 {
				// we don't have any info. wait until next check
				break
			}
			var current float64
			// use reflection to see if the first element of bids is a string
			if currPrice, ok := resp.Bids[0][0].(string); ok {
				current, _ = strconv.ParseFloat(currPrice, 64)
			}

			if !alerting && current > 0 && current <= alerts["BTC-USD"] {
				log.Printf("Bitcoin Alert! Bitcoin is less than %f (%v)", alerts["BTC-USD"], current)
				alerting = true
				break
			} else if current > alerts["BTC-USD"] {
				alerting = false
			}
		}
	}
}

func NewGdaxClient(client *http.Client) *GdaxClient {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &GdaxClient{
		Base:   "https://api.gdax.com",
		Client: client,
	}
}

func (g *GdaxClient) Fetch(product string) (*GdaxResponse, error) {
	url := fmt.Sprintf("%s/products/%s/book", g.Base, product)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to build request: %v", err)
		return nil, err
	}
	resp, err := g.Client.Do(req)
	if err != nil {
		log.Printf("Error requesting %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
	}
	gresp := &GdaxResponse{}
	err = json.Unmarshal(body, gresp)
	if err != nil {
		log.Printf("Failed to unmarshal gdax response: %v", err)
	}

	return gresp, nil
}

type GdaxResponse struct {
	Sequence int             `json:"sequence"`
	Bids     [][]interface{} `json:"bids"`
	Asks     [][]interface{} `json:"asks"`
}

//{"sequence":4655738901,"bids":[["14080.13","5.81923144",8]],"asks":[["14080.14","2.89340088",1]]}
