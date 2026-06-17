package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type ServicesStatus string

const (
	ServicesStatusOk       ServicesStatus = "ok"
	ServicesStatusDegraded ServicesStatus = "degraded"
)

type ServiceStatus struct {
	Name   string     `json:"name"`
	Status PingStatus `json:"status"`
}

type ServicesInfo struct {
	Status   ServicesStatus  `json:"status"`
	Services []ServiceStatus `json:"services"`
}
