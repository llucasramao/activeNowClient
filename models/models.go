package models

type VersionInfo struct {
	Version string `json:"version"`
}

type Port struct {
	Port      int
	Receiveds []Received
}

type App struct {
	Name      string
	Version   string
	Receiveds []Received
}

type Service struct {
	Name     string
	Received []Received
}

type Received struct {
	Ip           string
	Hostname     string
	Os           string
	Ports        []Port
	Apps         []App
	Services     []Service
	AgentVersion string
}
