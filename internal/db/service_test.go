package db

import (
	"os"
	"path"
	"testing"

	"github.com/broswen/kvs/internal/item"
)

func TestSqlite(t *testing.T) {
	os.Setenv("SQLITE_DB", path.Join(os.TempDir(), "kvs-test-db.db"))
	db, err := New()
	if err != nil {
		t.Fatalf("new db: %v\n", err)
	}
	err = db.Set(item.Item{Key: "test", Value: "test"})
	if err != nil {
		t.Fatalf("set item: %v\n", err)
	}
	item, err := db.Get("test")
	if err != nil {
		t.Fatalf("set item: %v\n", err)
	}
	if item.Value != "test" {
		t.Fatalf("wanted test but got %v\n", item.Value)
	}
	err = db.Delete("test")
	if err != nil {
		t.Fatalf("delete item: %v\n", err)
	}
	item, err = db.Get("test")
	if err == nil {
		t.Fatalf("item shouldn't exist")
	}
}
