package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	addr := flag.String("addr", "", "Address of vault service.")
	authtype := flag.String("authtype", "kubejwt", "Method of token asquition. \"token\" read VAULT_TOKEN env var. \"kubejwt\" use Kubernetes vault plugin and default way.")
	token := flag.String("token", "", "Token for work with vault service.")
	jwtpath := flag.String("jwtpath", "/var/run/secrets/kubernetes.io/serviceaccount/token", "Custom path to Kubernetes ServiceAccount JWT file. Useless for any non-\"kubejwt\" auth methods.")
	jwtrole := flag.String("jwtrole", "default", "Custom vault role to assume via Kubernetes ServiceAccount JWT. Useless for any non-\"kubejwt\" auth methods.")
	command := ""
	flag.Parse()

	// Check arguments
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("You must specify command!")
	}
	switch args[0] {
	case "backup":
		command = args[0]
		args = args[1:len(args)] // shift
		if len(args) < 1 {
			log.Printf("No one KVs to backup.")
		}
	case "restore":
		command = args[0]
		args = args[1:len(args)] // shift
		if len(args) < 1 {
			log.Fatalf("No one backup file specified to restore.")
		}
	default:
		log.Fatalf("Only \"backup\" and \"restore\" commands available.")
	}

	// Common args
	if len(*addr) < 7 {
		log.Fatalf("Bad vault service address: %s", *addr)
	}
	if (*addr)[0:4] != "http" {
		log.Fatalf("Bad vault service address: %s", *addr)
	}
	switch *authtype {
	case "kubejwt":
	case "token":
		if *token == "" {
			log.Fatalf("Bad vault token: %s", *token)
		}
	default:
		log.Fatalf("Unexpected vault auth method: %s", *authtype)
	}

	// Auth
	if *authtype == "kubejwt" {
		var err error
		*token, err = TokenByJWT( // for any error it crash here, lol
			*addr,
			*jwtpath,
			*jwtrole,
		)
		if err != nil {
			panic(err)
		}
	}

	switch command {
	case "backup":
		fmt.Printf("%s\n", YAMLizeStruct(
			Backup(
				*addr,
				*token,
				args,
			),
		))
	case "restore": // I fucking scare to use it in production!
		var handler BackupHandler

		yamlBytes, errReadFile := ioutil.ReadFile(args[0])
		if errReadFile != nil {
			log.Fatalf("Error to read backup file: %s\n", errReadFile.Error())
		}

		if errUnmarshal := yaml.Unmarshal([]byte(yamlBytes), &handler); errUnmarshal != nil {
			log.Fatalf("Error to parse backup file: %s\n", errUnmarshal.Error())
		}

		Restore(*addr, *token, handler)
	}

}
