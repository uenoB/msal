//go:build simple

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s config.json\n", os.Args[0])
		os.Exit(1)
	}

	config, err := OpenConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	client, err := public.New(config.ClientID, public.WithCache(config))
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := client.Accounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts) == 0 {
		log.Fatal("account not found")
	}

	result, err := client.AcquireTokenSilent(
		context.Background(),
		config.Scopes,
		public.WithSilentAccount(accounts[0]),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.AccessToken)

	config.Close()
}
