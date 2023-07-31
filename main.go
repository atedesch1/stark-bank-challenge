package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/starkbank/sdk-go/starkbank"
	"github.com/starkinfra/core-go/starkcore/user/project"
)

func setupIntegration() error {
	if os.Getenv("ENVIRONMENT") == "dev" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	user := project.Project{
		Id:          os.Getenv("PROJECT_ID"),
		PrivateKey:  os.Getenv("PRIVATE_KEY"),
		Environment: os.Getenv("ENV"),
	}

	starkbank.User = user

	return nil
}

func main() {
	if err := setupIntegration(); err != nil {
		panic(err)
	}

	// run webhook listener on go routine
	go serveWebHookServer()

	// run cron job on this go routine for 24hrs
	c := cron.New()
	c.AddFunc("0 */3 * * *", issueInvoices)
	c.Start()
	time.Sleep(24 * time.Hour)
	c.Stop()
}

func serveWebHookServer() {
	http.HandleFunc("/invoicehook", invoiceHookHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(fmt.Sprintln("Error starting the server: ", err))
	}
}

func invoiceHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invoice webhook should have POST method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	reqBody := string(body)

	digitalSignature, ok := r.Header["Digital-Signature"]
	if !ok || len(digitalSignature) == 0 {
		http.Error(w, "No digital signature found", http.StatusBadRequest)
		return
	}

	err = verifyDigitalSignature(digitalSignature[0], reqBody)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintln("Could not verify digital signature: ", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	type requestBody struct {
		Event struct {
			Subscription string `json:"subscription"`
			Log          struct {
				Invoice struct {
					Amount int `json:"amount"`
				} `json:"invoice"`
			} `json:"log"`
		} `json:"event"`
	}

	var req requestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintln("Invalid request payload: ", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	event := req.Event

	if event.Subscription != "invoice" {
		http.Error(w, "Invalid subscription", http.StatusBadRequest)
		return
	} else if event.Log.Invoice.Amount <= 0 {
		http.Error(w, "Amount from invoice should be positive", http.StatusBadRequest)
		return
	}

	transferAmountToStarkBank(event.Log.Invoice.Amount)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Transfer concluded")
}
