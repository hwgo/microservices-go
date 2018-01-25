package cmd

import (
	"fmt"
)

var infoCmd = NewServerCommand(
	"info",
	func(e string) error {
		fmt.Printf("args: %#v", e)
		return nil
	},
)

func init() {
	infoCmd.AppendTo(RootCmd)
}
