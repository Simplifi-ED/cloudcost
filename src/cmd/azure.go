/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var vmType string
var region string
var service string
var pricingType string
var currency string
var period int
var bandwidth float64

// azureCmd represents the azure command
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Main command for Azure resource management.",
	Long: `The main command for interacting with Azure resources. This command has subcommands 
for searching and calculating the pricing of Azure resources`,
}

func init() {
	azureCmd.AddCommand(calculatorCmd)
	azureCmd.AddCommand(searchCmd)
}
