package main

import (
	"github.com/tarkalabs/tarkacoin/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
