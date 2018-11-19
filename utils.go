package main

import (
	"strings"
)

// CreateResourceLabels - Returns resource labels for a give resource ID.
func CreateResourceLabels(resourceID, resourceName, resourceType, environment string) map[string]string {
	labels := make(map[string]string)
	labels["resource_group"] = strings.Split(resourceID, "/")[4]
	labels["resource_type"] = resourceType
	labels["resource_name"] = resourceName
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
