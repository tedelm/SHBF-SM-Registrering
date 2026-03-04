package clientsession

import (
	"net/http"
	"net/http/cookiejar"
)

type ClientSession struct {
	Client *http.Client
}

func NewClientSession() (*ClientSession, error) {
	// Create cookie jar (equivalent to WebRequestSession)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &ClientSession{
		Client: &http.Client{
			Jar: jar,
		},
	}, nil
}
