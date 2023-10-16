package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GetJSON fetches JSON data from a URL and decodes it into the provided interface
func GetJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %s: %s", resp.Status, string(bodyBytes))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// Query creates the query string for the Azure pricing API
func Query(region string, service string, vmType string, pricingType string) string {
	var query string
	if service != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and contains(serviceName, '%s')", region, service)
	} else if vmType != "" {
		query = fmt.Sprintf("armRegionName eq '%s' and contains(armSkuName, '%s') and priceType eq '%s'", region, vmType, pricingType)
	}
	return query
}
