package main

import (
	"os"

	"kubevault.dev/cli/pkg/cmds"

	"kmodules.xyz/client-go/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
