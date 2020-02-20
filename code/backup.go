package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// Backup func
func Backup(Address string, Token string, KVs []string) BackupHandler {

	// We do it in single thread for avoid service overload

	ret := BackupHandler{}

	vaultClient, errNewClient := vault.NewClient(
		&vault.Config{
			Address: Address,
			HttpClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	)
	if errNewClient != nil {
		panic(errNewClient)
	}

	vaultClient.SetToken(Token)

	for _, v := range KVs {
		entries, errDump := DumpItems(vaultClient, v)
		if errDump != nil {
			panic(errDump)
		}

		ret.KVBuckets = append(make([]struct {
			Name    string        `yaml:"name"`
			Entries []ItemHandler `yaml:"entries"`
		}, 0), struct {
			Name    string        `yaml:"name"`
			Entries []ItemHandler `yaml:"entries"`
		}{
			Name:    "mykv",
			Entries: entries,
		})
	}

	pols, errDP := DumpPolicies(vaultClient)
	if errDP != nil {
		panic(errDP)
	}

	ret.Policies = append(
		make([]PolicyHandler, 0),
		pols...,
	)

	return ret
}

// DumpItems func
func DumpItems(vc *vault.Client, kvname string) ([]ItemHandler, error) {
	ret := make([]ItemHandler, 0)

	tree, errTree := ListTree(
		vc,
		kvname,
		"",
	)
	if errTree != nil {
		return ret, errTree
	}

	for _, v := range tree {
		data, errRead := ReadValue(vc, kvname, v)
		if errRead != nil {
			return ret, errRead
		}
		if len(data) != 0 {
			ret = append(ret, ItemHandler{
				Path:    v,
				Content: data,
			})
		}
	}

	return ret, nil
}

// ListTree just recursively search all values in specified kv and return founded item paths
func ListTree(vc *vault.Client, kvname string, path string) ([]string, error) {
	retList := make([]string, 0)

	log.Printf("[DEBUG] Listing \"%s/%s\"", kvname, strings.Trim(path, "/"))
	data, errList := vc.Logical().List(kvname + "/metadata/" + strings.Trim(path, "/"))
	if errList != nil {
		log.Printf("[ERROR] Cannot list \"%s/%s\" : %s", kvname, path, errList.Error())
		return retList, fmt.Errorf("Cannot list items in %s: %s", kvname, errList.Error())
	}

	if data == nil {
		return retList, nil
	}

	keys := (*data).Data["keys"].([]interface{})

	for _, v := range keys {
		entity := path + "/" + strings.Trim(v.(string), "/")

		retList = append(
			retList,
			entity,
		)

		lr, errList := ListTree(vc, kvname, entity)
		if errList != nil {
			return retList, errList
		}

		for _, vv := range lr {
			retList = append(
				retList,
				vv,
			)
		}
	}

	return retList, nil
}

// ReadValue func
func ReadValue(vc *vault.Client, kvname string, path string) (map[string]interface{}, error) {

	log.Printf("[DEBUG] Read \"%s/%s\" value...", kvname, strings.Trim(path, "/"))
	data, errRead := vc.Logical().Read(kvname + "/data/" + strings.Trim(path, "/"))
	if errRead != nil {
		// log.Printf("[ERROR] Cannot read bucket \"%s\" value \"%s\": %s", kvname, path, errRead.Error())
		return make(map[string]interface{}), fmt.Errorf("Cannot read bucket \"%s\" value \"%s\": %s", kvname, path, errRead.Error())
	}

	if data == nil {
		return make(map[string]interface{}), nil // empty record
	}

	return (*data).Data["data"].(map[string]interface{}), nil
}

// DumpPolicies func
func DumpPolicies(vc *vault.Client) ([]PolicyHandler, error) {
	ret := make([]PolicyHandler, 0)

	log.Printf("[DEBUG] List policies")
	list, errList := vc.Logical().List("sys/policies/acl")
	if errList != nil {
		log.Printf("[ERROR] Cannot list policies: %s", errList.Error())
		return ret, fmt.Errorf("Cannot list policies: %s", errList.Error())
	}

	if list == nil {
		return ret, nil
	}

	for _, p := range (*list).Data["keys"].([]interface{}) {
		policyName := p.(string)

		// Skip root policies
		if policyName == "root" {
			continue
		}

		// Skip default policies
		if policyName == "default" {
			continue
		}

		log.Printf("[DEBUG] Read policy \"%s\"", policyName)
		acl, errRead := vc.Logical().Read("sys/policies/acl/" + policyName)
		if errRead != nil {
			log.Printf("[ERROR] Cannot read \"%s\" policy ACL: %s", policyName, errRead.Error())
		}

		policyACL := (*acl).Data["policy"].(string)

		ret = append(ret, PolicyHandler{
			Name: policyName,
			ACL:  policyACL,
		})
	}

	return ret, nil
}
