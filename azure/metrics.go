package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dasa-health/azure_metrics_exporter/logger"
)

// GetMetricTypes Loop through all specified resource targets and get their respective metric definitions.
func (ac *Client) GetMetricTypes(resourceName, resourceType string) (MetricDefinitionResponse, error) {
	err := ac.validateAccesssToken()

	if err != nil {
		logger.Error("[GetMetricTypes] - Error in validation access token", err)
		return MetricDefinitionResponse{}, fmt.Errorf("Error refreshing access token: %v", err)
	}

	apiVersion := "2018-01-01"

	metricsTarget := fmt.Sprintf("%s%s/providers/microsoft.insights/metricDefinitions", ac.resource, resourceName)
	req, err := http.NewRequest("GET", metricsTarget, nil)
	if err != nil {
		return MetricDefinitionResponse{}, fmt.Errorf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+ac.accessToken)
	values := url.Values{}
	values.Add("api-version", apiVersion)

	req.URL.RawQuery = values.Encode()

	resp, err := ac.client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetricTypes] - Error in GET %s", req.URL), err)
		return MetricDefinitionResponse{}, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("[GetMetricTypes] - Error in GET %s", req.URL), resp.StatusCode)
		return MetricDefinitionResponse{}, fmt.Errorf("Unable to query metrics API with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetricTypes] - Error in GET %s", req.URL), err)
		return MetricDefinitionResponse{}, fmt.Errorf("Error reading body of response: %v", err)
	}

	var data MetricDefinitionResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetricTypes] - Error in GET %s", req.URL), err)
		return MetricDefinitionResponse{}, fmt.Errorf("Error unmarshalling response body: %v", err)
	}

	return data, nil
}

// GetMetric retrieves resource metrics in azure
func (ac *Client) GetMetric(resource, metricNames, aggregation string) (MetricValueResponse, error) {

	err := ac.validateAccesssToken()

	if err != nil {
		logger.Error("[GetMetric] - Error in validation access token", err)
		return MetricValueResponse{}, fmt.Errorf("Error refreshing access token: %v", err)
	}

	apiVersion := "2018-01-01"

	endTime, startTime := getTimes()

	metricValueEndpoint := fmt.Sprintf("%s%s/providers/microsoft.insights/metrics", ac.resource, resource)

	req, err := http.NewRequest("GET", metricValueEndpoint, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetric] - Error in GET %s", req.URL), req)
		return MetricValueResponse{}, fmt.Errorf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+ac.accessToken)

	values := url.Values{}
	if metricNames != "" {
		values.Add("metricnames", metricNames)
	}

	values.Add("aggregation", aggregation)
	values.Add("timespan", fmt.Sprintf("%s/%s", startTime, endTime))
	values.Add("api-version", apiVersion)

	req.URL.RawQuery = values.Encode()

	resp, err := ac.client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetric] - Error in GET %s", req.URL), err)
		return MetricValueResponse{}, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("[GetMetric] - Error in GET %s", req.URL), resp.StatusCode)
		return MetricValueResponse{}, fmt.Errorf("Unable to query metrics API with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetric] - Error in GET %s", req.URL), err)
		return MetricValueResponse{}, fmt.Errorf("Error reading body of response: %v", err)
	}

	var data MetricValueResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(fmt.Sprintf("[GetMetric] - Error in GET %s", req.URL), err)
		return MetricValueResponse{}, fmt.Errorf("Error unmarshalling response body: %v", err)
	}

	return data, nil
}

// GetTimes - Returns the endTime and startTime used for querying Azure Metrics API
func getTimes() (string, string) {
	// Make sure we are using UTC
	now := time.Now().UTC()

	// Use query delay of 3 minutes when querying for latest metric data
	endTime := now.Add(time.Minute * time.Duration(-3)).Format(time.RFC3339)
	startTime := now.Add(time.Minute * time.Duration(-4)).Format(time.RFC3339)
	return endTime, startTime
}

// TreatTypeMetric performs metric type api return processing for use in metric api
func TreatTypeMetric(typeMetrics MetricDefinitionResponse) []string {
	if len(typeMetrics.MetricDefinitionResponses) <= 0 {
		return []string{}
	}
	definitions := []string{}
	totalMetrics := len(typeMetrics.MetricDefinitionResponses)
	internalIndex := 0
	metric := ""

	for index, typeMetric := range typeMetrics.MetricDefinitionResponses {

		if typeMetric.Name.Value == "" || strings.Trim(typeMetric.Name.Value, " ") == "" {
			internalIndex++
		} else if internalIndex <= 19 {
			internalIndex++
			metric += typeMetric.Name.Value + ","
		} else if internalIndex > 19 {
			internalIndex = 0
			definitions = append(definitions, metric[0:(len(metric)-1)])
			metric = typeMetric.Name.Value + ","
		}

		if index == (totalMetrics-1) && metric != "" {
			definitions = append(definitions, metric[0:(len(metric)-1)])
		}
	}

	return definitions
}

const typesAllowed = "Microsoft.AnalysisServices/servers,CloudSimple.PrivateCloudIaaS/virtualMachines,Microsoft.Web/serverFarms,Microsoft.Web/sites,Microsoft.Web/sites/slots,Microsoft.Web/hostingEnvironments/multiRolePools,Microsoft.Web/hostingEnvironments/workerPools,test.shoebox/testresources,test.shoebox/testresources2,Microsoft.ServiceBus/namespaces,Microsoft.Network/virtualNetworks,Microsoft.Network/publicIPAddresses,Microsoft.Network/networkInterfaces,Microsoft.Network/loadBalancers,Microsoft.Network/networkWatchers/connectionMonitors,Microsoft.Network/virtualNetworkGateways,Microsoft.Network/connections,Microsoft.Network/applicationGateways,Microsoft.Network/dnszones,Microsoft.Network/trafficmanagerprofiles,Microsoft.Network/expressRouteCircuits,Microsoft.Network/expressRoutePorts,Microsoft.Network/azureFirewalls,Microsoft.Network/frontdoors,Microsoft.DBforMySQL/servers,Microsoft.Sql/servers,Microsoft.Sql/servers/databases,Microsoft.Sql/servers/elasticpools,Microsoft.Sql/managedInstances,microsoft.insights/components,microsoft.insights/autoscalesettings,Microsoft.KeyVault/vaults,Microsoft.Cache/Redis,Microsoft.ContainerRegistry/registries,Microsoft.LocationBasedServices/accounts,Microsoft.DocumentDB/databaseAccounts,Microsoft.ContainerInstance/containerGroups,Microsoft.Devices/IotHubs,Microsoft.Devices/ElasticPools,Microsoft.Devices/ElasticPools/IotHubTenants,Microsoft.Devices/ProvisioningServices,Microsoft.Compute/virtualMachines,Microsoft.Compute/virtualMachineScaleSets,Microsoft.Compute/virtualMachineScaleSets/virtualMachines,Microsoft.ClassicCompute/domainNames/slots/roles,Microsoft.ClassicCompute/virtualMachines,Microsoft.SignalRService/SignalR,Microsoft.DataBoxEdge/DataBoxEdgeDevices,Microsoft.Search/searchServices,Microsoft.Logic/workflows,Microsoft.Logic/integrationServiceEnvironments,Microsoft.HDInsight/clusters,Microsoft.Relay/namespaces,Microsoft.EventHub/namespaces,Microsoft.EventHub/clusters,Microsoft.Kusto/clusters,Microsoft.OperationalInsights/workspaces,Microsoft.Maps/accounts,Microsoft.DBforMariaDB/servers,Microsoft.TimeSeriesInsights/environments,Microsoft.TimeSeriesInsights/environments/eventsources,Microsoft.DBforPostgreSQL/servers,Microsoft.StreamAnalytics/streamingjobs,Microsoft.NotificationHubs/namespaces/notificationHubs,Microsoft.ApiManagement/service,Microsoft.Storage/storageAccounts,Microsoft.Storage/storageAccounts/blobServices,Microsoft.Storage/storageAccounts/tableServices,Microsoft.Storage/storageAccounts/queueServices,Microsoft.Storage/storageAccounts/fileServices,Microsoft.DataLakeAnalytics/accounts,Microsoft.PowerBIDedicated/capacities,Microsoft.IoTSpaces/Graph,Microsoft.Automation/automationAccounts,Microsoft.DataLakeStore/accounts,Microsoft.DataFactory/dataFactories,Microsoft.DataFactory/factories,Microsoft.NetApp/netAppAccounts/capacityPools,Microsoft.NetApp/netAppAccounts/capacityPools/volumes,Microsoft.StorageSync/storageSyncServices,Microsoft.StorageSync/storageSyncServices/syncGroups,Microsoft.StorageSync/storageSyncServices/syncGroups/serverEndpoints,Microsoft.StorageSync/storageSyncServices/registeredServers,Microsoft.ContainerService/managedClusters,Microsoft.CustomerInsights/hubs,Microsoft.Batch/batchAccounts,Microsoft.EventGrid/eventSubscriptions,Microsoft.EventGrid/topics,Microsoft.EventGrid/domains,Microsoft.EventGrid/extensionTopics,Microsoft.CognitiveServices/accounts"

// ValidateTypeMetric valid if the resource type has some metric definition in the azure api
func ValidateTypeMetric(metricType string) bool {

	if metricType == "" || strings.Trim(metricType, " ") == "" {
		return false
	}

	if !strings.Contains(typesAllowed, metricType) {
		return false
	}

	return true
}

// SanitizeMetric is the method responsible for performing all treatments in the metrics recovered in azure
func (value *MetricValueResponseValue) SanitizeMetric(resourceType string) error {

	metricValue := value.Timeseries[0].Data[len(value.Timeseries[0].Data)-1]

	value.Unit = strings.ToLower(value.Unit)
	if value.Unit != "milliseconds" {
		metricName, err := sanitizeMetricName(value.Name.Value, value.Unit, resourceType)

		if err != nil {
			return err
		}

		value.Name.Value = metricName
		return err
	}

	value.Unit = "seconds"
	metricValue.Total = convertMillisToSeconds(metricValue.Total)
	metricValue.Average = convertMillisToSeconds(metricValue.Average)
	metricValue.Maximum = convertMillisToSeconds(metricValue.Maximum)
	metricValue.Minimum = convertMillisToSeconds(metricValue.Minimum)

	metricName, err := sanitizeMetricName(value.Name.Value, value.Unit, resourceType)

	if err != nil {
		return err
	}

	value.Name.Value = metricName

	return nil
}

// SanitizeMetricName ensure Azure metric names conform to Prometheus metric name conventions
func sanitizeMetricName(name, unit, resourceType string) (string, error) {

	if name == "" || unit == "" {
		return "", fmt.Errorf("Metric name or metric unit not found")
	}

	invalidMetricChars := regexp.MustCompile("[^a-zA-Z0-9_:]")

	if unit == "total" || unit == "count" {
		unit = "amount"
	}

	metricName := "azure_" + strings.Replace(resourceType, ".", "_", -1)
	metricName = strings.Replace(metricName, "/", "_", -1)
	metricName = metricName + "_" + strings.Replace(name, " ", "_", -1)
	metricName = strings.ToLower(metricName + "_" + unit)
	metricName = strings.Replace(metricName, "/", "_per_", -1)
	metricName = invalidMetricChars.ReplaceAllString(metricName, "_")

	return metricName, nil
}

// convertMillisToSeconds convete milisegundos para segundos
func convertMillisToSeconds(millis float64) float64 {
	return millis / 1000
}
