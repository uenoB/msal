//go:build !simple

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func main() {
	var (
		refresh = flag.Bool("R", false, "force refresh")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr,
			    "usage: %s [-R] config.json\n",
			    os.Args[0])
		os.Exit(1)
	}

	config, err := OpenConfig(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	if (*refresh) {
		config.Token = nil
	}

	client, err := public.New(config.ClientID, public.WithCache(config))
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := client.Accounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var result public.AuthResult
	if len(accounts) > 0 {
		result, err = client.AcquireTokenSilent(
			context.Background(),
			config.Scopes,
			public.WithSilentAccount(accounts[0]),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		reqUrl, err := client.AuthCodeURL(
			context.Background(),
			config.ClientID,
			config.RedirectURI,
			config.Scopes,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Access this URL by your browser:")
		fmt.Println(reqUrl)
		fmt.Println("Paste redirected URL:")

		var resUrl string
		fmt.Scan(&resUrl)
		u, err := url.Parse(resUrl)
		if err != nil {
			log.Fatal(err)
		}
		authCode, found := u.Query()["code"]
		if !found {
			log.Fatal("auth code not found")
		}

		result, err = client.AcquireTokenByAuthCode(
			context.Background(),
			authCode[0],
			config.RedirectURI,
			config.Scopes,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(result.AccessToken)
}
