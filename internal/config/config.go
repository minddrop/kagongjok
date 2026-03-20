package config

import "time"

const (
	// HealthCheckTarget is the URL to ping for connectivity checks.
	HealthCheckTarget = "https://www.google.com"

	// CheckInterval is how often to check for connectivity.
	CheckInterval = 10 * time.Second

	// LoginRetryDelay is how long to wait after a failed login attempt or between loops.
	LoginRetryDelay = 5 * time.Second
)
