package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/ecdsa"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/publickey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/signature"
)

func verifyDigitalSignature(digitalSignature string, message string) error {
	log.Println("signature in base64", digitalSignature)
	signature := signature.FromBase64(digitalSignature)
	log.Println("signature from base64", signature)

	publicKey, err := getStarkPublicKey()
	if err != nil {
		return err
	}
	log.Println("publicKey", publicKey)

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
		return "", errors.New("No public keys were found")
	}
	return res.PublicKeys[0].Content, nil
}
