package betfair

import (
	"testing"
)

func Test_NewCredentials(t *testing.T) {
	c, err := NewCredentials("userName", "passWord", "UK", "appKey")
	if err != nil {
		t.Fatal(err)
	}
	u := c.(*InteractiveCredentials)

	if u.Username != "userName" || u.Password != "passWord" ||
		u.ApplicationKey != "appKey" || u.Exchange != "UK" {
		t.Error("interactive credential errors")
	}

	c, err = NewCredentials("userName", "passWord", "UK", "crt", "key")
	if err != nil {
		t.Fatal(err)
	}

	v := c.(*NonInteractiveCredentials)
	if v.Username != "userName" || v.Password != "passWord" ||
		v.Exchange != "UK" || v.CertPath.CrtFile != "crt" ||
		v.CertPath.KeyFile != "key" {
		t.Error("non-interactive credential errors")
	}

	c, err = NewCredentials()
	if err == nil {
		t.Fatal("not returned error")
	}
}

func Test_getHttpClient(t *testing.T) {
	c, _ := NewCredentials("username", "pass", "UK", "appKey")
	client, err := getHttpClient(c)
	if err != nil {
		t.Fatal(err)
	}

	if client.Transport != nil {
		t.Error("client error")
	}
	
	v, _ := NewCredentials("username", "pass", "UK",
		"/home/baris/whore/betfair-certs/client-2048.crt",
		"/home/baris/whore/betfair-certs/client-2048.key")
	client, err = getHttpClient(v)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_prepareEndpoint(t *testing.T) {
	url, err := prepareEndpoint("betting", "listEvents", "UK")
	if err != nil {
		t.Fatal(err)
	}

	if url != "https://api.betfair.com/exchange/betting/rest/v1.0/listEvents" {
		t.Error("preparing url wrong")
	}
}

/*
func Test_NewSession(t *testing.T) {
	
	}
func Test_doRequest(t *testing.T)  {}
*/