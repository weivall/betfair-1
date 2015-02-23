package betfair

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type developerApp struct {
	AppName     string
	AppId       uint32
	AppVersions []struct {
		Owner                string
		VersionId            uint32
		Version              string
		ApplicationKey       string
		DelayData            bool
		SubscriptionRequired bool
		OwnerManaged         bool
		Active               bool
	}
}

// Sets request application key by provided name and delay info
func (s *Session) SetUsedApplication(name string, delay bool) error {
	// first set session developer apps
	if s.developerApps == nil {
		var apps []developerApp
		resp, err := s.GetDeveloperAppKeys()
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(resp), &apps); err != nil {
			return err
		}
		s.developerApps = &apps
	}

	// find application by provided app name and delay info
	for _, app := range *s.developerApps {
		if app.AppName != name {
			continue
		}

		for _, versions := range app.AppVersions {
			if versions.DelayData != delay || !versions.Active {
				continue
			}

			s.requestCredentials.ApplicationKey = versions.ApplicationKey
		}
	}

	if s.requestCredentials.ApplicationKey == "" {
		return errors.New(fmt.Sprintf("%s application not found", name))
	}

	return nil
}

func (s *Session) GetDeveloperAppKeys() (string, error) {
	payload := strings.NewReader("")
	resp, err := doRequest(s, "account", "getDeveloperAppKeys", payload)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func (s *Session) GetAccountFunds() (string, error) {
	payload := strings.NewReader("")
	resp, err := doRequest(s, "account", "getAccountFunds", payload)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func (s *Session) GetAccountDetails() (string, error) {
	payload := strings.NewReader("")
	resp, err := doRequest(s, "account", "getAccountDetails", payload)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
