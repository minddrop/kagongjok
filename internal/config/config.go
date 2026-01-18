package config

import "time"

const (
	// HealthCheckTarget is the URL to ping for connectivity checks.
	HealthCheckTarget = "https://www.google.com"

	// StarbucksLoginURLTemplate is the format string for the login redirect URL.
	StarbucksLoginURLTemplate = "https://service.wi2.ne.jp/wi2auth/redirect?cmd=login&mac=%s&essid=%%20&apname=tunnel%%201&apgroup=&url=http%%3A%%2F%%2Fexample%%2Ecom%%2F%%3F%d"

	// CheckInterval is how often to check for connectivity.
	CheckInterval = 10 * time.Second

	// LoginRetryDelay is how long to wait after a failed login attempt or between loops.
	LoginRetryDelay = 5 * time.Second
)
