package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goombaio/namegenerator"
	cpf "github.com/mvrilo/go-cpf"
	"github.com/starkbank/sdk-go/starkbank/invoice"
)

const (
	maxInvoiceAmount = 9999999999
)

func issueInvoices() {
	// Random number of invoices from 8 to 12
	numberOfInvoices := 8 + rand.Intn(5)
	invoices := make([]invoice.Invoice, numberOfInvoices)
	for i := 0; i < numberOfInvoices; i++ {
		invoices[i] = generateRandomInvoice()
	}

	invoices, err := invoice.Create(invoices, nil)

	if err.Errors != nil {
		for _, e := range err.Errors {
			panic(fmt.Sprintf("code: %s, message: %s", e.Code, e.Message))
		}
	}

	fmt.Printf(
		"Issued %d invoices at: %v\n",
		numberOfInvoices,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

func generateRandomInvoice() invoice.Invoice {
	return invoice.Invoice{
		Amount: rand.Intn(maxInvoiceAmount),
		Name:   namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate(),
		TaxId:  cpf.Generate(),
	}
}
