package cmd

import (
	"fmt"

	"github.com/logvoyage/logvoyage/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: " Starts backend API server",
	Long:  "",

	Run: func(cmd *cobra.Command, args []string) {
		host := cmd.Flags().Lookup("host")
		port := cmd.Flags().Lookup("port")
		fmt.Println("Starting API server at port", port.Value)
		api.Start(host.Value.String(), port.Value.String())
	},
}

func init() {
	startCmd := &cobra.Command{Use: "start", Short: "Set of commands to work with LogVoyage services"}

	apiCmd.Flags().String("port", "3000", "Port to open for API")
	apiCmd.Flags().String("host", "localhost", "Server host")

	startCmd.AddCommand(
		apiCmd,
	)

	RootCmd.AddCommand(startCmd)
}
