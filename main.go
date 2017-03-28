package main

import (
	"fmt"
	"os"

	"bitbucket.org/firstrow/logvoyage/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
