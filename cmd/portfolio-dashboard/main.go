package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	data "github.com/mgm103/portfolio-dashboard/data"
	"github.com/mgm103/portfolio-dashboard/tui"
)

func main() {
	store := &data.Store{}
	m := tui.NewModel(store)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// import (
// 	"bufio"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"strings"
//
// 	"github.com/joho/godotenv"
// 	_ "github.com/mattn/go-sqlite3"
// )
//
// type Operation int
//
// const (
// 	GetPortfolioValue Operation = iota
// 	AddPosition
// 	RemovePosition
// 	ShowWatchList
// 	AddAssetsWatchList
// 	RemoveAssetsWatchList
// 	Exit
// )
//
// func (o Operation) String() string {
// 	switch o {
// 	case GetPortfolioValue:
// 		return "Portfolio worth"
// 	case AddPosition:
// 		return "Add position(s)"
// 	case RemovePosition:
// 		return "Remove position(s)"
// 	case ShowWatchList:
// 		return "Show watch list"
// 	case AddAssetsWatchList:
// 		return "Add asset(s) to watch list"
// 	case RemoveAssetsWatchList:
// 		return "Remove asset(s) from watch list"
// 	case Exit:
// 		return "Exit"
// 	default:
// 		return "unknown"
// 	}
// }
//
// type ApiResponse struct {
// 	Data map[string]AssetResponse `json:"data"`
// }
//
// type AssetResponse struct {
// 	Id     int              `json:"id"`
// 	Symbol string           `json:"symbol"`
// 	Quote  map[string]Quote `json:"quote"`
// }
//
// type Quote struct {
// 	Price     float64 `json:"price"`
// 	MarketCap float64 `json:"market_cap"`
// }
//
// type AssetDetail struct {
// 	Id        int     `json:"id"`
// 	Symbol    string  `json:"symbol"`
// 	Price     float64 `json:"price"`
// 	MarketCap float64 `json:"market_cap"`
// }
//
// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	assetDb, err := initAssetDB("./assetData.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer assetDb.Close()
//
// programLoop:
// 	for {
// 		currentOperation, err := readOperatonInput()
// 		if err != nil {
// 			fmt.Println("An error occurred: ", err)
// 			continue
// 		}
//
// 		switch currentOperation {
// 		case GetPortfolioValue:
// 			err = getPortfolioValue(assetDb)
// 			if err != nil {
// 				fmt.Println("Could not print portfolio value")
// 			}
// 			fmt.Println()
// 		case AddPosition:
// 			reader := bufio.NewReader(os.Stdin)
// 			fmt.Println("Enter an asset(s) (id): ")
// 			positionString, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Println("Could not read inputted assets")
// 				continue programLoop
// 			}
// 			positionString = strings.TrimSpace(positionString)
//
// 			err = addPosition(assetDb, positionString)
// 			if err != nil {
// 				fmt.Println("Could not add position: ", err)
// 				continue programLoop
// 			}
// 		case RemovePosition:
// 			reader := bufio.NewReader(os.Stdin)
// 			fmt.Println("Enter an asset(s) (id): ")
// 			positionString, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Println("Could not read inputted assets to remove")
// 				continue programLoop
// 			}
// 			positionString = strings.TrimSpace(positionString)
//
// 			err = removePosition(assetDb, positionString)
// 			if err != nil {
// 				fmt.Println("Could not remove position: ", err)
// 				continue programLoop
// 			}
// 		case ShowWatchList:
// 			assets, err := getAllAssets(assetDb)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			sort.Slice(assets, func(i, j int) bool {
// 				return assets[i].MarketCap > assets[j].MarketCap
// 			})
//
// 			fmt.Println()
// 			for _, asset := range assets {
// 				fmt.Printf("%s (%d): $%.2f\t$%.2f\n", asset.Symbol, asset.Id, asset.Price, asset.MarketCap)
// 			}
//
// 			fmt.Println()
// 		case AddAssetsWatchList:
// 			reader := bufio.NewReader(os.Stdin)
// 			fmt.Println("Enter an asset(s) (id): ")
// 			assetsString, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Println("Could not read inputted assets")
// 				continue programLoop
// 			}
// 			assetsString = strings.TrimSpace(assetsString)
//
// 			err = addAssets(assetDb, assetsString)
// 			if err != nil {
// 				log.Fatal("Adding new asset failed: ", err)
// 			}
//
// 			fmt.Println()
// 		case RemoveAssetsWatchList:
// 			reader := bufio.NewReader(os.Stdin)
// 			fmt.Println("Enter an asset(s) (id): ")
// 			assetsString, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Println("Could not read inputted assets")
// 				continue programLoop
// 			}
// 			assetsString = strings.TrimSpace(assetsString)
//
// 			err = removeAssets(assetDb, assetsString)
// 			if err != nil {
// 				log.Fatal("Deleting asset(s) failed: ", err)
// 			}
//
// 			fmt.Println()
// 		case Exit:
// 			break programLoop
// 		default:
// 			fmt.Println("Please enter a valid input")
// 		}
// 	}
// }
//
// func initAssetDB(path string) (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open database: %w", err)
// 	}
//
// 	// Make sure the connection is valid
// 	if err := db.Ping(); err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("failed to connect to database: %w", err)
// 	}
//
// 	asssetTableStmt := `
// 	CREATE TABLE IF NOT EXISTS assets (
// 		id INTEGER PRIMARY KEY,
// 		symbol TEXT NOT NULL,
// 		price REAL,
// 		market_cap REAL
// 	)`
//
// 	positionsTableStmt := `
// 	CREATE TABLE IF NOT EXISTS positions (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		assetId INTEGER NOT NULL,
// 		amount REAL NOT NULL,
// 		UNIQUE(assetId)
// 	)`
//
// 	tableStmts := []string{asssetTableStmt, positionsTableStmt}
// 	tx, err := db.Begin()
// 	if err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("Failed to start transaction for table creation: %s", err)
// 	}
//
// 	for _, stmt := range tableStmts {
// 		if _, err := tx.Exec(stmt); err != nil {
// 			tx.Rollback()
// 			db.Close()
// 			return nil, fmt.Errorf("Could not create table: %s", err)
// 		}
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("Transaction failed to commit: %s", err)
// 	}
//
// 	return db, nil
// }
//
// func readOperatonInput() (Operation, error) {
// 	for i := Operation(0); i <= Exit; i++ {
// 		fmt.Printf("[%d]  %s\n", i, i.String())
// 	}
//
// 	scanner := bufio.NewScanner(os.Stdin)
// 	if !scanner.Scan() {
// 		return Exit, fmt.Errorf("Invalid input provided.")
// 	}
//
// 	input := strings.TrimSpace(scanner.Text())
// 	operation, err := strconv.Atoi(input)
// 	if err != nil {
// 		return Exit, fmt.Errorf("Input should be numeric")
// 	}
//
// 	return Operation(operation), nil
// }
//
// func addAssets(db *sql.DB, id string) error {
// 	assetIds := strings.Fields(id)
// 	assetData, err := getAssetData(assetIds, "AUD")
// 	if err != nil {
// 		return fmt.Errorf("Could not retrieve asset data: %s", err)
// 	}
//
// 	err = addAssetsToDb(db, assetData)
// 	if err != nil {
// 		log.Fatal("Could not add assets to db: ", err)
// 	}
//
// 	return nil
// }
//
// func getAssetData(ids []string, currency string) ([]AssetDetail, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	q := url.Values{}
// 	q.Add("id", strings.Join(ids, ","))
// 	q.Add("convert", "AUD")
//
// 	ApiKey := os.Getenv("CMC_API_KEY")
// 	req.Header.Set("Accepts", "application/json")
// 	req.Header.Add("X-CMC_PRO_API_KEY", ApiKey)
// 	req.URL.RawQuery = q.Encode()
//
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
//
// 	var response ApiResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	var result []AssetDetail
// 	for _, priceResponse := range response.Data {
// 		assetQuote, ok := priceResponse.Quote[currency]
// 		if !ok {
// 			fmt.Println("Failed to get asset.")
// 		}
//
// 		result = append(result, AssetDetail{Id: priceResponse.Id, Symbol: priceResponse.Symbol, Price: assetQuote.Price, MarketCap: assetQuote.MarketCap})
// 	}
//
// 	return result, nil
// }
//
// func addAssetsToDb(db *sql.DB, assetData []AssetDetail) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("Failed to start db transaction: %s", err)
// 	}
//
// 	stmt, err := tx.Prepare(`INSERT INTO assets (id, symbol, price, market_cap) VALUES (?, ?, ?, ?)`)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("Could not prepare transaction to add assets to db: %s", err)
// 	}
// 	defer stmt.Close()
//
// 	for _, asset := range assetData {
// 		_, err = stmt.Exec(
// 			asset.Id,
// 			asset.Symbol,
// 			asset.Price,
// 			asset.MarketCap,
// 		)
// 		if err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("Could not insert asset data into db: %s", err)
// 		}
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("Could not commit transaction to db: %s", err)
// 	}
//
// 	return nil
// }
//
// func removeAssets(db *sql.DB, assets string) error {
// 	assetIds := strings.Fields(assets)
// 	if len(assetIds) == 0 {
// 		return nil
// 	}
//
// 	idPlaceholders := strings.Repeat("?,", len(assetIds))
// 	idPlaceholders = idPlaceholders[:len(idPlaceholders)-1]
//
// 	stmtArgs := make([]any, len(assetIds))
// 	for i, v := range assetIds {
// 		stmtArgs[i] = v
// 	}
//
// 	stmt := fmt.Sprintf("DELETE FROM assets WHERE id IN (%s)", idPlaceholders)
// 	_, err := db.Exec(stmt, stmtArgs...)
// 	if err != nil {
// 		return fmt.Errorf("Could not delete assets from db: %s", err)
// 	}
//
// 	return nil
// }
//
// func getAllAssets(db *sql.DB) ([]AssetDetail, error) {
// 	rows, err := db.Query(`SELECT * FROM assets`)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query assets: %w", err)
// 	}
// 	defer rows.Close()
//
// 	var assets []AssetDetail
//
// 	for rows.Next() {
// 		var asset AssetDetail
// 		if err := rows.Scan(&asset.Id, &asset.Symbol, &asset.Price, &asset.MarketCap); err != nil {
// 			return nil, fmt.Errorf("failed to scan asset row: %w", err)
// 		}
// 		assets = append(assets, asset)
// 	}
//
// 	// check for row iteration errors
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("row iteration error: %w", err)
// 	}
//
// 	return assets, nil
// }
//
// func getPortfolioValue(db *sql.DB) error {
// 	type position struct {
// 		assetId int
// 		symbol  string
// 		amount  float64
// 		price   float64
// 		mktCap  float64
// 	}
//
// 	positions := make([]position, 0)
// 	assetIds := make([]string, 0)
// 	rows, err := db.Query(`SELECT assetId, amount FROM positions`)
// 	if err != nil {
// 		return fmt.Errorf("failed to query assets: %w", err)
// 	}
// 	defer rows.Close()
//
// 	for rows.Next() {
// 		var portfolioPosition position
// 		if err := rows.Scan(&portfolioPosition.assetId, &portfolioPosition.amount); err != nil {
// 			return fmt.Errorf("Failed to read position in query: %w", err)
// 		}
// 		assetIds = append(assetIds, strconv.Itoa(portfolioPosition.assetId))
// 		positions = append(positions, portfolioPosition)
// 	}
//
// 	assetPriceData, err := getAssetData(assetIds, "AUD")
// 	if err != nil {
// 		return fmt.Errorf("Could not get price data for portfolio positions: %w", err)
// 	}
//
// 	assetIdToPriceData := make(map[int]AssetDetail, 0)
// 	for _, priceInfo := range assetPriceData {
// 		assetIdToPriceData[priceInfo.Id] = priceInfo
// 	}
//
// 	var portfolioValue float64
// 	fmt.Println("Asset  Amount  Price MktCap Size")
// 	for _, pos := range positions {
// 		priceData := assetIdToPriceData[pos.assetId]
// 		pos.mktCap, pos.price, pos.symbol = priceData.MarketCap, priceData.Price, priceData.Symbol
// 		fmt.Printf("%s %f %f %f %f", pos.symbol, pos.amount, pos.price, pos.mktCap, pos.amount*pos.price)
// 		portfolioValue += pos.amount * pos.price
// 	}
// 	fmt.Println("\nTotal value: ", portfolioValue)
//
// 	return nil
// }
//
// func addPosition(db *sql.DB, positionString string) error {
// 	positionArgs := strings.Fields(positionString)
// 	if len(positionArgs) != 2 {
// 		return fmt.Errorf("You did not supply the expected input (AssetID Amount)")
// 	}
// 	id := positionArgs[0]
// 	amount, err := strconv.Atoi(positionArgs[1])
// 	if err != nil {
// 		return fmt.Errorf("Could not add position, invalid amount inputted: %s", err)
// 	}
//
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("could not begin transaction to add position: %s", err)
// 	}
//
// 	_, err = tx.Exec(`INSERT INTO positions (assetId, amount) VALUES (?, ?)`, id, amount)
// 	if err != nil {
// 		return fmt.Errorf("Could not add position into db: %s", err)
// 	}
//
// 	var assetId string
// 	err = tx.QueryRow(`SELECT id FROM assets WHERE id = ?`, id).Scan(&assetId)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		tx.Rollback()
// 		return fmt.Errorf("Failed to query asset db: %s", err)
// 	}
// 	if err != nil {
// 		err = addAssets(db, assetId)
// 		if err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("Failed to add asset to db: %s", err)
// 		}
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		db.Close()
// 		return fmt.Errorf("Failed to commit position to db: %s", err)
// 	}
//
// 	return nil
// }
//
// func removePosition(db *sql.DB, positionString string) error {
// 	positionArgs := strings.Fields(positionString)
// 	if len(positionArgs) != 1 {
// 		return fmt.Errorf("Must supply only one id of the position asset to remove")
// 	}
//
// 	id, err := strconv.Atoi(positionArgs[0])
// 	if err != nil {
// 		return fmt.Errorf("Could not convert id string into integer: %s", err)
// 	}
//
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("could not begin transaction to add position: %s", err)
// 	}
//
// 	_, err = tx.Exec(`DELETE FROM positions WHERE assetId=?`, id)
// 	if err != nil {
// 		return fmt.Errorf("Could not delete position from db: %s", err)
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("Failed to execute delete transaction on positions: %s", err)
// 	}
//
// 	return nil
// }
