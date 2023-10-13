package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

func main() {
	vmType := flag.String("t", "Standard_B4ms", "VM type")
	region := flag.String("r", "westus", "Region")
	service := flag.String("s", "", "Azure service (e.g., 'D' for D series vms, Private for Private links)")
	pricingType := flag.String("p", "Consumption", "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	currency := flag.String("c", "USD", "Price Currency (e.g., 'USD' or 'EUR')")
	flag.StringVar(service, "service", "", "Azure service (e.g., 'D' for D series vms, Private for Private links)")
	flag.StringVar(vmType, "type", "Standard_B4ms", "VM type")
	flag.StringVar(region, "region", "westus", "Region")
	flag.StringVar(pricingType, "pricing-type", "Consumption", "Pricing Type (e.g., 'Consumption' or 'Reservation')")
	flag.StringVar(currency, "currency", "USD", "Price Currency (e.g., 'USD' or 'EUR')")
	flag.Parse()

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	typeColors := map[string]lipgloss.Color{
		"Spot":   lipgloss.Color("#FF7698"),
		"Normal": lipgloss.Color("#75FBAB"),
		"Low":    lipgloss.Color("#FDFF90"),
	}

	var query string
	if *service != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and contains(serviceName, '%s')", *region, *service)
	} else if *vmType != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and contains(armSkuName, '%s') and priceType eq '%s'", *region, *vmType, *pricingType)
	} else {
		fmt.Println("Please provide either a series or type flag.")
		return
	}

	tableData := [][]string{{"SKU", "Retail Price", "Unit of Measure", "Region", "Meter", "Product Name"}}
	apiURL := "https://prices.azure.com/api/retail/prices?"
	currencyType := fmt.Sprintf("currencyCode='%s'", *currency)

	for {
		var resp Response
		err := getJSON(apiURL+currencyType+"&$filter="+url.QueryEscape(query), &resp)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, item := range resp.Items {
			tableData = append(tableData, []string{item.ArmSkuName, fmt.Sprintf("%f", item.RetailPrice), item.UnitOfMeasure, item.ArmRegionName, item.MeterName, item.ProductName})
		}
		if resp.NextPageLink == "" {
			break
		}
		apiURL = resp.NextPageLink
	}

	headers := []string{"SKU", "Retail Price", "Unit of Measure", "Region", "Meter", "Product Name"}
	CapitalizeHeaders := func(tableData []string) []string {
		for i := range tableData {
			tableData[i] = strings.ToUpper(tableData[i])
		}
		return tableData
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(CapitalizeHeaders(headers)...).
		Width(100).
		Rows(tableData[1:]...). // Pass only the rows to the Rows function
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			if col == 4 {
				// Check if the "Meter" column contains "Spot" or "Low"
				meter := tableData[row-0][4]   // The "Meter" column is the 5th column (index 4)
				color := lipgloss.Color("252") // Default color

				if strings.Contains(meter, "Spot") {
					color = typeColors["Spot"]
				} else if strings.Contains(meter, "Low") {
					color = typeColors["Low"]
				} else {
					color = typeColors["Normal"]
				}
				return baseStyle.Copy().Foreground(color)
			}
			return baseStyle.Copy().Foreground(lipgloss.Color("252"))
		})
	fmt.Println(t)
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
