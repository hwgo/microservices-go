package cmd

import (
	"github.com/hwgo/customer"
)

var customerCmd *ServerCommand

func init() {
	customerCmd = NewServerCommand(
		customer.ServiceName,
		func(endpoint string) error {
			server := customer.NewServer(endpoint)
			return server.Run()
		},
	)
	customerCmd.AppendTo(RootCmd)
}
