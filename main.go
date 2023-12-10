package main

import (
	"activeNow/functions"
	logger "activeNow/log"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

const version = "0.3.5"
const manager = "http://localhost:7654"

func main() {
	if runtime.GOOS == "linux" {
		Cron("5s")
		//listen()
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

func Init() {
	if isUpdate() {
		functions.Finder()
	} else {
		time.Sleep(time.Second * 3)
		logger.Log("Tentando auto atualizar", false)
		//autoUpdate()
	}
}

type VersionInfo struct {
	Version string `json:"version"`
}

func isUpdate() bool {
	resp, err := http.Get(manager + "/version")
	if err != nil {
		logger.Log(err.Error(), true)
		return false
	}
	defer resp.Body.Close()

	var versionInfo VersionInfo
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
