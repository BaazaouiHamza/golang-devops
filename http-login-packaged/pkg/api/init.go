package api

import "net/http"

type Options struct {
	Password string
	LoginUrl string
}

type APIIface interface {
	DoGetRequest(requestURL string) (Response, error)
}

type API struct {
	Options Options
	Client  http.Client
}

func New(Options Options) APIIface {
	return API{
		Options: Options,
		Client: http.Client{
			Transport: &MyJWTTransport{
				transport: http.DefaultTransport,
				password:  Options.Password,
				loginURL:  Options.LoginUrl,
			},
		},
	}
}
