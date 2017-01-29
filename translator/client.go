package translator

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
	Token      string
}

func NewClient(rawURL string, logger *log.Logger) (*Client, error) {
	url, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	var discardLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	c := &Client{
		URL:        url,
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}

	return c, nil
}

func newHTTPRequest(method, url, q string) (*http.Request, error) {
	if method == http.MethodPost {
		return http.NewRequest(method, url, strings.NewReader(q))
	}

	req, err := http.NewRequest(method, url, nil)
	req.URL.RawQuery = q
	return req, err
}

func (c *Client) newRequest(ctx context.Context, q, method, from, to string) (*http.Request, error) {
	req, err := newHTTPRequest(method, c.URL.String(), q)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", userAgent)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Accept", "application/xml")

	return req, nil
}
