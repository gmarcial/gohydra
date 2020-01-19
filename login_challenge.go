package gohydra

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// LoginData group all data necessaries in the login challenge and process
// LoginChallenge is the identification of the a login challenge in open.
// CookieCRSF cookie to validate if the request it was not forgery.
type LoginData struct {
	LoginChallenge string
}

// NewLoginData construct a new loginData
func NewLoginData(loginChallenge string) (*LoginData, error) {

	if len(loginChallenge) == 0 {
		return nil, errors.New("the login challenge is empty")
	}

	return &LoginData{LoginChallenge: loginChallenge}, nil
}

// GetLoginRequest say to a oryhydra service that a client
// wish accomplish a login, initialize the login process/challenge
// for more: https://www.ory.sh/docs/hydra/sdk/api#get-an-login-request
func (client *Client) GetLoginRequest(loginData *LoginData) (map[string]interface{}, error) {

	getLoginRequestUri, parseErr := client.challengeUriParse(GetLoginRequestUrn, loginData.LoginChallenge)
	if parseErr != nil {
		return nil, parseErr
	}

	response, getErr := client.httpClient.Get(getLoginRequestUri.String())
	if getErr != nil {
		return nil, getErr
	}

	loginRequestSerialized, readAllBytesErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if readAllBytesErr != nil {
		return nil, readAllBytesErr
	}

	var loginRequestDeserialized map[string]interface{}
	jsonUnmarshalErr := json.Unmarshal(loginRequestSerialized, &loginRequestDeserialized)
	if jsonUnmarshalErr != nil {
		return nil, jsonUnmarshalErr
	}

	return loginRequestDeserialized, nil
}

// AcceptLoginParameters required parameters and that
// represent the accepted of a login challenge.
// Docs of the all properties and official: https://www.ory.sh/docs/hydra/sdk/api#accept-an-login-request
type AcceptLoginParameters struct {
	Remember   bool   `json:"remember"`
	RememberTo int    `json:"remember_for"`
	Subject    string `json:"subject"`
}

func NewAcceptLoginParameters(remember bool, rememberTo int, subject string) (*AcceptLoginParameters, error) {

	if len(subject) == 0 {
		return nil, errors.New("the subject it is empty")
	}

	return &AcceptLoginParameters{Remember: remember, RememberTo: rememberTo, Subject: subject}, nil
}

// AcceptLoginRequest accept a login challenge of the oryhydra.
// For more information in offcial oryhydra documentation:
// https://www.ory.sh/docs/hydra/sdk/api#accept-an-login-request
func (client *Client) AcceptLoginRequest(loginData *LoginData, acceptLoginParameters *AcceptLoginParameters) (map[string]string, error) {

	if acceptLoginParameters == nil {
		return nil, errors.New("the acceptLoginParameters is empty")
	}

	acceptLoginUri, parseErr := client.challengeUriParse(AcceptLoginRequestUrn, loginData.LoginChallenge)
	if parseErr != nil {
		return nil, parseErr
	}

	return acceptChallengeRequest(client, acceptLoginParameters, acceptLoginUri)
}
