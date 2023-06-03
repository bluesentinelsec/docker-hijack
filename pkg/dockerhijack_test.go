package dockerhijack

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProxyDockerArgs(t *testing.T) {
	assert := assert.New(t)

	args := []string{"docker", "--help"}

	ret := ProxyDockerArgs(args)

	assert.Equal(ExitSuccess, ret)
}

func TestGetDockerfileArg(t *testing.T) {

	assert := assert.New(t)

	cmd1 := []string{"docker", "build", "."}
	assert.Equal("Dockerfile", ExtractBuildFileFromArgs(cmd1))

	cmd2 := []string{"docker", "build", "-f", "ctx/Dockerfile", "http://server/ctx.tar.gz"}
	assert.Equal("ctx/Dockerfile", ExtractBuildFileFromArgs(cmd2))

	cmd3 := []string{"docker", "build", "-", "<", "Dockerfile"}
	assert.Equal("Dockerfile", ExtractBuildFileFromArgs(cmd3))

}

func TestDoHijack(t *testing.T) {
	assert := assert.New(t)

	err := os.Chdir("../testData")
	assert.Nil(err)

	cmd := []string{"docker", "build", "."}
	DoHijack(cmd)
}
