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
	vmType := flag.String("t", "", "VM type")
	region := flag.String("r", "", "Region")
	pricingPeriod := flag.String("p", "", "Pricing period (hour or month)")
	series := flag.String("s", "", "VM series (e.g., 'D' for D series)")
	flag.StringVar(series, "series", "", "VM series (e.g., 'D' for D series)")
	flag.StringVar(vmType, "type", "", "VM type")
	flag.StringVar(region, "region", "", "Region")
	flag.StringVar(pricingPeriod, "pricing", "", "Pricing period (hour or month)")
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
	if *series != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and contains(armSkuName, '%s') and priceType eq 'Consumption'", *region, *series)
	} else if *vmType != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and armSkuName eq '%s' and priceType eq 'Consumption'", *region, *vmType)
	} else {
		fmt.Println("Please provide either a series or type flag.")
		return
	}

	tableData := [][]string{{"SKU", "Retail Price", "Unit of Measure", "Region", "Meter", "Product Name"}}
	apiURL := "https://prices.azure.com/api/retail/prices?"

	for {
		var resp Response
		err := getJSON(apiURL+"&$filter="+url.QueryEscape(query), &resp)
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
