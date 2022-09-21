package util

import "time"

type TouchRequest struct {
	Id     string
	Target string
}

type TouchResponse struct {
	YourIp          string
	TargetIp        string
	TargetLastTouch time.Time
}
