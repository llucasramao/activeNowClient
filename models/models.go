package models

type Port struct {
	Port      int
	Receiveds []Received
}

// Reciber representa a tabela de receptores
type Received struct {
	Ip       string
	Hostname string
	Os       string
	Ports    []Port
}
