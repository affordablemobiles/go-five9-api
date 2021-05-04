package five9

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

var LoggingEnabled bool = false

type Five9APIClient struct {
	token *Token
	cjar  http.CookieJar
}

func NewFive9APIClient(ctx context.Context, username, password string) (*Five9APIClient, error) {
	data := &Five9APIClient{}

	var err error
	data.cjar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("Failed to load cookie jar")
	}

	if LoggingEnabled {
		log.Printf("Performing initial login")
	}
	data.token, err = data.performLogin(ctx, username, password)
	if err != nil {
		return nil, err
	}

	err = data.handleStateChange(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}
