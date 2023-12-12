package logger

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Log(message string, isError bool) {
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
