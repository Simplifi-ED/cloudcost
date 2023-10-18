package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(azureCmd)
}

var rootCmd = &cobra.Command{
	Use:   "cloudcost",
	Short: "Cloud Costs CLI",
	Long:  `cloudcost is a Go CLI that retrieves pricing information for Azure services using the Azure pricing API.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
