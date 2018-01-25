package cmd

import (
	"github.com/hwgo/frontend"
)

var frontendCmd *ServerCommand

func init() {
	frontendCmd = NewServerCommand(
		frontend.ServiceName,
		func(endpoint string) error {
			server := frontend.NewServer(endpoint)
			return server.Run()
		},
	)
	frontendCmd.AppendTo(RootCmd)
}
