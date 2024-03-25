package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/baazaouihamza/http-login-tests/pkg/api"
)

func main() {
	var (
		requestURL string
		password   string
	)

	flag.StringVar(&requestURL, "url", "", "target site url")
	flag.StringVar(&password, "password", "", "use a password to access our api")

	flag.Parse()

	parserdURL, err := url.ParseRequestURI(requestURL)

	if err != nil {
		fmt.Printf("validation error: URL is not valid: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	apiInstance := api.New(api.Options{
		Password: password,
		LoginUrl: parserdURL.Scheme + "://" + parserdURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parserdURL.String())
	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response %s\n", res.GetResponse())
}
