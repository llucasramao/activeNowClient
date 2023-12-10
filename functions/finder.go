package functions

import (
	logger "activeNow/log"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Finder() {
	logger.Log("Buscando softwares e versoes instaladas COMMAND: 'dpkg -l'", false)
	cmd := exec.Command("dpkg", "-l")
	output, err := cmd.Output()
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	findings := []map[string]string{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			software := fields[1]
			version := fields[2]

			if strings.Contains(software, "Err?=(none)") || (strings.Contains(software, "Name")) {
				logger.Log("Pulando "+software+":"+version, false)
				continue
			}
			newObject := map[string]string{
				"software": software,
				"version":  version,
			}

			findings = append(findings, newObject)
		}
	}

	requestBody := Body{
		Ip:       findIP(),
		Ports:    findPorts(),
		Hostname: findHostname(),
		Os:       findOS(),
		Findings: findings,
	}
	postRequest("https://192.168.1.14:7654/receiver", requestBody)
}

func findIP() string {
	logger.Log("Buscando IP eth0 da maquina", false)
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		logger.Log(err.Error(), true)
		return "nil"
	}

	addrs, err := iface.Addrs()
	if err != nil {
		logger.Log(err.Error(), true)
		return "nil"
	}

	// Percorrer os endereços e imprimir o primeiro que for IPv4
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			logger.Log(err.Error(), true)
			continue
		}
		if ip.To4() != nil {
			return ip.String()
			break
		}
	}
	return "ok"
}

func findHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		logger.Log(err.Error(), true)
		return "nil"
	} else {
		return hostname
	}
}

func findOS() string {
	logger.Log("Funcao de buscar OS em manutencao", false)
	return "nil"
}

func findPorts() []int {
	var openPorts []int
	for port := 1; port <= 1024; port++ {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		} else {
			openPorts = append(openPorts, port)
			conn.Close()
		}
	}
	return openPorts
}

type Body struct {
	Ip       string              `json:"ip"`
	Hostname string              `json:"hostname"`
	Os       string              `json:"os"`
	Ports    []int               `json:"ports"`
	Findings []map[string]string `json:"findings"`
}

func postRequest(url string, requestBody Body) {
	//const url = "https://webhook.site/36cc7ae1-ddca-4fcb-bb75-896b43540fb0"

	// requestBody := []map[string]string{
	// 	{
	// 		"teste":  "teste",
	// 		"teste2": "teste2",
	// 	},
	// }

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.Log(err.Error(), true)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Log("Erro ao fazer requisição", true)
		content, _ := ioutil.ReadAll(resp.Body)
		logger.Log(string(content), true)
	} else {
		logger.Log("Dados enviados a API", false)
		return
	}
}
