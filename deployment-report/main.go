package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

const sampleEvents = `[
  { "service": "api-gateway", "namespace": "production", "status": "success", "duration_seconds": 42, "deployed_by": "alice" },
  { "service": "auth-service", "namespace": "production", "status": "failure", "duration_seconds": 130, "deployed_by": "bob" },
  { "service": "api-gateway", "namespace": "staging", "status": "success", "duration_seconds": 18, "deployed_by": "alice" },
  { "service": "payment-service", "namespace": "production", "status": "success", "duration_seconds": 95, "deployed_by": "carol" },
  { "service": "auth-service", "namespace": "staging", "status": "success", "duration_seconds": 60, "deployed_by": "alice" },
  { "service": "payment-service", "namespace": "production", "status": "failure", "duration_seconds": 200, "deployed_by": "bob" },
  { "service": "api-gateway", "namespace": "production", "status": "success", "duration_seconds": 55, "deployed_by": "carol" },
  { "service": "auth-service", "namespace": "production", "status": "success", "duration_seconds": 75, "deployed_by": "alice" }
]`

type DeploymentEvent struct {
	Service          string `json:"service"`
	Namespace        string `json:"namespace"`
	Status           string `json:"status"`
	Duration_seconds int    `json:"duration_seconds"`
	Deployed_by      string `json:"deployed_by"`
}

type DeploymentStats struct {
	TotalDeployments int     `json:"total_deployments"`
	SuccessCount     int     `json:"success_count"`
	FailureCount     int     `json:"failure_count"`
	AvgDuration      float64 `json:"avg_duration"`
	TotalDuration    float64 `json:"-"`
}

func aggregateServiceStats(events []DeploymentEvent) map[string]*DeploymentStats {
	stats := map[string]*DeploymentStats{}
	for _, event := range events {
		if stats[event.Service] == nil {
			stats[event.Service] = &DeploymentStats{}
		}
		stats[event.Service].TotalDeployments++
		if event.Status == "success" {
			stats[event.Service].SuccessCount++
		}
		if event.Status == "failure" {
			stats[event.Service].FailureCount++
		}
		stats[event.Service].TotalDuration += float64(event.Duration_seconds)
	}

	// compute average duration per service in DeploymentStats
	for _, s := range stats {
		s.AvgDuration = math.Round(s.TotalDuration/float64(s.TotalDeployments)*10) / 10

	}

	return stats
}

func main() {

	var events []DeploymentEvent

	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Error checking standard input: %v", err)
	}
	// decoder is a json decoder used for reading in data
	var decoder *json.Decoder

	// check for piped input on stdin
	if info.Mode()&os.ModeCharDevice != 0 {
		decoder = json.NewDecoder(strings.NewReader(sampleEvents))
	} else {
		decoder = json.NewDecoder(os.Stdin)
	}

	if err = decoder.Decode(&events); err != nil {
		log.Fatalf("Error decoding JSON input: %v", err)
	}

	serviceStats := aggregateServiceStats(events)

	jsonOutput, err := json.Marshal(serviceStats)
	fmt.Printf("%s\n", jsonOutput)

	os.Exit(0)

}
