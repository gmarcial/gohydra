package tests

import (
	"github.com/gmarcial/gohydra"
	"testing"
)

func TestMustCreateAValidConsentData(t *testing.T) {
	//Arrange
	consentChallenge := "oidas90i9daklm1-"

	//Action
	consentData, err := gohydra.NewConsentData(consentChallenge)

	//Assert
	if err != nil {
		t.Errorf("Creation is invalid, occurred a error dont expected")
	}

	if consentData == nil {
		t.Errorf("Creation is invalid, consentData is pointer nil.")
	}
}

func TestMustCreateAInvalidConsentDataWhenTheConsentChallengeIsEmpty(t *testing.T) {
	//Arrange
	consentChallenge := ""

	//Action
	consentData, err := gohydra.NewConsentData(consentChallenge)

	//Assert
	if err == nil {
		t.Errorf("Creation is valid, dont occurred a error.")
	}

	if consentData != nil {
		t.Errorf("Creation is valid, consentData dont is nil.")
	}
}

func createAcceptConsentParametersDefault() *gohydra.AcceptConsentParameters {

	acceptConsentParameters := gohydra.AcceptConsentParameters{
		Remember:    true,
		RememberFor: 3600,
	}

	return &acceptConsentParameters
}

func TestAcceptConsentRequestWithSuccess(t *testing.T) {
	//Arrange
	consentData, _ := gohydra.NewConsentData("oisdoijoasd909123n")
	acceptConsentParameters := createAcceptConsentParametersDefault()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(server.URL, true, server.Client())

	//Action
	responseDeserialized, errAcceptConsentRequest := oryhydraClient.AcceptConsentRequest(consentData, acceptConsentParameters)

	//Assert
	if errAcceptConsentRequest != nil {
		t.Errorf("Occurred a error to accept consent valid.")
	}

	if responseDeserialized == nil {
		t.Errorf("The response is pointer nil.")
	}

	if len(responseDeserialized) == 0 {
		t.Errorf("The response is empy.")
	}

	if len(responseDeserialized["redirect_to"]) == 0 {
		t.Errorf("Dont contain the property to redirect.")
	}
}

func TestAcceptConsentRequestFailedWhenTheServerItsUnavailable(t *testing.T) {
	//Arrange
	consentData, _ := gohydra.NewConsentData("oisdoijoasd909123n")
	acceptConsentParameters := createAcceptConsentParametersDefault()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()

	server.Close()

	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptConsentRequest(consentData, acceptConsentParameters)

	//Log the cause of error
	t.Logf(err.Error())

	//Assert
	if err == nil {
		t.Errorf("Dont occurred a error to accept consent request.")
	}

	if responseDeserialized != nil {
		t.Errorf("The consent was accepted of invalid form.")
	}
}

func TestAcceptConsentRequestFailedWhenWhenUnmarshalBodyFailed(t *testing.T) {
	//Arrange
	consentData, _ := gohydra.NewConsentData("oisdoijoasd909123n")
	acceptConsentParameters := createAcceptConsentParametersDefault()
	payload := []byte(`{"redirect_to":: "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()
	defer server.Close()
	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptConsentRequest(consentData, acceptConsentParameters)

	//Log the cause of error
	t.Logf(err.Error())

	//Assert
	if err == nil {
		t.Errorf("Dont occurred a error to accept login request.")
	}

	if responseDeserialized != nil {
		t.Errorf("The login was accepted of invalid form.")
	}
}
