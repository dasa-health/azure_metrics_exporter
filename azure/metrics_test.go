package azure

import "testing"

const errorMessageQuantity = "Error of quantity of items returned. Expectation %s Return %s"
const errorMessageData = "Data error returned incorrectly. Expectation %s Return %s"

func TestTreatTypeMetricUpTo20Metrics(t *testing.T) {

	metricDefinitionResponse := MetricDefinitionResponse{
		MetricDefinitionResponses: []metricDefinitionResponse{
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "CpuTime"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Requests"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "BytesReceived"}},
		},
	}

	response := TreatTypeMetric(metricDefinitionResponse)

	if response == nil || len(response) != 1 {
		t.Error(errorMessageQuantity, 1, len(response))
	}

	expectedDataReturn := "CpuTime,Requests,BytesReceived"
	if expectedDataReturn != response[0] {
		t.Error(errorMessageData, expectedDataReturn, response[0])
	}
}

func TestTreatTypeMetricPlusTo20Metrics(t *testing.T) {

	metricDefinitionResponse := MetricDefinitionResponse{
		MetricDefinitionResponses: []metricDefinitionResponse{
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "CpuTime"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Requests"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "BytesReceived"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "BytesSent"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http101"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http2xx"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http3xx"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http401"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http403"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http404"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http406"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http4xx"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Http5xx"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "MemoryWorkingSet"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "AverageMemoryWorkingSet"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "AverageResponseTime"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "AppConnections"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Handles"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Threads"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "PrivateBytes"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoReadBytesPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoWriteBytesPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoOtherBytesPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoReadOperationsPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoWriteOperationsPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "IoOtherOperationsPerSecond"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "RequestsInApplicationQueue"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "CurrentAssemblies"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "TotalAppDomains"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "TotalAppDomainsUnloaded"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Gen0Collections"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Gen1Collections"}},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "Gen2Collections"}},
		},
	}

	response := TreatTypeMetric(metricDefinitionResponse)

	expectedFirstDataReturn := "CpuTime,Requests,BytesReceived,BytesSent,Http101,Http2xx,Http3xx,Http401,Http403,Http404,Http406,Http4xx,Http5xx,MemoryWorkingSet,AverageMemoryWorkingSet,AverageResponseTime,AppConnections,Handles,Threads,PrivateBytes"
	expectedSecondDataReturn := "IoReadBytesPerSecond,IoWriteBytesPerSecond,IoOtherBytesPerSecond,IoReadOperationsPerSecond,IoWriteOperationsPerSecond,IoOtherOperationsPerSecond,RequestsInApplicationQueue,CurrentAssemblies,TotalAppDomains,TotalAppDomainsUnloaded,Gen0Collections,Gen1Collections,Gen2Collections"

	if response == nil || len(response) != 2 {
		t.Error(errorMessageQuantity, 2, len(response))
	}

	if expectedFirstDataReturn != response[0] {
		t.Error(errorMessageData, expectedFirstDataReturn, response[0])
	}

	if expectedSecondDataReturn != response[1] {
		t.Error(errorMessageData, expectedSecondDataReturn, response[1])
	}
}
func TestTreatTypeMetricEmptyMetrics(t *testing.T) {

	metricDefinitionResponse := MetricDefinitionResponse{
		MetricDefinitionResponses: []metricDefinitionResponse{},
	}

	response := TreatTypeMetric(metricDefinitionResponse)

	if response == nil || len(response) != 0 {
		t.Error(errorMessageQuantity, 0, len(response))
	}
}

func TestTreatTypeMetricNillMetrics(t *testing.T) {

	metricDefinitionResponse := MetricDefinitionResponse{}

	response := TreatTypeMetric(metricDefinitionResponse)

	if response == nil || len(response) != 0 {
		t.Error(errorMessageQuantity, 0, len(response))
	}
}

func TestTreatTypeMetricEmptyValues(t *testing.T) {

	metricDefinitionResponse := MetricDefinitionResponse{
		MetricDefinitionResponses: []metricDefinitionResponse{
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: ""}},
			metricDefinitionResponse{},
			metricDefinitionResponse{Name: metricDefinitionResponseName{Value: "   "}},
		},
	}

	response := TreatTypeMetric(metricDefinitionResponse)

	if response == nil || len(response) != 0 {
		t.Error(errorMessageQuantity, 0, len(response))
	}
}

func TestValidateTypeMetricRealMetric(t *testing.T) {

	type testValidateTypeMetric struct {
		typeMetric  string
		expectative bool
	}
	conditions := [5]testValidateTypeMetric{
		{"Microsoft.AnalysisServices/servers", true},
		{"Microsoft.Compute/availabilitySets", false},
		{"xxxxxxxxxxxxxxxxxx", false},
		{"", false},
		{"         ", false},
	}

	for _, condition := range conditions {

		dataReturn := ValidateTypeMetric(condition.typeMetric)
		if dataReturn != condition.expectative {
			t.Error(errorMessageData, dataReturn, condition.expectative)
		}
	}
}

func TestConvertMillisToSecondsValidScenarios(t *testing.T) {

	type testConvertMillisToSeconds struct {
		millis  float64
		seconds float64
	}
	conditions := [6]testConvertMillisToSeconds{
		{0, 0},
		{1, 0.001},
		{11, 0.011},
		{145, 0.145},
		{99999999, 99999.999},
		{102939455, 102939.455},
	}

	for _, condition := range conditions {

		dataReturn := convertMillisToSeconds(condition.millis)
		if dataReturn != condition.seconds {
			t.Error(errorMessageData, dataReturn, condition.seconds)
		}
	}

}
