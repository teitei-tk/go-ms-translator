package translator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	OAuthFetchTokenURL        = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	MicrosoftTranslatorAPIURL = "https://api.microsofttranslator.com/v2/http.svc/Translate"
)

type Response struct {
	Content string `xml:",chardata"`
}

func TranslateRequest(subscriptionKey, text, from, to string) (string, error) {
	token, err := getAccessToken(subscriptionKey)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, MicrosoftTranslatorAPIURL, nil)
	if err != nil {
		return "", nil
	}

	v := &url.Values{}
	v.Add("text", text)
	v.Add("from", from)
	v.Add("to", to)

	req.URL.RawQuery = v.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)

	response := Response{}
	if err = decoder.Decode(&response); err != nil {
		return "", nil
	}
	return response.Content, nil
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
