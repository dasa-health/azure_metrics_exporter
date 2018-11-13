package azure

import (
	"net/http"
	"time"
)

// Client represents our client to talk to the Azure api
type Client struct {
	client               *http.Client
	accessToken          string
	accessTokenExpiresOn time.Time
	resource             string
}

func newAzureClient() Client {
	return Client{
		client:               &http.Client{},
		accessToken:          "",
		accessTokenExpiresOn: time.Time{},
		resource:             "",
	}
}

// ResourceResponse represents generic resource for Azure
type ResourceResponse struct {
	Value []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"value"`
}

// MetricDefinitionResponse represents metric definition response for a given resource from Azure.
type MetricDefinitionResponse struct {
	MetricDefinitionResponses []metricDefinitionResponse `json:"value"`
}
type metricDefinitionResponse struct {
	Dimensions []struct {
		LocalizedValue string `json:"localizedValue"`
		Value          string `json:"value"`
	} `json:"dimensions"`
	ID                   string `json:"id"`
	IsDimensionRequired  bool   `json:"isDimensionRequired"`
	MetricAvailabilities []struct {
		Retention string `json:"retention"`
		TimeGrain string `json:"timeGrain"`
	}
	Name struct {
		LocalizedValue string `json:"localizedValue"`
		Value          string `json:"value"`
	} `json:"name"`
	PrimaryAggregationType string `json:"primaryAggregationType"`
	ResourceID             string `json:"resourceId"`
	Unit                   string `json:"unit"`
}

type metricDefinitionResponseName struct {
	LocalizedValue string `json:"localizedValue"`
	Value          string `json:"value"`
}

// MetricValueResponse represents a metric value response for a given metric definition.
type MetricValueResponse struct {
	Value []struct {
		Timeseries []struct {
			Data []struct {
				TimeStamp string  `json:"timeStamp"`
				Total     float64 `json:"total"`
				Average   float64 `json:"average"`
				Minimum   float64 `json:"minimum"`
				Maximum   float64 `json:"maximum"`
			} `json:"data"`
		} `json:"timeseries"`
		ID   string `json:"id"`
		Name struct {
			LocalizedValue string `json:"localizedValue"`
			Value          string `json:"value"`
		} `json:"name"`
		Type string `json:"type"`
		Unit string `json:"unit"`
	} `json:"value"`
	APIError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
