package main

import (
	"fmt"
	dockerhijack "github.com/bluesentinelsec/docker-hijack/pkg"
	"os"
)

func main() {
	if isDockerHijack(os.Args) {
		fmt.Println("you executed docker-hijack")
	} else {
		dockerhijack.ProxyDockerArgs(os.Args)
	}
}

func isDockerHijack(osArgs []string) bool {
	for _, arg := range osArgs {
		if arg == "--infected" {
			return true
		}
	}
	return false
}
