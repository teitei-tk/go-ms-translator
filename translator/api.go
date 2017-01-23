package translator

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	OAuthFetchTokenURL   = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	TranslateAPIURL      = "https://api.microsofttranslator.com/v2/http.svc/Translate"
	TranslateArrayAPIURL = "https://api.microsofttranslator.com/v2/http.svc/TranslateArray"
)

var userAgent = fmt.Sprintf("XXXGoClient/ (%s)", runtime.Version())

type Response struct {
	Content string `xml:",chardata"`
}

func Translate(subscriptionKey, text, from, to string) (string, error) {
	c, err := NewClient(TranslateAPIURL, nil)
	if err != nil {
		return "", err
	}

	token, err := getAccessToken(subscriptionKey)
	if err != nil {
		return "", err
	}
	c.Token = token

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	v := &url.Values{}
	v.Add("text", text)

	req, err := c.newRequest(ctx, v, http.MethodGet, from, to)
	if err != nil {
		return "", nil
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errResponse(res)
	}

	response := Response{}
	if err = decodeXML(res, &response); err != nil {
		return "", nil
	}

	return response.Content, nil
}

func decodeXML(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	return xml.NewDecoder(res.Body).Decode(out)
}

func decodeJSON(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(out)
}

func errResponse(res *http.Response) error {
	type errRes struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}

	e := &errRes{}
	if err := decodeJSON(res, e); err != nil {
		return err
	}
	return errors.New(e.Message)
}

func getAccessToken(subscriptionKey string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, OAuthFetchTokenURL, nil)
	if err != nil {
		return "", err
	}

	v := &url.Values{}
	v.Add("Subscription-Key", subscriptionKey)
	req.URL.RawQuery = v.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errResponse(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
