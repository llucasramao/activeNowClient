package base

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func base() {
	// Registra a função buscarSoftware para a URL /buscar
	http.HandleFunc("/buscar", buscarSoftware)

	// Inicia o servidor na porta 8080
	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarSoftware(w http.ResponseWriter, r *http.Request) {
	// Executa o comando dpkg --list para listar os softwares e versões
	cmd := exec.Command("dpkg", "--list")
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return
	}

	// Converte a saída do comando em um slice de strings
	softwares := strings.Split(string(output), "\n")

	// Cria um buffer de bytes com os dados dos softwares e versões
	var data bytes.Buffer
	for _, software := range softwares {
		data.WriteString(software + "\n")
	}

	// Envia uma requisição POST para o servidor remoto com os dados
	resp, err := http.Post("http://example.com/receber", "text/plain", &data)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta e imprime uma mensagem
	if resp.StatusCode == http.StatusOK {
		fmt.Fprintln(w, "Requisição enviada com sucesso!")
	} else {
		fmt.Fprintln(w, "Ocorreu um erro ao enviar a requisição.")
	}
}
