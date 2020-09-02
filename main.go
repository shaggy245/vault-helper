package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func updateTokens(tokens map[string]string, tokenPath string) error {
	tokenJ, _ := json.Marshal(tokens)
	err := ioutil.WriteFile(tokenPath, tokenJ, 0644)
	return err
}

func eraseToken(tokens map[string]string, vaultAddr string, tokenPath string) error {
	delete(tokens, vaultAddr)
	err := updateTokens(tokens, tokenPath)
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

	rawVaultAddr := os.Getenv("VAULT_ADDR")
	vaultAddr := strings.ToLower(strings.TrimSpace(rawVaultAddr))
	vaultAddr = strings.TrimSuffix(vaultAddr, "/")
	vaultAddr = strings.TrimPrefix(vaultAddr, "http://")
	vaultAddr = strings.TrimPrefix(vaultAddr, "https://")
	if vaultAddr == "" {
		log.Fatal("Environment variable $VAULT_ADDR must be set")
	}

	switch os.Args[1] {
	case "get":
		// get
		token := tokenData[vaultAddr]
		fmt.Print(token)
	case "erase":
		// erase
		err := eraseToken(tokenData, vaultAddr, tokenPath)
		if err != nil {
			log.Fatal(err)
		}
	case "store":
		// store
		var rawInput string
		fmt.Scanln(&rawInput)
		input := strings.TrimSpace(rawInput)
		if input == "" {
			err := eraseToken(tokenData, vaultAddr, tokenPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			tokenData[vaultAddr] = input
			err := updateTokens(tokenData, tokenPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatal("Requires one arg of 'get', 'erase', or 'store'")
	}
}
