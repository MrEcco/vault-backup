package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// Restore func
func Restore(Address string, Token string, Backup BackupHandler) bool {

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

	// KVs
	for _, v := range Backup.KVBuckets {
		CreateBucket(vaultClient, v.Name)
		for _, vv := range v.Entries {
			WriteRecord(
				vaultClient,
				v.Name,
				vv.Path,
				vv.Content,
			)
		}
	}

	// Policies
	for _, v := range Backup.Policies {
		CreatePolicy(
			vaultClient,
			v.Name,
			v.ACL,
		)
	}

	return true
}

// CreateBucket func
func CreateBucket(vc *vault.Client, Name string) bool {
	// Sanityze
	Name = strings.Trim(Name, "/")

	_, errList := vc.Logical().Write(
		// Path
		"sys/mounts/"+Name,
		// Data
		map[string]interface{}{
			"path":   Name,
			"type":   "kv",
			"config": map[string]interface{}{}, // Empty untyped map
			"options": map[string]interface{}{
				"version": 2,
			},
			"generate_signing_key": true,
		},
	)
	if errList != nil {
		log.Printf("[ERROR] Cannot create bucket \"%s\": %s", Name, errList.Error())
		return false
	}
	return true
}

// WriteRecord func
func WriteRecord(vc *vault.Client, Bucket string, Path string, Content map[string]interface{}) bool {
	// Sanityze
	Path = strings.Trim(Path, "/")
	Bucket = strings.Trim(Bucket, "/")

	// Create bucket if not exist
	log.Printf("[DEBUG] Write record \"%s/%s\"", Bucket, Path)
	_, errWrite := vc.Logical().Write(
		// Path
		Bucket+"/data/"+Path,
		// Data
		map[string]interface{}{
			"data": Content,
			"options": map[string]interface{}{
				"cas": 0,
			},
		},
	)
	if errWrite != nil {
		log.Printf("[ERROR] Cannot write record \"%s/%s\": %s", Bucket, Path, errWrite.Error())
		return false
	}
	return true
}

// CreatePolicy func
func CreatePolicy(vc *vault.Client, Name string, Policy string) bool {
	// Sanityze
	Name = strings.Trim(Name, "/")

	// Create bucket if not exist
	log.Printf("[DEBUG] Create policy \"%s\"", Name)
	_, errWrite := vc.Logical().Write(
		// Path
		"sys/policies/acl/"+Name,
		// Data
		map[string]interface{}{
			"name":   Name,
			"policy": Policy,
		},
	)
	if errWrite != nil {
		log.Printf("[ERROR] Cannot create policy \"%s\": %s", Name, errWrite.Error())
		return false
	}
	return true
}
