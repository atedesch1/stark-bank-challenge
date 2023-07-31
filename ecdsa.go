package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/ecdsa"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/publickey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/signature"
)

func verifyDigitalSignature(header http.Header, message string) error {
	digitalSignature, ok := header["Digital-Signature"]
	if !ok || len(digitalSignature) == 0 {
		return errors.New("no digital signature found")
	}

	signature := signature.FromBase64(digitalSignature[0])

	publicKey, err := getStarkPublicKey()
	if err != nil {
		return err
	}

	pk := publickey.FromPem(publicKey)

	if !ecdsa.Verify(
		message,
		signature,
		&pk,
	) {
		return errors.New("could not verify incoming webhook digital signature")
	}

	return nil
}

func getStarkPublicKey() (string, error) {
	publicKeyResponse, err := http.Get("https://sandbox.api.starkbank.com/v2/public-key")
	if err != nil {
		return "", err
	}

	type responseBody struct {
		PublicKeys []struct {
			Content string `json:"content"`
		} `json:"publicKeys"`
	}

	var res responseBody

	err = json.NewDecoder(publicKeyResponse.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	if len(res.PublicKeys) == 0 {
		return "", fmt.Errorf("No public keys were found")
	}
	return res.PublicKeys[0].Content, nil
}
