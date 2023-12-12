package functions

import (
	logger "activeNow/log"
	"activeNow/models"
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

// Finder searches for installed software and versions using the 'dpkg -l' command.
func Finder() {
	const dpkgCommand = "dpkg"
	const dpkgArgs = "-l"
	const managerURL = "http://192.168.1.14:7654/NewReceived"

	logger.Log("Buscando softwares e versões instaladas COMMAND: 'dpkg -l'", false)

	cmd := exec.Command(dpkgCommand, dpkgArgs)
	output, err := cmd.Output()
	if err != nil {
		logger.Log(fmt.Sprintf("Erro ao executar comando: %v", err), true)
		return
	}

	apps := parseDpkgOutput(string(output))

	requestBody := models.Received{
		Ip:       findIP(),
		Ports:    findPorts(),
		Hostname: findHostname(),
		Os:       findOS(),
		Apps:     apps,
	}

	postRequest(managerURL, requestBody)
}

func parseDpkgOutput(output string) []models.App {
	var apps []models.App
	lines := strings.Split(output, "")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			software, version := fields[1], fields[2]

			if strings.Contains(software, "Err?=(none)") || strings.Contains(software, "Name") {
				continue
			}

			apps = append(apps, models.App{Name: software, Version: version})
		}
	}
	return apps
}

func findIP() string {
	logger.Log("Buscando IP eth0 da máquina", false)
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

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			logger.Log(err.Error(), true)
			continue
		}
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return "nil"
}

func findHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Log(err.Error(), true)
		return "nil"
	}
	return hostname
}

func findOS() string {
	logger.Log("Função de buscar OS em manutenção", false)
	// Potentially more logic can be added here for different OS detection
	return "linux"
}

func findPorts() []models.Port {
	var openPorts []models.Port
	for port := 1; port <= 65500; port++ {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}
		openPorts = append(openPorts, models.Port{Port: port})
		conn.Close()
	}
	return openPorts
}

func postRequest(url string, requestBody models.Received) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.Log(fmt.Sprintf("Erro ao serializar JSON: %v", err), true)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log(fmt.Sprintf("Erro ao fazer requisição HTTP: %v", err), true)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log("Erro ao fazer requisição", true)
		content, _ := ioutil.ReadAll(resp.Body)
		logger.Log(string(content), true)
	} else {
		logger.Log("Dados enviados a API com sucesso", false)
	}
}
