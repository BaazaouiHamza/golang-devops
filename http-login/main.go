package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return strings.Join(w.Words, ", ")
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}

	return strings.Join(out, ", ")
}

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
	client := http.Client{}
	if password != "" {
		token, err := doLoginRequest(client, parserdURL.Scheme+"://"+parserdURL.Host+"/login", password)
		if err != nil {
			if requestErr, ok := err.(RequestError); ok {
				fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
				os.Exit(1)
			}
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		client.Transport = MyJWTTransport{
			transport: http.DefaultTransport,
			token:     token,
		}
	}

	res, err := doRequest(client, parserdURL.String())
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
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

func doRequest(client http.Client, requestURL string) (Response, error) {

	response, err := client.Get(requestURL)

	if err != nil {
		return nil, fmt.Errorf("http Get error %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll error %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (HTTP Code %d): %s", response.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "No valid JSON returned",
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}

	switch page.Name {

	case "words":
		var words Words

		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words Unmarshal error: %s", err),
			}
		}
		return words, nil

	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurences Unmarshal error: %s", err),
			}
		}
		return occurrence, nil
	}

	return nil, nil

}
