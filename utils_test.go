package main

import "testing"

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
