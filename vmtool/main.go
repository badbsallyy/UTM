package main

import (
	"fmt"
	"os"

	"github.com/utmapp/vmtool/cmd/vmtool"
)

func main() {
	if err := vmtool.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
