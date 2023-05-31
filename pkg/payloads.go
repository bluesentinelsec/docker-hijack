package dockerhijack

import (
	_ "embed"
)

const PAYLOAD_NAME = "meterpreter"

//go:embed payload/build_commands.txt
var installPayload []byte

//go:embed payload/meterpreter
var payload []byte
