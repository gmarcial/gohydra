package gohydra

import "errors"

// ConsentData represent the data required to consent challenge ory hydra.
// ConsentChallenge is the identification of the a consent challenge in open.
// For more informations: https://www.ory.sh/docs/next/hydra/sdk/api#get-consent-request-information
type consentData struct {
	ConsentChallenge string
}

// NewConsentData consctruct a new consentData.
func NewConsentData(consentChallenge string) (*consentData, error) {

	if len(consentChallenge) == 0 {
		return nil, errors.New("the consent challenge is empty")
	}

	return &consentData{ConsentChallenge: consentChallenge}, nil
}

// AcceptConsentParameters required parameters and that
// represent the accepted of a consent challenge.
// Docs of the all properties and official: https://www.ory.sh/docs/hydra/sdk/api#schemaacceptconsentrequest
type AcceptConsentParameters struct {
	GrantAccessTokenAudience []string       `json:"grant_access_token_audience"`
	GrantScope               []string       `json:"grant_scope"`
	Remember                 bool           `json:"remember"`
	RememberFor              int            `json:"remember_for"`
	Session                  ConsentSession `json:"session"`
}

// ConsentSession represent the session granted in AccessToken or/and IdToken
// Docs of the all properties and official: https://www.ory.sh/docs/hydra/sdk/api#schemaconsentrequestsession
type ConsentSession struct {
	AccessToken map[string]interface{} `json:"access_token"`
	IdToken     map[string]interface{} `json:"id_token"`
}

// AcceptConsentRequest accept a consent challenge of the oryhydra.
// For more information in offcial oryhydra documentation:
// https://www.ory.sh/docs/hydra/sdk/api#accept-an-consent-request
func (client *Client) AcceptConsentRequest(consentChallenge *consentData, acceptConsentParameters *AcceptConsentParameters) (map[string]string, error) {

	if acceptConsentParameters == nil {
		return nil, errors.New("the acceptConsentParameters is empty")
	}

	acceptConsentUri, urlParseErr := client.challengeUriParse(AcceptConsentRequestUrn, consentChallenge.ConsentChallenge)
	if urlParseErr != nil {
		return nil, urlParseErr
	}

	return acceptChallengeRequest(client, acceptConsentParameters, acceptConsentUri)
}
