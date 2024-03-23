package models

import "time"

type Options struct {
	Name                string
	Address             string
	Port                int
	ConsulAddress       string
	HealthCheckInterval time.Duration
}
