package main

import (
	logger "activeNow/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

const version = "0.3.5"

func main() {
	if runtime.GOOS != "linux" {
		logger.Log("OK", false)
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
	fmt.Println(1)
	if isUpdate() {
		fmt.Println(2)
		//finder()
	} else {
		fmt.Println(3)
		time.Sleep(time.Second * 3)
		logger.Log("Tentando auto atualizar", false)
		//autoUpdate()
	}
}

func isUpdate() bool {
	resp, err := http.Get("http://192.168.1.31/lastVersion")
	if err != nil {
		logger.Log(err.Error(), true)
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	lastVersion := strings.TrimSpace(string(content))

	if lastVersion == version {
		return true
	} else {
		logger.Log("\n\n\nSoftware não atualizado!\nAtual: "+version+"\nUltima: "+lastVersion+"\n", true)
		return false
	}
}
