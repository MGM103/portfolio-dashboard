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

	assetDb, err := initAssetDB("./assetData.db")
	if err != nil {
		log.Fatal(err)
	}
	defer assetDb.Close()

	var newAssetId string
	fmt.Println("Enter an asset (id): ")
	fmt.Scanf("%s", &newAssetId)

	err = addAsset(assetDb, newAssetId)
	if err != nil {
		log.Fatal("Adding new asset failed: ", err)
	}

	assets, err := getAllAssets(assetDb)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].MarketCap > assets[j].MarketCap
	})

	for _, asset := range assets {
		fmt.Printf("%s (%d): $%.2f\t$%.2f\n", asset.Symbol, asset.Id, asset.Price, asset.MarketCap)
	}
}

func getAssetData(id string) (AssetDetail, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Values{}
	q.Add("id", id)
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

	audQuote, ok := response.Data[id].Quote["AUD"]
	if !ok {
		return AssetDetail{}, fmt.Errorf("Error with quote.")
	}

	return AssetDetail{
		Id:        response.Data[id].Id,
		Symbol:    response.Data[id].Symbol,
		Price:     audQuote.Price,
		MarketCap: audQuote.MarketCap,
	}, nil
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

func addAsset(db *sql.DB, id string) error {
	assetData, err := getAssetData(id)
	if err != nil {
		return fmt.Errorf("Could not retrieve asset data: %s", err)
	}

	_, err = db.Exec(`INSERT INTO assets (id, symbol, price, market_cap) VALUES (?, ?, ?, ?)`,
		assetData.Id,
		assetData.Symbol,
		assetData.Price,
		assetData.MarketCap,
	)
	if err != nil {
		return fmt.Errorf("Could not insert asset data into db: %s", err)
	}

	return nil
}

func getAsset(db *sql.DB, id string) (AssetDetail, error) {
	var asset AssetDetail

	row := db.QueryRow(`SELECT id, symbol, price, market_cap FROM assets WHERE id = ?`, id)
	err := row.Scan(&asset.Id, &asset.Symbol, &asset.Price, &asset.MarketCap)
	if err != nil {
		if err == sql.ErrNoRows {
			return AssetDetail{}, fmt.Errorf("no asset found with ID %s", id)
		}
		return AssetDetail{}, fmt.Errorf("failed to query asset: %w", err)
	}

	return asset, nil
}

func getAllAssets(db *sql.DB) ([]AssetDetail, error) {
	rows, err := db.Query(`SELECT * FROM assets`)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []AssetDetail

	for rows.Next() {
		var asset AssetDetail
		if err := rows.Scan(&asset.Id, &asset.Symbol, &asset.Price, &asset.MarketCap); err != nil {
			return nil, fmt.Errorf("failed to scan asset row: %w", err)
		}
		assets = append(assets, asset)
	}

	// check for row iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return assets, nil
}
