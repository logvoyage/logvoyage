package cmd

import (
	"fmt"
	"os"

	"github.com/logvoyage/logvoyage/shared/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "logvoyage",
	Short: "Logvoyage",
	Long:  "Logvoyage - Open Source Log Management",
}

func Execute() {
	config.InitConfig()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
