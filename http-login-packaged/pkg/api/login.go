package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, loginURL, password string) (string, error) {

	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %s", err)
	}

	response, err := client.Post(loginURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("http Post error %s", err)
	}

	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("invalid output (HTTP Code %d): %s", response.StatusCode, string(resBody))
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      "No valid JSON returned",
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      "Empty token replied",
		}
	}

	return loginResponse.Token, nil
}
