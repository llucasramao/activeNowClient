package main

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"log"
	"os"
	"runtime"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"net"
	"time"
	"os/signal"
	
	"github.com/robfig/cron/v3"
)

const version = "0.3.5"

func main() {
	if runtime.GOOS == "linux" {
			Cron("1h")
			listen()
	} else {
		Log("Esse agente so pode ser utilizado em sistema linux!", true)
		return
	}
}

func Cron(sched string){
	Log("\n\nIniciando agente activeNow "+version+" - "+sched, false)
	c := cron.New()
	id, _ := c.AddFunc("@every "+sched, func() {
		Init()
	})
	c.Entry(id).Job.Run()
	c.Start()
}

func Init(){
	if isUpdate(){
		finder()
	} else {
		time.Sleep(time.Second * 3)
		Log("Tentando auto atualizar", false)
		autoUpdate()
	}
}

func listen() {
	// Função cirada para manter o cron ininterrupto
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	fmt.Println(time.Now().String() + " - Closed")
}

func Log(message string, isError bool){
	logDir := exec.Command("mkdir", "-p", "/var/log/activeNow")
	_, err := logDir.Output()
	if err != nil {
		Log(err.Error(), true)
		return
	}
	file, err := os.OpenFile("/var/log/activeNow/history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)

	if isError {
		message = "[-] ERROR: " + message
	} else {
		message = "[+] INFO: " + message
	}

	logger.Println(message)
	fmt.Println(message)
}

func isUpdate() bool{
	resp, err := http.Get("http://192.168.1.31/lastVersion")
	if err != nil{
		Log(err.Error(), true)
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	lastVersion := strings.TrimSpace(string(content))
	
	if lastVersion == version{
		return true
	} else {
		Log("\n\n\nSoftware não atualizado!\nAtual: "+version+"\nUltima: "+lastVersion+"\n", true)
		return false
	}
}

func autoUpdate(){
	cmd := exec.Command("bash", "-c", "curl -fsSL http://192.168.1.31/install.sh | sudo sh")
	output, err := cmd.Output()
	if err != nil {
		Log(err.Error(), true)
		return
	}
	outputStr := string(output)
	fmt.Println(outputStr)
}

func finder(){
  	Log("Buscando softwares e versoes instaladas COMMAND: 'dpkg -l'", false)
	cmd := exec.Command("dpkg", "-l")
	output, err := cmd.Output()
	if err != nil {
		Log(err.Error(), true)
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

			if strings.Contains(software, "Err?=(none)") || (strings.Contains(software, "Name")){
				Log("Pulando "+software+":"+version,false)
				continue
			}
			newObject := map[string]string{
        		"software": software,
        		"version": version,
      		}

      		findings = append(findings, newObject)
		}
	}

	requestBody := Body{
		Ip:findIP(),
    	Ports:findPorts(),
		Hostname:findHostname(),
		Os:findOS(),
		Findings:findings,
	}
	postRequest("https://webhook.site/36cc7ae1-ddca-4fcb-bb75-896b43540fb0", requestBody)
}

func findIP() string{
	Log("Buscando IP eth0 da maquina", false)
	iface, err := net.InterfaceByName("eth0")
    if err != nil {
        Log(err.Error(), true)
        return "nil"
    }

	addrs, err := iface.Addrs()
    if err != nil {
        Log(err.Error(), true)
        return "nil"
    }

    // Percorrer os endereços e imprimir o primeiro que for IPv4
    for _, addr := range addrs {
        ip, _, err := net.ParseCIDR(addr.String())
        if err != nil {
            Log(err.Error(), true)
            continue
        }
        if ip.To4() != nil {
			return ip.String()
            break
        }
    }
	return "ok"
}

func findHostname() string{
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		Log(err.Error(), true)
		return "nil"
	} else {
		return hostname
	}
}

func findOS() string{
	Log("Funcao de buscar OS em manutencao", false)
	return "nil"
}

func findPorts() []int{
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
	Ip          string              `json:"ip"`
	Hostname	string				`json:"hostname"`
	Os			string 				`json:"os"`
	Ports []int `json:"ports"`
	Findings []map[string]string 	`json:"findings"`
	
  
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
		Log(err.Error(), true)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		Log(err.Error(), true)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200{
		Log("Erro ao fazer requisição", true)
		content, _ := ioutil.ReadAll(resp.Body)
		Log(string(content), true)
	} else {
		Log("Dados enviados a API", false)
		return
	}

}
