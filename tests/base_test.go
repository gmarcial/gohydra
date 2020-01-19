package tests

import (
	"crypto/tls"
	"github.com/gmarcial/gohydra"
	"net/http"
	"testing"
	"time"
)

const (
	baseUrlNull  = ""
	badBaseUrl   = "htp::///localhost:4445"
	baseUrlHttp  = "http://localhost:4445"
	baseUrlHttps = "https://localhost:4445"
)

func createDefaultClientHttp() *http.Client {
	temporaryInsecureTls := tls.Config{
		InsecureSkipVerify: true}

	roundTripper := http.Transport{
		TLSClientConfig:     &temporaryInsecureTls,
		TLSHandshakeTimeout: time.Minute * 3,
		DisableKeepAlives:   false,
		MaxIdleConns:        5}

	return &http.Client{
		Timeout:   time.Minute * 5,
		Transport: &roundTripper}
}

func TestMustCreateAValidClient(t *testing.T) {
	//Arrange
	baseUrl := baseUrlHttp
	httpsRequired := false
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err != nil {
		t.Errorf(err.Error())
	}

	if client == nil {
		t.Errorf("The client that was create is nil.")
	}
}

func TestMustCreateAValidClientWithTls(t *testing.T) {
	//Arrange
	baseUrl := baseUrlHttps
	httpsRequired := true
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err != nil {
		t.Errorf(err.Error())
	}

	if client == nil {
		t.Errorf("The client that was create is nil.")
	}
}

func TestMustCreateAInvalidClientWhenBaseUrlIsEmpty(t *testing.T) {
	//Arrange
	baseUrl := baseUrlNull
	httpsRequired := true
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err == nil {
		t.Errorf("The client was created with the empty base url.")
	}

	if client != nil {
		t.Errorf("The client that was create is different of nil.")
	}

}

func TestMustCreateInvalidClientWhenHttpsRequiridButBaseUrlDontIsHttps(t *testing.T) {
	//Arrange
	baseUrl := baseUrlHttp
	httpsRequired := true
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err == nil {
		t.Errorf("The client was created with base url http, but is required https.")
	}

	if client != nil {
		t.Errorf("The client that was create is different of nil.")
	}
}

func TestMustCreateInvalidClientWhenHttpsNotRequiridButBaseUrlDontIsHttp(t *testing.T) {
	//Arrange
	baseUrl := baseUrlHttps
	httpsRequired := false
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err == nil {
		t.Errorf("The client was created with base url https, but is required http.")
	}

	if client != nil {
		t.Errorf("The client that was create is different of nil.")
	}
}

func TestMustCreateInvalidClientWhenBaseUrlIsAnBadUrl(t *testing.T) {
	//Arrange
	baseUrl := badBaseUrl
	httpsRequired := false
	clientHttp := createDefaultClientHttp()

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, clientHttp)

	//Assert
	if err == nil {
		t.Errorf("The client was created with an bad base url.")
	}

	if client != nil {
		t.Errorf("The client that was create is different of nil.")
	}
}

func TestMustCreateInvalidClientWhenTheHttpClientBeInvalid(t *testing.T) {
	//Arrange
	baseUrl := baseUrlHttps
	httpsRequired := true

	//Action
	client, err := gohydra.NewClient(baseUrl, httpsRequired, nil)

	//Assert
	if err == nil {
		t.Errorf("The client was created with an http client invalid.")
	}

	if client != nil {
		t.Errorf("The client that was create is different of nil.")
	}
}
