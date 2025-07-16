package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	assetDb, err := initAssetDB("./assetData.db")
	if err != nil {
		log.Fatal(err)
	}
	defer assetDb.Close()

	assetIds := []string{"1", "1027", "5426", "22861", "32684", "30494", "29587", "26997", "6953", "12220", "7429", "28932", "11396", "30661", "18934", "32492"}
	currencies := []string{"AUD"}
	q := url.Values{}
	q.Add("id", strings.Join(assetIds, ","))
	q.Add("convert", strings.Join(currencies, ","))

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

func initAssetDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Make sure the connection is valid
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ensure table exists
	createTableStmt := `
	CREATE TABLE IF NOT EXISTS assets (
		id INTEGER PRIMARY KEY,
		symbol TEXT NOT NULL,
		price REAL,
		market_cap REAL
	)`
	if _, err := db.Exec(createTableStmt); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}
