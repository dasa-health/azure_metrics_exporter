package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RobustPerception/azure_metrics_exporter/config"
)

// PrintPrettyJSON - Prints structs nicely for debugging.
func PrintPrettyJSON(input map[string]interface{}) {
	out, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		log.Fatalf("Error indenting JSON: %v", err)
	}
	fmt.Println(string(out))
}

// GetTimes - Returns the endTime and startTime used for querying Azure Metrics API
func GetTimes() (string, string) {
	// Make sure we are using UTC
	now := time.Now().UTC()

	// Use query delay of 3 minutes when querying for latest metric data
	endTime := now.Add(time.Minute * time.Duration(-3)).Format(time.RFC3339)
	startTime := now.Add(time.Minute * time.Duration(-4)).Format(time.RFC3339)
	return endTime, startTime
}

// CreateResourceLabels - Returns resource labels for a give resource ID.
func CreateResourceLabels(resourceID string) map[string]string {
	labels := make(map[string]string)
	labels["resource_group"] = strings.Split(resourceID, "/")[4]
	labels["resource_type"] = strings.Split(resourceID, "/")[6]
	labels["resource_name"] = strings.Split(resourceID, "/")[8]
	return labels
}

func hasAggregation(t config.Target, aggregation string) bool {
	for _, aggr := range t.Aggregations {
		if aggr == aggregation {
			return true
		}
	}
	return false
}
