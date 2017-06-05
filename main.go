package main

import (
	"fmt"
	"os"

	"github.com/logvoyage/logvoyage/cmd"
	"github.com/logvoyage/logvoyage/shared/config"
)

func main() {
	config.InitConfig()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
