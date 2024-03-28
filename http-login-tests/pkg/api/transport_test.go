package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	RoundTripperOutput *http.Response
}

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer abc" {
		return nil, fmt.Errorf("wrong authorization header: %s", req.Header.Get("Authorization"))
	}
	return m.RoundTripperOutput, nil
}

func TestRoundTrip(t *testing.T) {
	loginResponse := LoginResponse{
		Token: "abc",
	}

	loginResponseBytes, err := json.Marshal(loginResponse)
	if err != nil {
		t.Errorf("Marshal error: %s", err)
	}
	myJWTTransport := MyJWTTransport{
		transport: MockRoundTripper{
			RoundTripperOutput: &http.Response{
				StatusCode: 200,
			},
		},
		HTTPClient: MockClient{PostResponseOutput: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
		}},
		password: "xyz",
	}
	req := &http.Request{
		Header: make(http.Header),
	}
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code is not 200, got %d", res.StatusCode)
	}
}
