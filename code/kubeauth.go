package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// TokenByJWT func
func TokenByJWT(Address, JWTPath, Role string) (string, error) {

	// Read config
	jwtBytes, errReadFile := ioutil.ReadFile(JWTPath)
	if errReadFile != nil {
		log.Fatalf("Error to read JWT token: %s\n", errReadFile.Error())
	}

	// Request master token
	resp, errPost := http.Post(
		Address+"/v1/auth/kubernetes/login",
		"application/json",
		bytes.NewBuffer(
			[]byte(
				JSONizeStruct(struct {
					JWT  string `json:"jwt"`
					Role string `json:"role"`
				}{
					JWT:  string(jwtBytes),
					Role: Role,
				}),
			),
		),
	)
	if errPost != nil {
		log.Fatalf("Error to connect to Vault: %s\n", errPost.Error())

	}

	// Read response
	reply, errRespReadAll := ioutil.ReadAll(resp.Body)
	if errRespReadAll != nil {
		log.Fatalf("Error to read response from Vault: %s\n", errRespReadAll.Error())

	}

	// Check responce
	var VaultResp struct {
		Auth struct {
			Token                string `json:"client_token"`
			LeaseDurationSeconds int    `json:"lease_duration"`
		} `json:"auth"`
		Errors []string `json:"errors"`
	}
	errUnmarshal := json.Unmarshal(reply, &VaultResp)
	if errUnmarshal != nil {
		log.Fatalf("Bad response from Vault: %s\n", errUnmarshal.Error())

	}
	// Check Vault errors
	if len(VaultResp.Errors) > 0 {
		log.Fatalf("Vault responded with errors: %v", VaultResp.Errors)
	}

	return VaultResp.Auth.Token, nil
}
