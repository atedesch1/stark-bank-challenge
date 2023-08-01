package main

import (
	"fmt"
	"time"

	"github.com/starkbank/sdk-go/starkbank/transfer"
)

func createTransferToStarkBank(amount int) transfer.Transfer {
	return transfer.Transfer{
		Amount:        amount,
		Name:          "Stark Bank S.A",
		TaxId:         "20.018.183/0001-80",
		BankCode:      "20018183",
		BranchCode:    "0001",
		AccountNumber: "6341320293482496",
		AccountType:   "payment",
	}
}

func transferAmountToStarkBank(amount int) {
	_, err := transfer.Create(
		[]transfer.Transfer{
			createTransferToStarkBank(amount),
		}, nil)

	if err.Errors != nil {
		for _, e := range err.Errors {
			panic(fmt.Sprintf("code: %s, message: %s", e.Code, e.Message))
		}
	}

	fmt.Printf(
		"Tranfered %d to Stark Bank at: %v\n",
		amount,
		time.Now().Format("2000-01-01 01:01:01"),
	)
}
