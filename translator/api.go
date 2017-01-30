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
	// FetchTokenURL is fetch api token url.
	FetchTokenURL = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"

	// TranslateAPIURL is Text Translation API /translate URI
	TranslateAPIURL = "https://api.microsofttranslator.com/v2/http.svc/Translate"

	// TranslateArrayAPIURL is Text Translation API /translateArray URI
	TranslateArrayAPIURL = "https://api.microsofttranslator.com/v2/http.svc/TranslateArray"
)

// request user agenet.
var userAgent = fmt.Sprintf("XXXGoClient/ (%s)", runtime.Version())

type (
	// TranslateResponse is /translate Response struct.
	// see http://docs.microsofttranslator.com/text-translate.html#!/default/get_Translate
	TranslateResponse struct {
		TranslatedText string `xml:",chardata"`
	}

	// TranslateArrayResponse is each of translated result.
	TranslateArrayResponse struct {
		From           string
		TranslatedText string
	}

	// ArrayOfTranslateArrayResponse is /translateArray Response struct.
	ArrayOfTranslateArrayResponse struct {
		TranslateArrayResponse []TranslateArrayResponse
	}
)

// Fetch Authentication Token API for Microsoft Cognitive Services Translator API
// see http://docs.microsofttranslator.com/oauth-token.html
func getAccessToken(subscriptionKey string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, FetchTokenURL, nil)
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

// Translate is Execute /translate API request. and fetch API response,
// API Reference is http://docs.microsofttranslator.com/text-translate.html#!/default/get_Translate
// return as TranslateResponse struct
func Translate(subscriptionKey, text, from, to string) (*TranslateResponse, error) {
	c, err := NewClient(TranslateAPIURL, nil)
	if err != nil {
		return nil, err
	}

	token, err := getAccessToken(subscriptionKey)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errResponse(res)
	}

	response := &TranslateResponse{}
	if err = decodeXML(res, response); err != nil {
		return nil, err
	}

	return response, nil
}

// TranslateArray is Execute /translateArray API request. and fetch API response,
// API Reference is http://docs.microsofttranslator.com/text-translate.html#!/default/post_TranslateArray
// return as TranslateResponse struct
func TranslateArray(subscriptionKey string, texts []string, from, to string) (*ArrayOfTranslateArrayResponse, error) {
	c, err := NewClient(TranslateArrayAPIURL, nil)
	if err != nil {
		return nil, err
	}

	token, err := getAccessToken(subscriptionKey)
	if err != nil {
		return nil, err
	}
	c.Token = token

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	reqBody := genTranslateArraayReqXML(texts, from, to)
	req, err := c.newRequest(ctx, reqBody, http.MethodPost, from, to)
	req.Header.Set("Content-Type", "text/xml")
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &ArrayOfTranslateArrayResponse{}
	if err = xml.Unmarshal(b, response); err != nil {
		return nil, err
	}

	return response, nil
}

// Generate XML to be used by /translateArray Api.
// sample is https://github.com/MicrosoftTranslator/HTTP-Code-Samples/blob/master/CSharp/TranslateArrayMethod.cs#L48
func genTranslateArraayReqXML(texts []string, from, to string) string {
	var reqText []string
	for _, v := range texts {
		s := `<string xmlns='http://schemas.microsoft.com/2003/10/Serialization/Arrays'>%s</string>`
		reqText = append(reqText, fmt.Sprintf(s, v))
	}

	return fmt.Sprintf(`<TranslateArrayRequest>
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
}

// decode request.body to xml as struct.
func decodeXML(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	return xml.NewDecoder(res.Body).Decode(out)
}

// decode request.body to json as struct.
func decodeJSON(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(out)
}

func errResponse(res *http.Response) error {
	type errRes struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}

	defer res.Body.Close()

	e := &errRes{}
	if err := decodeJSON(res, e); err != nil {
		return err
	}
	return errors.New(e.Message)
}
