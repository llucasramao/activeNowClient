package main

import (
	"activeNow/config"
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

func main() {
	if runtime.GOOS == "linux" {
		Cron(config.CronTime)
		listen()
	} else {
		logger.Log("Esse agente so pode ser utilizado em sistema linux!", true)
		return
	}
}

func Cron(sched string) {
	logger.Log("Iniciando agente activeNow "+config.AgentVersion+" - "+sched+"\n\n", false)
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
		//functions.GetCommands()
		functions.Finder()
	} else {
		time.Sleep(time.Second * 3)
		logger.Log("Tentando auto atualizar", false)
		autoUpdate()
	}
}

func autoUpdate() {
	cmd := exec.Command("bash", "-c", "curl -fsSL "+config.Manager+"/linuxClient | sudo sh")
	output, err := cmd.Output()
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}
	outputStr := string(output)
	fmt.Println(outputStr)
}

func isUpdate() bool {
	resp, err := http.Get(config.Manager + "/version")
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

	if lastVersion == config.AgentVersion {
		return true
	} else {
		logger.Log("\n\n\nSoftware não atualizado!\nAtual: "+config.AgentVersion+"\nUltima: "+lastVersion+"\n", true)
		return false
	}
}
