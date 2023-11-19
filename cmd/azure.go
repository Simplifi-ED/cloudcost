/*
Azure Cloud Costs CMD
*/
package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type Colors struct {
	Spot   lipgloss.AdaptiveColor
	Normal lipgloss.AdaptiveColor
	Low    lipgloss.AdaptiveColor
}

var vmType string
var region string
var service string
var pricingType string
var currency string
var period int
var bandwidth float64
var eventCount float64
var typeColors = Colors{
	Spot:   lipgloss.AdaptiveColor{Light: "#D83F31", Dark: "#D83F31"},
	Normal: lipgloss.AdaptiveColor{Light: "#116D6E", Dark: "#00DFA2"},
	Low:    lipgloss.AdaptiveColor{Light: "#EE9322", Dark: "#E9B824"},
}

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
