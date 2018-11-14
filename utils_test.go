package main

import (
	"strings"
	"testing"
)

const errorMessageQuantity = "Error of quantity of items returned. Expectation %s Return %s"
const errorMessageData = "Data error returned incorrectly. Expectation %s Return %s"

func TestIdentifyEnvironmentResourceValidScenarios(t *testing.T) {

	type testIdentifyEnvironmentResource struct {
		resource    string
		expectative string
	}
	conditions := [8]testIdentifyEnvironmentResource{
		{"resource-dev", "dev"},
		{"resource-DEV", "dev"},
		{"resource-hml", "hml"},
		{"resource-HML", "hml"},
		{"resource-prod", "prd"},
		{"resource-PROD", "prd"},
		{"resource-prd", "prd"},
		{"resource-PROD", "prd"},
		// {"resource", "undefined"},
		// {"", ""},
	}

	for _, condition := range conditions {

		dataReturn := IdentifyEnvironmentResource(condition.resource)
		if dataReturn != condition.expectative {
			t.Error(errorMessageData, dataReturn, condition.expectative)
		}
	}
}

func TestIdentifyEnvironmentResourceInvalidScenarios(t *testing.T) {

	type testIdentifyEnvironmentResource struct {
		resource    string
		expectative string
	}
	conditions := [5]testIdentifyEnvironmentResource{
		{"resource", "undefined"},
		{"XXXXXXXXXXXX", "undefined"},
		{"Xxxxx", "undefined"},
		{"", ""},
		{"               ", "undefined"},
	}

	for _, condition := range conditions {

		dataReturn := IdentifyEnvironmentResource(condition.resource)
		if dataReturn != condition.expectative {
			t.Error(errorMessageData, dataReturn, condition.expectative)
		}
	}
}

func TestCreateResourceLabelsValidScenarios(t *testing.T) {

	type testCreateResourceLabels struct {
		ID           string
		name         string
		resourceType string
		environment  string
	}

	conditions := [2]testCreateResourceLabels{
		{"resource/subscriptions/your_subscription/resourceGroups/YOUR_RSG/providers/Microsoft.Cache/Redis/my-redis-cache", "my-redis-cache", "Microsoft.Cache/Redis", "undefined"},
		{"resource/subscriptions/your_subscription/resourceGroups/YOUR_RSG/providers/Microsoft.Web/sites/my-frontend/appServices", "my-frontend", "Microsoft.Web/sites", "dev"},
	}

	for _, condition := range conditions {

		dataReturn := CreateResourceLabels(condition.ID, condition.name, condition.resourceType, condition.environment)

		if dataReturn["resource_group"] == "" || dataReturn["resource_group"] != strings.Split(condition.ID, "/")[4] {
			t.Error(errorMessageData, dataReturn["resource_group"], strings.Split(condition.ID, "/")[4])
		}

		if dataReturn["resource_type"] == "" || dataReturn["resource_type"] != condition.resourceType {
			t.Error(errorMessageData, dataReturn["resource_type"], condition.resourceType)
		}

		if dataReturn["resource_name"] == "" || dataReturn["resource_name"] != condition.name {
			t.Error(errorMessageData, dataReturn["resource_name"], condition.name)
		}

		if dataReturn["resource_environment"] == "" || dataReturn["resource_environment"] != condition.environment {
			t.Error(errorMessageData, dataReturn["resource_environment"], condition.environment)
		}
	}
}
