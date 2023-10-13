package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

type Item struct {
	ArmRegionName string  `json:"armRegionName"`
	ArmSkuName    string  `json:"armSkuName"`
	MeterName     string  `json:"meterName"`
	ProductName   string  `json:"productName"`
	RetailPrice   float64 `json:"retailPrice"`
	UnitOfMeasure string  `json:"unitOfMeasure"`
}

type Response struct {
	Items        []Item `json:"Items"`
	NextPageLink string `json:"NextPageLink"`
}

var vmType string
var region string
var service string
var pricingType string
var currency string

var rootCmd = &cobra.Command{
	Use:   "azureprice",
	Short: "Azure Prices CLI",
	Long:  `azureprice is a Go CLI that retrieves pricing information for Azure services using the Azure pricing API.`,
	Run: func(cmd *cobra.Command, args []string) {
		re := lipgloss.NewRenderer(os.Stdout)
		baseStyle := re.NewStyle().Padding(0, 1)
		headerStyle := baseStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#186F65", Dark: "#1AACAC"}).Bold(true)
		typeColors := map[string]lipgloss.AdaptiveColor{
			"Spot":   lipgloss.AdaptiveColor{Light: "#D83F31", Dark: "#D83F31"},
			"Normal": lipgloss.AdaptiveColor{Light: "#116D6E", Dark: "#00DFA2"},
			"Low":    lipgloss.AdaptiveColor{Light: "#EE9322", Dark: "#E9B824"},
		}

		var query string
		if service != "" {
			query = fmt.Sprintf("armRegionName eq '%s' and contains(serviceName, '%s')", region, service)
		} else if vmType != "" {
			query = fmt.Sprintf("armRegionName eq '%s' and contains(armSkuName, '%s') and priceType eq '%s'", region, vmType, pricingType)
		} else {
			fmt.Println("Please provide either a series or type flag.")
			return
		}

		tableData := [][]string{{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Region", "Meter", "Product Name"}}
		apiURL := "https://prices.azure.com/api/retail/prices?"
		currencyType := fmt.Sprintf("currencyCode='%s'", currency)

		for {
			var resp Response
			err := getJSON(apiURL+currencyType+"&$filter="+url.QueryEscape(query), &resp)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, item := range resp.Items {
				var monthlyPrice string
				if pricingType != "Reservation" {
					monthlyPrice = fmt.Sprintf("%v", item.RetailPrice*730) // Calculate the monthly price
				} else {
					monthlyPrice = "---"
				}
				tableData = append(tableData, []string{item.ArmSkuName, fmt.Sprintf("%f", item.RetailPrice), item.UnitOfMeasure, fmt.Sprintf("%v", monthlyPrice), item.ArmRegionName, item.MeterName, item.ProductName})
			}
			if resp.NextPageLink == "" {
				break
			}
			apiURL = resp.NextPageLink
		}

		headers := []string{"SKU", "Retail Price", "Unit of Measure", "Monthly Price", "Region", "Meter", "Product Name"}
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
						color = typeColors["Spot"]
					} else if strings.Contains(meter, "Low") {
						color = typeColors["Low"]
					} else {
						color = typeColors["Normal"]
					}
					return baseStyle.Copy().Foreground(color)
				}
				return baseStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#053B50", Dark: "#F1EFEF"})
			})
		fmt.Println(t)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&vmType, "type", "t", "Standard_B4ms", "VM type")
	rootCmd.Flags().StringVarP(&region, "region", "r", "westus", "Region")
	rootCmd.Flags().StringVarP(&service, "service", "s", "", "Azure service (e.g., 'D' for D series vms, Private for Private links)")
	rootCmd.Flags().StringVarP(&pricingType, "pricing-type", "p", "Consumption", "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	rootCmd.Flags().StringVarP(&currency, "currency", "c", "USD", "Price Currency (e.g., 'USD' or 'EUR')")
}

func getJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
