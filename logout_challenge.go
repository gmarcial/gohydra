package gohydra

import "errors"

type LogoutData struct {
	LogoutChallenge string
}

// NewLoginData construct a new loginData
func NewLogoutData(logoutChallenge string) (*LogoutData, error) {

	if len(logoutChallenge) == 0 {
		return nil, errors.New("the logout challenge is empty")
	}

	return &LogoutData{LogoutChallenge: logoutChallenge}, nil
}

func (client *Client) AcceptLogoutRequest(logoutData *LogoutData) (map[string]string, error) {

	if logoutData == nil {
		return nil, errors.New("the logout data is empty")
	}

	acceptLogoutUri, parseErr := client.challengeUriParse(AcceptLogoutRequestUrn, logoutData.LogoutChallenge)
	if parseErr != nil {
		return nil, parseErr
	}

	return acceptChallengeRequest(client, nil, acceptLogoutUri)
}
