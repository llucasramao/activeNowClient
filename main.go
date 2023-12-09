package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Iniciando script de busca de ativos!")
	os := runtime.GOOS
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	switch os {
	case "windows":
		fmt.Println("Windows")
	case "darwin":
		fmt.Println("MAC operating system")
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("%s.\n", os)
	}
}

func searchActivesDebian() {
	fmt.Println("Buscando ativos dpkg")
	//cmd := exec.Command("dpkg", "--list")
}
