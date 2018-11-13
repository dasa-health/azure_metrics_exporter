package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// GetResources get all resoures from azure
func (ac *Client) GetResources(tagValue string) (ResourceResponse, error) {

	if tagValue == "" {
		return ResourceResponse{}, fmt.Errorf("Tag value is empty")
	}

	err := ac.validateAccesssToken()

	if err != nil {
		return ResourceResponse{}, fmt.Errorf("Error refreshing access token: %v", err)
	}
	apiVersion := "2018-05-01"
	subscriptionID := os.Getenv("subscription_id")
	metricValueEndpoint := fmt.Sprintf("%ssubscriptions/%s/resources", ac.resource, subscriptionID)

	log.Print(metricValueEndpoint)
	req, err := http.NewRequest("GET", metricValueEndpoint, nil)
	if err != nil {
		return ResourceResponse{}, fmt.Errorf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+ac.accessToken)

	values := url.Values{}
	resourceQuery := fmt.Sprintf("tagName eq '%s' and tagValue eq '%s'", os.Getenv("resource_query_tag_name"), tagValue)
	if resourceQuery != "" {
		values.Add("$filter", resourceQuery)
	}
	values.Add("api-version", apiVersion)

	req.URL.RawQuery = values.Encode()

	resp, err := ac.client.Do(req)
	if err != nil {
		return ResourceResponse{}, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ResourceResponse{}, fmt.Errorf("Unable to query metrics API with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ResourceResponse{}, fmt.Errorf("Error reading body of response: %v", err)
	}

	var data ResourceResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ResourceResponse{}, fmt.Errorf("Error unmarshalling response body: %v", err)
	}

	return data, nil
}
