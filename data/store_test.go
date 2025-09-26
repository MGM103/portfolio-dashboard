package data_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	data "github.com/mgm103/portfolio-dashboard/data"
)

func setupTestStore(t *testing.T) *data.Store {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	store := &data.Store{Conn: db}

	// Call createTables to set up schema
	if err := store.CreateTables(); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	return store
}

func TestSaveAndGetWatchlist(t *testing.T) {
	store := setupTestStore(t)

	assets := []data.Asset{
		{ID: "1", Ticker: "BTC"},
		{ID: "2", Ticker: "ETH"},
	}

	if err := store.SaveToWatchlist(assets); err != nil {
		t.Fatalf("SaveToWatchlist failed: %v", err)
	}

	watchlist, err := store.GetWatchlist()
	if err != nil {
		t.Fatalf("GetWatchlist failed: %v", err)
	}

	if len(watchlist) != 2 {
		t.Errorf("expected 2 assets, got %d", len(watchlist))
	}
}

func TestSaveToPositionsAndGetPositions(t *testing.T) {
	store := setupTestStore(t)

	asset := data.Asset{
		ID:     "1",
		Ticker: "BTC",
		Amount: 3,
	}

	err := store.SaveToPositions(asset)
	if err != nil {
		t.Fatalf("SaveToPositions failed: %v", err)
	}

	positions, err := store.GetPositions()
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}

	if len(positions) != 1 {
		t.Errorf("expected 1 position, got %d", len(positions))
	}

	if positions[0].ID != asset.ID || positions[0].Ticker != asset.Ticker || positions[0].Amount != asset.Amount {
		t.Errorf("unexpected asset data: %+v", positions[0])
	}
}

func TestSaveToWatchlist_Upsert(t *testing.T) {
	store := setupTestStore(t)

	initial := data.Asset{ID: "1", Ticker: "BTC"}
	updated := data.Asset{ID: "1", Ticker: "XBTC"}

	err := store.SaveToWatchlist([]data.Asset{initial})
	if err != nil {
		t.Fatalf("initial insert failed: %v", err)
	}

	err = store.SaveToWatchlist([]data.Asset{updated})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	watchlist, err := store.GetWatchlist()
	if err != nil {
		t.Fatalf("GetWatchlist failed: %v", err)
	}

	if len(watchlist) != 1 {
		t.Fatalf("expected 1 item after upsert, got %d", len(watchlist))
	}
	if watchlist[0].Ticker != "XBTC" {
		t.Errorf("expected updated ticker 'XBTC', got '%s'", watchlist[0].Ticker)
	}
}

func TestSaveToPositions_Upsert(t *testing.T) {
	store := setupTestStore(t)

	initial := data.Asset{ID: "2", Ticker: "ETH", Amount: 1}
	updated := data.Asset{ID: "2", Ticker: "ETH", Amount: 5}

	err := store.SaveToPositions(initial)
	if err != nil {
		t.Fatalf("initial insert failed: %v", err)
	}

	err = store.SaveToPositions(updated)
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	positions, err := store.GetPositions()
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}

	if len(positions) != 1 {
		t.Fatalf("expected 1 position after upsert, got %d", len(positions))
	}
	if positions[0].Amount != 5 {
		t.Errorf("expected amount 5, got %f", positions[0].Amount)
	}
}

func TestDeleteFromWatchlistByIDs(t *testing.T) {
	store := setupTestStore(t)

	assets := []data.Asset{
		{ID: "1", Ticker: "BTC"},
		{ID: "2", Ticker: "ETH"},
		{ID: "3", Ticker: "XRP"},
	}
	if err := store.SaveToWatchlist(assets); err != nil {
		t.Fatalf("SaveToWatchlist failed: %v", err)
	}

	// Delete 2 of the 3 assets
	toDelete := []string{"1", "3"}
	if err := store.RemoveFromWatchlist(toDelete); err != nil {
		t.Fatalf("DeleteFromWatchlistByIDs failed: %v", err)
	}

	// Fetch remaining items
	remaining, err := store.GetWatchlist()
	if err != nil {
		t.Fatalf("GetWatchlist failed: %v", err)
	}

	// Expect only the ETH asset to remain
	if len(remaining) != 1 {
		t.Fatalf("expected 1 asset after delete, got %d", len(remaining))
	}
	if remaining[0].ID != "2" || remaining[0].Ticker != "ETH" {
		t.Errorf("unexpected remaining asset: %+v", remaining[0])
	}
}
