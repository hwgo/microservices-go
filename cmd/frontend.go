package cmd

import (
	"github.com/hwgo/frontend"
)

var frontendCmd *ServerCommand

func init() {
	name := "frontend"
	frontendCmd = NewServerCommand(
		name,
		func(endpoint string) error {
			server := frontend.NewServer(name, endpoint)
			return server.Run()
		},
	)
	frontendCmd.AppendTo(RootCmd)
}
