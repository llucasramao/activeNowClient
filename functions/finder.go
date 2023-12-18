package functions

import (
	"activeNow/config"
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
	"regexp"
	"strings"
)

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
		Ip:           findIP(),
		Ports:        findPorts(),
		Hostname:     findHostname(),
		Os:           findOS(),
		Apps:         apps,
		Services:     ServicesList(),
		AgentVersion: config.AgentVersion,
	}

	postRequest(managerURL, requestBody)
}

func parseDpkgOutput(output string) []models.App {
	var apps []models.App
	lines := strings.Split(output, "\n")

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

func ServicesList() []models.Service {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--plain", "--no-legend")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil
	}

	re := regexp.MustCompile(`^(\S+)\.service`)
	var services []models.Service
	for _, line := range bytes.Split(out.Bytes(), []byte("\n")) {
		matches := re.FindSubmatch(line)
		if len(matches) > 1 {
			service := models.Service{Name: string(matches[1])}
			services = append(services, service)
		}
	}

	return services
}

func findIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "nil"
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "nil"
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // não é um endereço IPv4
			}

			return ip.String()
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
