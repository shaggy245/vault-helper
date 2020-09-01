package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

func updateTokenHelper(tokens map[string]string, tokenPath string) error {
	tokenJ, _ := json.Marshal(tokens)
	err := ioutil.WriteFile(tokenPath, tokenJ, 0644)
	return err
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Requires one arg of 'get', 'erase', or 'store'")
	}

	tokenPath, err := homedir.Expand("~/.vault-tokens")
	if err != nil {
		log.Fatal(err)
	}

	tokenData := make(map[string]string)

	data, err := ioutil.ReadFile(tokenPath)
	if err == nil {
		json.Unmarshal(data, &tokenData)
	}

	fmt.Println(tokenData)

	vault_addr := os.Getenv("VAULT_ADDR")
	if vault_addr == "" {
		log.Fatal("Environment variable $VAULT_ADDR must be set")
	}

	switch os.Args[1] {
	case "get":
		// get
		token := tokenData[vault_addr]
		fmt.Print(token)
	case "erase":
		// erase
		delete(tokenData, vault_addr)
		err := updateTokenHelper(tokenData, tokenPath)
		if err != nil {
			log.Fatal(err)
		}
	case "store":
		// store
		var input string
		fmt.Scanln(&input)
		tokenData[vault_addr] = input
		err := updateTokenHelper(tokenData, tokenPath)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Requires one arg of 'get', 'erase', or 'store'")
	}
}
