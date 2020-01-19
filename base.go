package gohydra

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	GetLoginRequestUrn      = "/oauth2/auth/requests/login?login_challenge="
	AcceptLoginRequestUrn   = "/oauth2/auth/requests/login/accept?login_challenge="
	AcceptConsentRequestUrn = "/oauth2/auth/requests/consent/accept?consent_challenge="
	AcceptLogoutRequestUrn = "/oauth2/auth/requests/logout/accept?logout_challenge="
	Https                   = "https"
	Http                    = "http"
)

func validateUrl(baseUrl string, httpsRequired bool) error {
	if len(baseUrl) == 0 {
		return errors.New("the base url to ory hydra service it is empty")
	}

	baseUrlParsed, parseErr := url.Parse(baseUrl)
	if parseErr != nil {
		return errors.New("the base url informed dont is valid")
	}

	if httpsRequired {
		if !(strings.EqualFold(baseUrlParsed.Scheme, Https)) {
			return errors.New("https protocol is required, but the base url dont is")
		}
	} else {
		if !(strings.EqualFold(baseUrlParsed.Scheme, Http)) {
			return errors.New("http protocol is required, but the base url dont is")
		}
	}

	return nil
}

// Client Encapsule all round trip with a service oryhydra,
// referent the request and challenge in login and consent.
// baseUrl is the localization of oryhydrahost used.
// httpClient is a http client configured to communicate with the oryhydra service.
type Client struct {
	baseUrl    string
	httpClient *http.Client
}

// NewClient create a new valid oryhydra client
func NewClient(baseUrl string, httpsRequired bool, httpClient *http.Client) (*Client, error) {
	validadeBaseUrlsErr := validateUrl(baseUrl, httpsRequired)
	if validadeBaseUrlsErr != nil {
		return nil, validadeBaseUrlsErr
	}

	//TODO(Shineray): Validade if the http client meets the requirements, deep in your configuration...
	if httpClient == nil {
		return nil, errors.New("the http client is invalid to operations")
	}

	return &Client{
		baseUrl:    baseUrl,
		httpClient: httpClient}, nil
}

// uriParse parse the oryhydra url with urn and o challenge id
// for a valid uri.
func (client *Client) challengeUriParse(urn string, challenge string) (*url.URL, error) {

	ChallengeRequestUriRaw := fmt.Sprint(client.baseUrl, urn, challenge)

	ChallengeRequestUri, parseErr := url.ParseRequestURI(ChallengeRequestUriRaw)
	if parseErr != nil {
		return nil, parseErr
	}

	return ChallengeRequestUri, nil
}

// acceptChallengeRequest accept a challenge of the oryhydra, can be login or consent.
// For more information in offcial oryhydra documentation: https://www.ory.sh/docs/next/hydra/sdk/api
func acceptChallengeRequest(client *Client, parameters interface{}, acceptUri *url.URL) (map[string]string, error) {

	parametersSerialized, jsonMarshalErr := json.Marshal(parameters)
	if jsonMarshalErr != nil {
		return nil, jsonMarshalErr
	}

	parametersReader := bytes.NewReader(parametersSerialized)
	acceptLoginRequest, newRequestErr := http.NewRequest("PUT", acceptUri.String(), parametersReader)
	if newRequestErr != nil {
		return nil, newRequestErr
	}

	response, doError := client.httpClient.Do(acceptLoginRequest)
	if doError != nil {
		return nil, doError
	}

	jsonBody, readAllErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if readAllErr != nil {
		return nil, readAllErr
	}

	var desirializedBody map[string]string
	desirializationErr := json.Unmarshal(jsonBody, &desirializedBody)
	if desirializationErr != nil {
		return nil, desirializationErr
	}

	return desirializedBody, nil
}
