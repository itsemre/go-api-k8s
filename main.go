package main

import (
	"github.com/itsemre/go-api-k8s/pkg/cmd"
)

// The entrypoint to the API
func main() {
	if err := cmd.Execute(); err != nil {
		cmd.Logger.Panic(err)
		cmd.Logger.ExitFunc(1)
	}
}
