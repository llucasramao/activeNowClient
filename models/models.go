package models

type Port struct {
	Port      int
	Receiveds []Received
}

type Received struct {
	Ip       string
	Hostname string
	Os       string
	Ports    []Port
}

type VersionInfo struct {
	Version string `json:"version"`
}
