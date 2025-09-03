package data

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Recommended to place import on its own line
)

type Asset struct {
	ID     string
	Ticker string
	Amount uint64
}

type Store struct {
	Conn *sql.DB
}

func (s *Store) Init() error {
	dbPath := "./data/asset.db"
	var err error
	s.Conn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database at %s: %w", dbPath, err)
	}

	if err := s.CreateTables(); err != nil {
		s.Conn.Close()
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (s *Store) CreateTables() error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	createTablesStmt := `
	CREATE TABLE IF NOT EXISTS watchlist (
		id TEXT PRIMARY KEY,
		ticker TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS positions (
		id TEXT PRIMARY KEY,
		ticker TEXT NOT NULL UNIQUE,
		amount REAL NOT NULL
	);`

	_, err = tx.Exec(createTablesStmt)
	if err != nil {
		return fmt.Errorf("failed to execute create tables statement: %w", err)
	}

	return tx.Commit()
}

func (s *Store) GetWatchlist() ([]Asset, error) {
	rows, err := s.Conn.Query(`SELECT id, ticker FROM watchlist`)
	if err != nil {
		return nil, fmt.Errorf("Watchlist query failed: %w", err)
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		rows.Scan(&asset.ID, &asset.Ticker)
		assets = append(assets, asset)
	}

	return assets, nil
}

func (s *Store) GetPositions() ([]Asset, error) {
	rows, err := s.Conn.Query(`SELECT * FROM positions`)
	if err != nil {
		return nil, fmt.Errorf("Positions query failed: %w", err)
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		rows.Scan(&asset.ID, &asset.Ticker, &asset.Amount)
		assets = append(assets, asset)
	}

	return assets, nil
}

func (s *Store) SaveToWatchlist(assets []Asset) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	upsertQuery := ` INSERT INTO watchlist (id, ticker)
	VALUES(?, ?)
	ON CONFLICT(id) DO UPDATE
	SET ticker=excluded.ticker
	`

	stmt, err := tx.Prepare(upsertQuery)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, asset := range assets {
		_, err := stmt.Exec(asset.ID, asset.Ticker)
		if err != nil {
			return fmt.Errorf("Failed to add <%s> to db: %w", asset.ID, err)
		}
	}

	return tx.Commit()
}

func (s *Store) SaveToPositions(asset Asset) error {
	upsertQuery := ` INSERT INTO positions (id, ticker, amount)
	VALUES(?, ?, ?)
	ON CONFLICT(id) DO UPDATE
	SET ticker=excluded.ticker, amount=excluded.amount
	`

	_, err := s.Conn.Exec(upsertQuery, asset.ID, asset.Ticker, asset.Amount)
	if err != nil {
		return fmt.Errorf("Failed to add or update position: %w", err)
	}

	return nil
}

func (s *Store) RemoveFromWatchlist(assetIds []string) error {
	if len(assetIds) == 0 {
		return nil
	}

	placeholders := make([]string, len(assetIds))
	for i := range len(placeholders) {
		placeholders[i] = "?"
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM watchlist WHERE id IN (%s)`, strings.Join(placeholders, ","))
	args := make([]any, len(assetIds))
	for i, id := range assetIds {
		args[i] = id
	}

	_, err := s.Conn.Exec(deleteQuery, args...)
	if err != nil {
		return fmt.Errorf("Failed to execute delete query in db: %w", err)
	}

	return nil
}
