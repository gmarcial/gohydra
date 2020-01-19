package tests

import (
	"github.com/gmarcial/gohydra"
	"testing"
)

func createDefaultLogoutData() *gohydra.LogoutData {
	logoutData, _ := gohydra.NewLogoutData("klnda90091j123ksa")
	return logoutData
}

func TestMustCreateAValidLogoutData(t *testing.T) {
	//Arrange
	logoutChallenge := "0-123pokdskjldasls"

	//Action
	logoutData, err := gohydra.NewLogoutData(logoutChallenge)

	//Assert
	if err != nil {
		t.Errorf("A valid logout data dont was created, error returned.")
	}

	if logoutData == nil {
		t.Errorf("A valid logout data dont was created, pointer nil.")
	}
}

func TestMustCreateInvalidLogoutDataWhenLogoutChallengeBeEmpty(t *testing.T) {
	//Arrange
	logoutChallenge := ""

	//Action
	logoutData, err := gohydra.NewLogoutData(logoutChallenge)

	//Assert
	if err == nil {
		t.Errorf("A valid logout data was created, dont error returned .")
	}

	if logoutData != nil {
		t.Errorf("A invalid logout data was created.")
	}
}

func TestAcceptLogoutRequestWithSuccess(t *testing.T) {
	//Arrange
	logoutData := createDefaultLogoutData()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(server.URL, true, server.Client())

	//Action
	responseDeserialized, errAcceptLoginRequest := oryhydraClient.AcceptLogoutRequest(logoutData)

	//Assert
	if errAcceptLoginRequest != nil {
		t.Errorf("Occurred a error to accept logout valid.")
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

func TestAcceptLogoutRequestFailedWhenTheServerItsUnavailable(t *testing.T) {
	//Arrange
	logoutData := createDefaultLogoutData()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()

	server.Close()

	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptLogoutRequest(logoutData)

	//Log the cause of error
	t.Logf(err.Error())

	//Assert
	if err == nil {
		t.Errorf("Dont occurred a error to accept logout request.")
	}

	if responseDeserialized != nil {
		t.Errorf("The logout was accepted of invalid form.")
	}
}

func TestAcceptLogoutRequestFailedWhenWhenUnmarshalBodyFailed(t *testing.T) {
	//Arrange
	logoutData := createDefaultLogoutData()
	payload := []byte(`{"redirect_to":: "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptLogoutRequest(logoutData)

	//Log the cause of error
	t.Logf(err.Error())

	//Assert
	if err == nil {
		t.Errorf("Dont occurred a error to accept logout request.")
	}

	if responseDeserialized != nil {
		t.Errorf("The logout was accepted of invalid form.")
	}
}