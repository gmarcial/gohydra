package tests

import (
	"github.com/gmarcial/gohydra"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createDefaultLoginData() *gohydra.LoginData {
	loginData, _ := gohydra.NewLoginData("klnda90091j123ksa")
	return loginData
}

func TestMustCreateAValidLoginData(t *testing.T) {
	//Arrange
	loginChallenge := "0-123pokdskjldasls"

	//Action
	loginData, err := gohydra.NewLoginData(loginChallenge)

	//Assert
	if err != nil {
		t.Errorf("A valid login data dont was created, error returned.")
	}

	if loginData == nil {
		t.Errorf("A valid login data dont was created, pointer nil.")
	}
}

func TestMustCreateInvalidLoginDataWhenLoginChallengeBeEmpty(t *testing.T) {
	//Arrange
	loginChallenge := ""

	//Action
	loginData, err := gohydra.NewLoginData(loginChallenge)

	//Assert
	if err == nil {
		t.Errorf("A valid login data was created, dont error returned .")
	}

	if loginData != nil {
		t.Errorf("A invalid login data was created.")
	}
}

func createTLSServerDefault(payload []byte) *httptest.Server {
	getLoginRequestMock := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(payload)
		})
	return httptest.NewTLSServer(getLoginRequestMock)
}

func TestShouldGetALoginRequestWithSuccess(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	payload := []byte(`{"skip": false}`)

	server := createTLSServerDefault(payload)
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(server.URL, true, server.Client())

	//Action
	loginRequestDeserialized, err := oryhydraClient.GetLoginRequest(loginData)

	//Assert
	if err != nil {
		t.Errorf("Failed to get login request, error returned")
	}

	if loginRequestDeserialized == nil {
		t.Errorf("Failed to get login request, pointer nil.")
	}

	if len(loginRequestDeserialized) == 0 {
		t.Errorf("The login request is empy.")
	}

	if loginRequestDeserialized["skip"] == nil {
		t.Errorf("Dont contain the property to evaluate skip login.")
	}
}

func closeHttpServerTest(server *httptest.Server) {
	time.Sleep(1 * time.Millisecond)
	server.Close()
}

func TestGetALoginRequestShouldFailedWhenUnmarshalBodyFailed(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	payload := []byte(`{"skip":: false}`)

	server := createTLSServerDefault(payload)
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(server.URL, true, server.Client())

	//Action
	loginRequestDeserialized, err := oryhydraClient.GetLoginRequest(loginData)

	//Log the cause of error
	t.Logf(err.Error())

	//Assert
	if err == nil {
		t.Errorf("Not falied to get login request, dont error returned")
	}

	if loginRequestDeserialized != nil {
		t.Errorf("Not falied to get login request, login Request returned.")
	}
}

func TestMustCreateAValidAcceptLoginParameters(t *testing.T) {
	//Arrange
	remember := true
	rememberTo := 3600
	subject := "123"

	//Action
	acceptLoginParameters, err := gohydra.NewAcceptLoginParameters(remember, rememberTo, subject)

	//Assert
	if err != nil {
		t.Errorf("Occurred a error to create a accept login parameters valid.")
	}

	if acceptLoginParameters == nil {
		t.Errorf("The accept login parameters dont created with success.")
	}
}

func TestDontCreateAAcceptLoginParametersWithEmptySubject(t *testing.T) {
	//Arrange
	remember := true
	rememberTo := 3600
	subject := ""

	//Action
	acceptLoginParameters, err := gohydra.NewAcceptLoginParameters(remember, rememberTo, subject)

	//Assert
	if err == nil {
		t.Errorf("Dont occurred a error to create a accept login parameters.")
	}

	if acceptLoginParameters != nil {
		t.Errorf("A accept login parameters invalid was created.")
	}
}

func createAcceptLoginParametersDefault() *gohydra.AcceptLoginParameters {
	parameters, _ := gohydra.NewAcceptLoginParameters(true, 3600, "unit-test")
	return parameters
}

func TestAcceptLoginRequestWithSuccess(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	acceptLoginParameters := createAcceptLoginParametersDefault()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(server.URL, true, server.Client())

	//Action
	responseDeserialized, errAcceptLoginRequest := oryhydraClient.AcceptLoginRequest(loginData, acceptLoginParameters)

	//Assert
	if errAcceptLoginRequest != nil {
		t.Errorf("Occurred a error to accept login valid.")
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

func TestAcceptLoginRequestFailedWhenTheServerItsUnavailable(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	acceptLoginParameters := createAcceptLoginParametersDefault()
	payload := []byte(`{"redirect_to": "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()

	server.Close()

	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptLoginRequest(loginData, acceptLoginParameters)

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

func TestAcceptLoginRequestFailedWhenABadBodyRead(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	acceptLoginParameters := createAcceptLoginParametersDefault()
	payload := make([]byte, 999)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()
	defer server.Close()

	go closeHttpServerTest(server)
	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptLoginRequest(loginData, acceptLoginParameters)

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

func TestAcceptLoginRequestFailedWhenWhenUnmarshalBodyFailed(t *testing.T) {
	//Arrange
	loginData := createDefaultLoginData()
	acceptLoginParameters := createAcceptLoginParametersDefault()
	payload := []byte(`{"redirect_to":: "string"}`)

	server := createTLSServerDefault(payload)
	baseUrl := server.URL
	httpClient := server.Client()
	defer server.Close()

	oryhydraClient, _ := gohydra.NewClient(baseUrl, true, httpClient)

	//Action
	responseDeserialized, err := oryhydraClient.AcceptLoginRequest(loginData, acceptLoginParameters)

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
