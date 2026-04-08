package main

import (
	"os"
)

const sampleEvent = `[
  { "service": "api-gateway", "namespace": "production", "status": "success", "duration_seconds": 42, "deployed_by": "alice" },
  { "service": "auth-service", "namespace": "production", "status": "failure", "duration_seconds": 130, "deployed_by": "bob" },
  { "service": "api-gateway", "namespace": "staging", "status": "success", "duration_seconds": 18, "deployed_by": "alice" },
  { "service": "payment-service", "namespace": "production", "status": "success", "duration_seconds": 95, "deployed_by": "carol" },
  { "service": "auth-service", "namespace": "staging", "status": "success", "duration_seconds": 60, "deployed_by": "alice" },
  { "service": "payment-service", "namespace": "production", "status": "failure", "duration_seconds": 200, "deployed_by": "bob" },
  { "service": "api-gateway", "namespace": "production", "status": "success", "duration_seconds": 55, "deployed_by": "carol" },
  { "service": "auth-service", "namespace": "production", "status": "success", "duration_seconds": 75, "deployed_by": "alice" }
]`

type deploymentEvent struct {
	Service          string `json:"service"`
	Namespace        string `json:"namespace"`
	Status           string `json:"status"`
	Duration_seconds int    `json:"duration_seconds"`
	Deployed_by      string `json:"deployed_by"`
}

func main() {

	os.Exit(0)

}
