package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	GetResponseOutPut  *http.Response
	PostResponseOutput *http.Response
}

func (m MockClient) Get(url string) (resp *http.Response, err error) {
	return m.GetResponseOutPut, nil
}
func (m MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return m.PostResponseOutput, err
}
func TestDoGetRequest(t *testing.T) {
	words := WordsPage{
		Page:  Page{Name: "words"},
		Words: Words{Input: "abc", Words: []string{"a", "b"}},
	}
	wordsBytes, err := json.Marshal(words)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}
	apiInstance := API{
		Options: Options{},
		Client: MockClient{
			GetResponseOutPut: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
			},
		},
	}

	response, err := apiInstance.DoGetRequest("http://localhost/words")
	if err != nil {
		t.Errorf("do get request error: %s", err)
	}

	if response == nil {
		t.Fatalf("response is empty")
	}

	if response.GetResponse() != strings.Join([]string{"a", "b"}, ", ") {
		t.Errorf("Unexpected response %s", response.GetResponse())
	}
}
