package main

import (
	"fmt"
	"regexp"
	"strings"
)

// CreateResourceLabels - Returns resource labels for a give resource ID.
func CreateResourceLabels(resourceID, environment string) map[string]string {
	labels := make(map[string]string)
	labels["resource_group"] = strings.Split(resourceID, "/")[4]
	labels["resource_type"] = strings.Split(resourceID, "/")[6]
	labels["resource_name"] = strings.Split(resourceID, "/")[8]
	labels["resource_environment"] = environment
	return labels
}

// IdentifyEnvironmentResource identifies the environment of a resource
func IdentifyEnvironmentResource(resource string) string {

	if resource == "" {
		return ""
	}

	resource = strings.ToLower(resource)

	if strings.Contains(resource, "prod") || strings.Contains(resource, "prd") {
		return "prd"
	}

	if strings.Contains(resource, "hml") || strings.Contains(resource, "homol") {
		return "hml"
	}

	if strings.Contains(resource, "dev") {
		return "dev"
	}

	return "undefined"
}

// SanitizeMetricName ensure Azure metric names conform to Prometheus metric name conventions
func SanitizeMetricName(name, unit string) (string, error) {

	if name == "" || unit == "" {
		return "", fmt.Errorf("Metric name or metric unit not found")
	}

	invalidMetricChars := regexp.MustCompile("[^a-zA-Z0-9_:]")

	metricName := strings.Replace(name, " ", "_", -1)
	metricName = strings.ToLower(metricName + "_" + unit)
	metricName = strings.Replace(metricName, "/", "_per_", -1)
	metricName = invalidMetricChars.ReplaceAllString(metricName, "_")
	return metricName, nil
}
