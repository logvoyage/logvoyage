package cmd

import (
	"fmt"

	"github.com/logvoyage/logvoyage/api"
	"github.com/logvoyage/logvoyage/consumer"
	"github.com/logvoyage/logvoyage/models"
	"github.com/logvoyage/logvoyage/producer"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: " Starts backend API server",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		models.InitDatabase()

		host := cmd.Flags().Lookup("host")
		port := cmd.Flags().Lookup("port")
		fmt.Println("Starting API server at port", port.Value)
		api.Start(host.Value.String(), port.Value.String())
	},
}

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Starts data consumer",
	Long:  "Consumer will verify all incoming data and send it to storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting consumer")
		consumer.Start()
	},
}

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Starts data producer",
	Long:  "Producer worker accepts user data and sends it to consumer worker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting producer")
		producer.Start()
	},
}

func init() {
	startCmd := &cobra.Command{Use: "start", Short: "Set of commands to work with LogVoyage services"}

	apiCmd.Flags().String("port", "3000", "Port to open for API")
	apiCmd.Flags().String("host", "localhost", "Server host")

	startCmd.AddCommand(
		apiCmd,
		consumerCmd,
		producerCmd,
	)

	RootCmd.AddCommand(startCmd)
}
