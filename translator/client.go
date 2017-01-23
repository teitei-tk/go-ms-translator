package translator

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func (c *Client) newRequest(ctx context.Context, v *url.Values, method, from, to string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	v.Add("from", from)
	v.Add("to", to)

	req = req.WithContext(ctx)
	req.URL.RawQuery = v.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/xml")

	return req, nil
}
