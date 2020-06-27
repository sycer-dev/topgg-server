package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spec-tacles/go/broker"
)

// VoteData the data recieved from top.gg
type VoteData struct {
	BotID     string `json:"bot"`
	UserID    string `json:"user"`
	Type      string `json:"type"`
	IsWeekend bool   `json:"isWeekend"`
	Query     string `json:"query"`
}

// Bytes converts the data to a byte array
func (d *VoteData) Bytes() []byte {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(d)
	return buffer.Bytes()
}

var (
	endpoint  = os.Getenv("ENDPOINT")
	address   = os.Getenv("ADDRESS")
	secret    = os.Getenv("DBL_SECRET")
	amqpURL   = os.Getenv("AMQP_URL")
	amqpGroup = os.Getenv("AMQP_GROUP")
)

var rabbit *broker.AMQP

func main() {
	validateEnv()

	rabbit = broker.NewAMQP(amqpGroup, "", nil)
	go rabbit.Connect(amqpURL)

	http.HandleFunc(endpoint, handleIncoming())
	log.Printf("webhook server started at %+s%+s", address, endpoint)

	err := http.ListenAndServe(address, nil)

	log.Panicf("server crashed with error: %+v", err)
	os.Exit(1)
}

func handleIncoming() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Forbidden")
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			log.Printf("incoming request from %+s had no authorization header", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")

			return
		}

		if auth != secret {
			log.Printf("incoming request from %+s had a mismatched authorization header", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")

			return
		}

		var payload *VoteData
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)

		if err != nil {
			log.Printf("incoming request from %+s sent malformed vote data", r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")

			log.Println(err)
			return
		}

		rabbit.PublishOptions(broker.PublishOptions{
			Event:   "VOTE",
			Data:    payload.Bytes(),
			Timeout: 2 * time.Minute,
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")

		return
	}
}

func validateEnv() bool {
	if endpoint == "" {
		endpoint = "/webhooks/vote"
	}

	if address == "" {
		address = ":4500"
	}

	if amqpGroup == "" {
		amqpGroup = "votes"
	}

	if amqpURL == "" {
		amqpURL = "amqp://localhost//"
	}

	if secret == "" {
		log.Panicf("required environment variable 'DBL_SECRET' was not provided")
		os.Exit(1)
	}

	return true
}
