package main

import (
	"activeNow/functions"
	logger "activeNow/log"
	"activeNow/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

const version = "0.3.8-15s-test"
const manager = "http://192.168.1.14:7654"

func main() {
	if runtime.GOOS == "linux" {
		Cron("15s")
		listen()
	} else {
		logger.Log("Esse agente so pode ser utilizado em sistema linux!", true)
		return
	}
}

func Cron(sched string) {
	logger.Log("Iniciando agente activeNow "+version+" - "+sched+"\n\n", false)
	c := cron.New()
	id, _ := c.AddFunc("@every "+sched, func() {
		Init()
	})
	c.Entry(id).Job.Run()
	c.Start()
}

func listen() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	fmt.Println(time.Now().String() + " - Closed")
}

func Init() {
	if isUpdate() {
		functions.GetCommands()
		functions.Finder()
	} else {
		time.Sleep(time.Second * 3)
		logger.Log("Tentando auto atualizar", false)
		autoUpdate()
	}
}

func autoUpdate() {
	cmd := exec.Command("bash", "-c", "curl -fsSL "+manager+"/linuxClient | sudo sh")
	output, err := cmd.Output()
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}
	outputStr := string(output)
	fmt.Println(outputStr)
}

func isUpdate() bool {
	resp, err := http.Get(manager + "/version")
	if err != nil {
		logger.Log(err.Error(), true)
		return false
	}
	defer resp.Body.Close()

	var versionInfo models.VersionInfo
	err = json.NewDecoder(resp.Body).Decode(&versionInfo)
	if err != nil {
		logger.Log(err.Error(), true)
		return false
	}

	lastVersion := strings.TrimSpace(versionInfo.Version)

	if lastVersion == version {
		return true
	} else {
		logger.Log("\n\n\nSoftware não atualizado!\nAtual: "+version+"\nUltima: "+lastVersion+"\n", true)
		return false
	}
}
