/*
Azure Price Search CMD
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/muandane/azureprice/utils"
	"github.com/spf13/cobra"
)

func init() {
	searchCmd.Flags().StringVarP(&vmType, "type", "t", "", "VM type")
	searchCmd.Flags().StringVarP(&region, "region", "r", "", "Region")
	searchCmd.Flags().StringVarP(&service, "service", "s", "", "Azure service (e.g., 'D' for D series vms, Private for Private links)")
	searchCmd.Flags().StringVarP(&pricingType, "pricing-type", "p", "Consumption", "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	searchCmd.Flags().StringVarP(&currency, "currency", "c", "", "Price Currency (e.g., 'USD' or 'EUR')")
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for an Azure resource and get pricing information.",
	Long: `Use the azure search subcommand to search for a specific Azure resource 
and retrieve its pricing information. Provide the resource name 
as an argument to this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		re := lipgloss.NewRenderer(os.Stdout)
		baseStyle := re.NewStyle().Padding(0, 1)
		headerStyle := baseStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#186F65", Dark: "#1AACAC"}).Bold(true)

		tableData := [][]string{{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Region", "Meter", "Product Name"}}
		apiURL := "https://prices.azure.com/api/retail/prices?"
		currencyType := fmt.Sprintf("currencyCode='%s'", currency)
		query := utils.Query(region, service, vmType, pricingType)
		escapedQuery := url.QueryEscape(query)
		for {
			var resp utils.Response
			err := utils.GetJSON(apiURL+currencyType+"&$filter="+escapedQuery, &resp)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, item := range resp.Items {
				var monthlyPrice string
				if pricingType != "Reservation" && !strings.Contains(item.UnitOfMeasure, "GB") && !strings.Contains(item.UnitOfMeasure, "Month") && !strings.Contains(item.UnitOfMeasure, "M") && !strings.Contains(item.UnitOfMeasure, "K") {
					monthlyPrice = fmt.Sprintf("%v", item.RetailPrice*730) // Calculate the monthly price
				} else {
					monthlyPrice = "---"
				}
				tableData = append(tableData, []string{item.ArmSkuName, fmt.Sprintf("%f", item.RetailPrice), item.UnitOfMeasure, fmt.Sprintf("%v", monthlyPrice), item.MeterName, item.ArmRegionName, item.ProductName})
			}
			if resp.NextPageLink == "" {
				break
			}
			apiURL = resp.NextPageLink
		}

		headers := []string{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Meter", "Region", "Product Name"}
		CapitalizeHeaders := func(tableData []string) []string {
			for i := range tableData {
				tableData[i] = strings.ToUpper(tableData[i])
			}
			return tableData
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(re.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#186F65", Dark: "#1AACAC"})).
			Headers(CapitalizeHeaders(headers)...).
			Width(120).
			Rows(tableData[1:]...). // Pass only the rows to the Rows function
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 {
					return headerStyle
				}
				if col == 4 {
					// Check if the "Meter" column contains "Spot" or "Low"
					meter := tableData[row-0][4]                                       // The "Meter" column is the 5th column (index 4)
					color := lipgloss.AdaptiveColor{Light: "#186F65", Dark: "#1AACAC"} // Default color
					if strings.Contains(meter, "Spot") {
						color = typeColors.Spot
					} else if strings.Contains(meter, "Low") {
						color = typeColors.Low
					} else {
						color = typeColors.Normal
					}
					return baseStyle.Copy().Foreground(color)
				}
				return baseStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#053B50", Dark: "#F1EFEF"})
			})
		fmt.Println(t)
	},
}
