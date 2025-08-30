package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	Ticker    string  `json:"symbol"`
	Price     float64 `json:"price"`
	MarketCap float64 `json:"market_cap"`
}

func GetAssetData(ids []string, currency string) ([]AssetDetail, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Values{}
	q.Add("id", strings.Join(ids, ","))
	q.Add("convert", "AUD")

	ApiKey := os.Getenv("CMC_API_KEY")
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", ApiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var response ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}

	var result []AssetDetail
	for _, priceResponse := range response.Data {
		assetQuote, ok := priceResponse.Quote[currency]
		if !ok {
			fmt.Println("Failed to get asset.")
		}

		result = append(result, AssetDetail{Id: priceResponse.Id, Ticker: priceResponse.Symbol, Price: assetQuote.Price, MarketCap: assetQuote.MarketCap})
	}

	return result, nil
}
