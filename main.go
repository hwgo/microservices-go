package main

import (
	"runtime"

	"github.com/hwgo/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
