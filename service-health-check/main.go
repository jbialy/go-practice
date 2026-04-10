package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ServiceData struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type ServiceStats struct {
	TotalChecks      int
	HealthyCount     int
	UnhealthyCount   int
	UptimePercentage float64
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("Usage: service-health-check <datafile>")
	}

	dataFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer dataFile.Close()

	decoder := json.NewDecoder(dataFile)

	var services []ServiceData
	if err := decoder.Decode(&services); err != nil {
		log.Fatalf("Error decoding JSON data: %v", err)
	}

	stats := map[string]*ServiceStats{}

	for _, service := range services {
		if stats[service.Service] == nil {
			stats[service.Service] = &ServiceStats{}
		}
		stats[service.Service].TotalChecks++
		// increment healthy/unhealthy counts based on status
		switch service.Status {
		case "healthy":
			stats[service.Service].HealthyCount++
		case "unhealthy":
			stats[service.Service].UnhealthyCount++
		}
		// calculate uptime percentage
		stats[service.Service].UptimePercentage = float64(stats[service.Service].HealthyCount) / float64(stats[service.Service].TotalChecks) * 100
	}

	for service, result := range stats {
		fmt.Printf("Service: %s, Total Checks: %d, Healthy: %d, Unhealthy: %d, Uptime: %.2f%%\n",
			service, result.TotalChecks, result.HealthyCount, result.UnhealthyCount, result.UptimePercentage)
	}

}
