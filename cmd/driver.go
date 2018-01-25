package cmd

import (
	"github.com/hwgo/driver"
)

var driverCmd *ServerCommand

func init() {
	driverCmd = NewServerCommand(
		driver.ServiceName,
		func(endpoint string) error {
			server := driver.NewServer(endpoint)
			return server.Run()
		},
	)
	driverCmd.AppendTo(RootCmd)
}
