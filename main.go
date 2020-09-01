package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

func updateTokenHelper(tokens map[string]string, tokenHelper string) error {
	tokenJ, _ := json.Marshal(tokens)
	err := ioutil.WriteFile(tokenHelper, tokenJ, 0644)
	return err
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Requires one arg of 'get', 'erase', or 'store'")
	}

	tokenHelper, err := homedir.Expand("~/.vault-tokens")
	if err != nil {
		log.Fatal(err)
	}

	var tokenData map[string]string

	if _, err := os.Stat(tokenHelper); err == nil {
		data, err := ioutil.ReadFile(tokenHelper)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(data, &tokenData)
	}

	vault_addr := os.Getenv("VAULT_ADDR")
	if len(vault_addr) == 0 {
		log.Fatal("Environment variable $VAULT_ADDR must be set")
	}

	switch os.Args[1] {
	case "get":
		// get
		token := tokenData[vault_addr]
		fmt.Println(token)
	case "erase":
		// erase
		delete(tokenData, vault_addr)
		err := updateTokenHelper(tokenData, tokenHelper)
		if err != nil {
			log.Fatal(err)
		}
	case "store":
		// store
		var input string
		fmt.Scanln(&input)
		tokenData[vault_addr] = input
		err := updateTokenHelper(tokenData, tokenHelper)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Requires one arg of 'get', 'erase', or 'store'")
	}
}
