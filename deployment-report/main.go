package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
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
	Namespace        string  `json:"namespace"`
	Deployed_by      string  `json:"deployed_by"`
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
		// add more metadata from the event
		stats[event.Service].Deployed_by = event.Deployed_by
		stats[event.Service].Namespace = event.Namespace
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

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Error opening input file: %v", err)
		}
		defer file.Close()
		decoder = json.NewDecoder(file)
	} else if info.Mode()&os.ModeCharDevice != 0 {
		// no data piped in
		fmt.Printf("Usage:\tcat events.json | go run main.go\n")
		fmt.Printf("or...\tgo run main.go events.json\n")
		os.Exit(1)
	} else {
		decoder = json.NewDecoder(os.Stdin)
	}

	if err = decoder.Decode(&events); err != nil {
		log.Fatalf("Error decoding JSON input: %v", err)
	}

	serviceStats := aggregateServiceStats(events)

	for service, stats := range serviceStats {
		fmt.Printf("Service: %s, Total Deployments: %d, Successes: %d, Failures: %d, Average Duration: %.1f seconds. In namespace %s, deployed by %s.\n", service, stats.TotalDeployments, stats.SuccessCount, stats.FailureCount, stats.AvgDuration, stats.Namespace, stats.Deployed_by)
	}

	slowestService := ""

	for service, stats := range serviceStats {
		if slowestService == "" {
			slowestService = service
		}
		if stats.AvgDuration < serviceStats[slowestService].AvgDuration {
			slowestService = service
		}
	}

	fmt.Printf("Slowest service: %s with average duration %.1f seconds. In namespace %s, deployed by %s.\n", slowestService, serviceStats[slowestService].AvgDuration, serviceStats[slowestService].Namespace, serviceStats[slowestService].Deployed_by)
}
