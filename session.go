package betfair

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var endpoints = map[string]map[string]string{
	"UK": map[string]string{
		"certLogin": "https://identitysso-api.betfair.com/api/certlogin",
		"restLogin": "https://identitysso.betfair.com/api/login",
		"betting":   "https://api.betfair.com/exchange/betting/rest/v1.0/",
		"account":   "https://api.betfair.com/exchange/account/rest/v1.0/",
	},
	"AU": map[string]string{
		"certLogin": "https://identitysso-api.betfair.com/api/certlogin",
		"restLogin": "https://identitysso.betfair.com/api/login",
		"betting":   "https://api-au.betfair.com/exchange/betting/rest/v1.0/",
		"account":   "https://api-au.betfair.com/exchange/account/rest/v1.0/",
	},
}

// NewCredentials func ret val
type CredentialInterface interface{}

// inherited user credentials type
type userCredentials struct {
	Username string
	Password string
	Exchange string
}

// Used for interactive login and request payload and headers
type InteractiveCredentials struct {
	*userCredentials
	ApplicationKey string
}

// Used for NonInteractive (cert-login)
type NonInteractiveCredentials struct {
	*userCredentials
	CertPath struct {
		CrtFile string
		KeyFile string
	}
}

type Session struct {
	token              string
	credentials        CredentialInterface
	requestCredentials *InteractiveCredentials
	httpClient         *http.Client
	logger             *log.Logger
	developerApps      *[]developerApp
}

// returns CredentialInterface
// NewCredentials(username, password, exchange, appkey string) InteractiveCredentials
// NewCredentials(username, password, exchange, crtfile, keyfile) NonInteractiveCredentials
func NewCredentials(params ...string) (CredentialInterface, error) {
	lp := len(params)
	if lp != 4 && lp != 5 {
		return nil, errors.New("invalid credential params")
	}

	var c CredentialInterface

	if lp == 5 {
		// assign NonInteractiveCredential pointer
		c = &NonInteractiveCredentials{
			userCredentials: &userCredentials{
				Username: params[0],
				Password: params[1],
				Exchange: params[2],
			},
			CertPath: struct {
				CrtFile string
				KeyFile string
			}{
				CrtFile: params[3],
				KeyFile: params[4],
			},
		}
		return c, nil
	}

	// if params len is 3 then
	// assign InteractiveCredentials pointer
	c = &InteractiveCredentials{
		userCredentials: &userCredentials{
			Username: params[0],
			Password: params[1],
			Exchange: params[2],
		},
		ApplicationKey: params[3],
	}
	return c, nil
}

// returns http.Client according to credentials
func getHttpClient(credentials CredentialInterface) (*http.Client, error) {
	client := &http.Client{}

	if c, ok := credentials.(*NonInteractiveCredentials); ok {
		// check crt and key file exists
		if _, err := os.Stat(c.CertPath.CrtFile); os.IsNotExist(err) {
			return nil, err
		}

		if _, err := os.Stat(c.CertPath.KeyFile); os.IsNotExist(err) {
			return nil, err
		}

		// load key
		cert, err := tls.LoadX509KeyPair(c.CertPath.CrtFile, c.CertPath.KeyFile)
		if err != nil {
			return nil, err
		}

		// set client transport
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		}
	}

	return client, nil
}

// returns Session struct
func NewSession(credentials CredentialInterface, out ...io.Writer) (*Session,
	error) {
	session := &Session{
		credentials: credentials,
	}

	// set logger if provided
	var logof io.Writer
	if len(out) > 0 {
		logof = out[0]
	} else {
		logof = os.Stderr
	}
	session.logger = log.New(logof, PKG_NAME+" ", log.LstdFlags|log.Lshortfile)

	// set client
	client, err := getHttpClient(session.credentials)
	if err != nil {
		return nil, err
	}
	session.httpClient = client
	// restruct credentials
	setRequestCredentials(session)

	// LOGIN OPS
	payload := strings.NewReader(
		fmt.Sprintf("username=%s&password=%s",
			session.requestCredentials.Username,
			session.requestCredentials.Password))

	var endpoint string = "restLogin"
	if _, ok := session.credentials.(*NonInteractiveCredentials); ok {
		endpoint = "certLogin"
	}

	resp, err := doRequest(session, endpoint, "", payload)
	if err != nil {
		return nil, err
	}

	switch session.credentials.(type) {
	case *InteractiveCredentials:
		var result struct {
			Token   string
			Product string
			Status  string
			Error   string
		}

		if err := json.Unmarshal(resp, &result); err != nil {
			return nil, err
		}

		if result.Status != "SUCCESS" {
			return nil, errors.New(result.Error)
		}
		session.token = result.Token
	case *NonInteractiveCredentials:
		var result struct {
			Token  string `json:"sessionToken"`
			Status string `json:"loginStatus"`
		}

		if err := json.Unmarshal(resp, &result); err != nil {
			return nil, err
		}

		if result.Status != "SUCCESS" {
			return nil, errors.New(result.Status)
		}
		session.token = result.Token
	}

	return session, nil
}

// restruct credentials for requests
func setRequestCredentials(s *Session) {
	if c, ok := s.credentials.(*InteractiveCredentials); ok {
		s.requestCredentials = c
	} else {
		c := s.credentials.(*NonInteractiveCredentials)
		s.requestCredentials = &InteractiveCredentials{
			userCredentials: c.userCredentials,
			ApplicationKey:  "",
		}
	}
}

// prepares betfair endpoints with exchange and method values
func prepareEndpoint(endpoint, method, exchange string) (string, error) {
	var url string
	if _, exists := endpoints[exchange][endpoint]; !exists {
		return "", errors.New(
			fmt.Sprintf("invalid endpoint params: %s %s %s", endpoint, method,
				exchange))
	} else {
		url = endpoints[exchange][endpoint]
	}

	if endpoint != "certLogin" && endpoint != "restLogin" {
		url += method + "/"
	}
	return url, nil
}

// performs request jobs
func doRequest(s *Session, endpoint, method string, body *strings.Reader) (
	[]byte, error) {

	// get completed url
	url, err := prepareEndpoint(endpoint, method, s.requestCredentials.Exchange)
	if err != nil {
		return nil, err
	}

	// prepare request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// setting X-Application header
	var xapph string = s.requestCredentials.ApplicationKey
	if endpoint == "certLogin" {
		xapph = PKG_NAME
	}
	
	req.Header.Set("X-Application", xapph)
	if method == "getDeveloperAppKeys" {
		req.Header.Del("X-Application")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if endpoint != "certLogin" && endpoint != "restLogin" {
		req.Header.Set("X-Authentication", s.token)
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	s.logger.Println(res.Status, url, req.Header, "body:", body, "\n")
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	// defer closing body reader
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
