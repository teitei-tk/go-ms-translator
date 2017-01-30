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
	"strings"
	"time"
)

const (
	OAuthFetchTokenURL   = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	TranslateAPIURL      = "https://api.microsofttranslator.com/v2/http.svc/Translate"
	TranslateArrayAPIURL = "https://api.microsofttranslator.com/v2/http.svc/TranslateArray"
)

var userAgent = fmt.Sprintf("XXXGoClient/ (%s)", runtime.Version())

type (
	TranslateArrayResponse struct {
		From           string
		TranslatedText string
	}

	ArrayOfTranslateArrayResponse struct {
		TranslateArrayResponse []TranslateArrayResponse
	}
)

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
	v.Add("from", from)
	v.Add("to", to)

	req, err := c.newRequest(ctx, v.Encode(), http.MethodGet, from, to)
	if err != nil {
		return "", err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errResponse(res)
	}

	type apiResponse struct {
		Content string `xml:",chardata"`
	}

	response := apiResponse{}
	if err = decodeXML(res, &response); err != nil {
		return "", err
	}

	return response.Content, nil
}

func TranslateArray(subscriptionKey string, texts []string, from, to string) ([]string, error) {
	var result []string
	c, err := NewClient(TranslateArrayAPIURL, nil)
	if err != nil {
		return result, err
	}

	token, err := getAccessToken(subscriptionKey)
	if err != nil {
		return result, err
	}
	c.Token = token

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var reqText []string
	for _, v := range texts {
		s := `<string xmlns='http://schemas.microsoft.com/2003/10/Serialization/Arrays'>%s</string>`
		reqText = append(reqText, fmt.Sprintf(s, v))
	}

	reqBody := fmt.Sprintf(`<TranslateArrayRequest>
                <AppId />
                <From>%s</From>
                <Options>
                        <Category xmlns='http://schemas.datacontract.org/2004/07/Microsoft.MT.Web.Service.V2' />
                        <ContentType>text/plain</ContentType>
                        <ReservedFlags xmlns='http://schemas.datacontract.org/2004/07/Microsoft.MT.Web.Service.V2' />
                        <State xmlns='http://schemas.datacontract.org/2004/07/Microsoft.MT.Web.Service.V2' />
                        <Uri xmlns='http://schemas.datacontract.org/2004/07/Microsoft.MT.Web.Service.V2' />
                        <User xmlns='http://schemas.datacontract.org/2004/07/Microsoft.MT.Web.Service.V2' />
                </Options>
                <Texts>%s</Texts>
                <To>%s</To>
        </TranslateArrayRequest>`, from, strings.Join(reqText, " "), to)

	req, err := c.newRequest(ctx, reqBody, http.MethodPost, from, to)
	req.Header.Set("Content-Type", "text/xml")
	if err != nil {
		return result, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	response := ArrayOfTranslateArrayResponse{}
	if err = xml.Unmarshal(b, &response); err != nil {
		return result, err
	}

	for _, v := range response.TranslateArrayResponse {
		result = append(result, v.TranslatedText)
	}

	return result, nil
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
	req.Header.Set("User-Agent", userAgent)

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
