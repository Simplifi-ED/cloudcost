/*
Azure Price Calculator CMD
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

// calculatorCmd represents the calculator command
var calculatorCmd = &cobra.Command{
	Use:   "calculator",
	Short: "Calculate Azure resource pricing based on parameters.",
	Long: `Use the azure calculator subcommand to calculate the pricing of an Azure resource.
You can specify the resource name and additional parameters to get accurate pricing details.`,
	Run: func(cmd *cobra.Command, args []string) {
		re := lipgloss.NewRenderer(os.Stdout)
		baseStyle := re.NewStyle().Padding(0, 1)
		headerStyle := baseStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#186F65", Dark: "#1AACAC"}).Bold(true)

		tableData := [][]string{{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Usage", "Region", "Product Name"}}
		apiURL := "https://prices.azure.com/api/retail/prices?"
		currencyType := fmt.Sprintf("currencyCode='%s'", currency)
		query := utils.Query(region, service, vmType, pricingType)
		for {
			var resp utils.Response
			err := utils.GetJSON(apiURL+currencyType+"&$filter="+url.QueryEscape(query), &resp)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, item := range resp.Items {
				var usage float64
				if strings.Contains(item.UnitOfMeasure, "GB") {
					usage = calculateUsageGB(bandwidth, period, item.RetailPrice) // Assuming a bandwidth of 1 GB/day and a span of 30 days
				} else if strings.Contains(item.UnitOfMeasure, "Hour") {
					usage = calculateUsageHourly(bandwidth, period, item.RetailPrice) // Assuming a bandwidth of 1 GB/day and
				} else if strings.Contains(item.UnitOfMeasure, "Month") {
					usage = calculateUsageMonthly(bandwidth, period, item.RetailPrice) // Assuming a bandwidth of 1 GB/day and
				} else if strings.Contains(item.UnitOfMeasure, "M") {
					usage = calculateUsageEvents(eventCount, item.RetailPrice)
				}

				var monthlyPrice string
				if pricingType != "Reservation" && !strings.Contains(item.UnitOfMeasure, "GB") && !strings.Contains(item.UnitOfMeasure, "Month") && !strings.Contains(item.UnitOfMeasure, "M") && !strings.Contains(item.UnitOfMeasure, "K") {
					monthlyPrice = fmt.Sprintf("%v", item.RetailPrice*730) // Calculate the monthly price
				} else {
					monthlyPrice = "---"
				}
				tableData = append(tableData, []string{item.ArmSkuName, fmt.Sprintf("%f", item.RetailPrice), item.UnitOfMeasure, fmt.Sprintf("%v", monthlyPrice), fmt.Sprintf("%f", usage), item.ArmRegionName, item.MeterName, item.ProductName})
			}
			if resp.NextPageLink == "" {
				break
			}
			apiURL = resp.NextPageLink
		}
		headers := []string{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Usage", "Region", "Meter Name", "Product Name"}
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

func init() {
	calculatorCmd.Flags().StringVarP(&vmType, "type", "t", "", "VM type")
	calculatorCmd.Flags().StringVarP(&region, "region", "r", "", "Region")
	calculatorCmd.Flags().StringVarP(&service, "service", "s", "", "Azure service (e.g., 'D' for D series vms, Private for Private links)")
	calculatorCmd.Flags().StringVarP(&pricingType, "pricing-type", "p", "Consumption", "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	calculatorCmd.Flags().StringVarP(&currency, "currency", "c", "", "Price Currency (e.g., 'USD' or 'EUR')")
	calculatorCmd.Flags().Float64VarP(&bandwidth, "bandwidth", "b", 1, "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	calculatorCmd.Flags().IntVarP(&period, "days", "d", 1, "period (e.g., '1' for 1 day, '7' for 7 days)")
	calculatorCmd.Flags().Float64VarP(&eventCount, "events", "e", 1, "Number of events")
}

func calculateUsageGB(bandwidth float64, days int, usagePerGB float64) float64 {
	totalUsage := bandwidth * float64(days) * usagePerGB
	return totalUsage
}
func calculateUsageEvents(eventCount float64, usagePerEvent float64) float64 {
	totalUsage := eventCount * usagePerEvent
	return totalUsage
}
func calculateUsageHourly(period float64, days int, usagePerHour float64) float64 {
	hours := days * 24 // convert days to hours
	totalUsage := period * float64(hours) * usagePerHour
	return totalUsage
}
func calculateUsageMonthly(period float64, days int, usagePerMonth float64) float64 {
	month := days / 30 // convert days to hours
	totalUsage := period * float64(month) * usagePerMonth
	return totalUsage
}
