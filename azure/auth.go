package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dasa-health/azure_metrics_exporter/logger"
)

// GetAccessToken autentica o exporter na azure
func GetAccessToken() (Client, error) {

	ac := newAzureClient()
	clientID := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	tenantID := os.Getenv("tenant_id")

	target := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantID)
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"resource":      {"https://management.azure.com/"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}
	resp, err := ac.client.PostForm(target, form)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetAccessToken] - Error in GET %s", target), err)
		return Client{}, fmt.Errorf("Error authenticating against Azure API: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("[GetAccessToken] - Error in GET %s", target), resp.StatusCode)
		return Client{}, fmt.Errorf("Did not get status code 200, got: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetAccessToken] - Error in GET %s", target), err)
		return Client{}, fmt.Errorf("Error reading body of response: %v", err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetAccessToken] - Error in GET %s", target), err)
		return Client{}, fmt.Errorf("Error unmarshalling response body: %v", err)
	}
	ac.accessToken = data["access_token"].(string)
	ac.resource = data["resource"].(string)
	expiresOn, err := strconv.ParseInt(data["expires_on"].(string), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetAccessToken] - Error in GET %s", target), err)
		return Client{}, fmt.Errorf("Error ParseInt of expires_on failed: %v", err)
	}
	ac.accessTokenExpiresOn = time.Unix(expiresOn, 0).UTC()

	return ac, nil
}

func (ac *Client) validateAccesssToken() error {

	now := time.Now().UTC()
	refreshAt := ac.accessTokenExpiresOn.Add(-10 * time.Minute)
	if now.After(refreshAt) {
		newAc, err := GetAccessToken()
		if err != nil {
			//log.Fatal("Error refreshing access token: %v", err)
			return err
		}

		ac.accessToken = newAc.accessToken
		ac.resource = newAc.resource
		ac.accessTokenExpiresOn = newAc.accessTokenExpiresOn

		return nil
	}

	return nil
}
