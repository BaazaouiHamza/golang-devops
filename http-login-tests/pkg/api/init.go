package api

import (
	"io"
	"net/http"
)

type ClientInterface interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Options struct {
	Password string
	LoginUrl string
}

type APIIface interface {
	DoGetRequest(requestURL string) (Response, error)
}

type API struct {
	Options Options
	Client  ClientInterface
}

func New(Options Options) APIIface {
	return API{
		Options: Options,
		Client: &http.Client{
			Transport: &MyJWTTransport{
				transport:  http.DefaultTransport,
				password:   Options.Password,
				loginURL:   Options.LoginUrl,
				HTTPClient: &http.Client{},
			},
		},
	}
}
