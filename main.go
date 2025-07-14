package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type ApiResponse struct {
	Data map[string]AssetResponse `json:"data"`
}

type AssetResponse struct {
	Id     int              `json:"id"`
	Symbol string           `json:"symbol"`
	Quote  map[string]Quote `json:"quote"`
}

type Quote struct {
	Price     float64 `json:"price"`
	MarketCap float64 `json:"market_cap"`
}

type AssetDetail struct {
	Id        int     `json:"id"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	MarketCap float64 `json:"market_cap"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	assetIds := []string{"1", "1027", "5426", "22861", "32684", "30494", "29587", "26997", "6953", "12220", "7429", "28932", "11396", "30661", "18934", "32492"}
	currencies := []string{"AUD"}
	q := url.Values{}
	q.Add("id", strings.Join(assetIds, ","))
	q.Add("convert", strings.Join(currencies, ","))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b25a7f35-2400-4192-84d7-71dba52a2cdd")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	defer resp.Body.Close()

	var response ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding JSON", err)
		os.Exit(1)
	}

	var result []AssetDetail
	for _, asset := range response.Data {
		audQuote, ok := asset.Quote["AUD"]

		if !ok {
			continue
		}

		result = append(result, AssetDetail{
			Id:        asset.Id,
			Symbol:    asset.Symbol,
			Price:     audQuote.Price,
			MarketCap: audQuote.MarketCap,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].MarketCap > result[j].MarketCap
	})

	for _, asset := range result {
		fmt.Printf("%s (%d): $%.2f\t$%.2f\n", asset.Symbol, asset.Id, asset.Price, asset.MarketCap)
	}
}
