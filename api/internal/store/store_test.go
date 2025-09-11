package store

import "testing"

func TestOpen_InMemorySQLite_AutoMigrate(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}

func TestOpenFromURL_FallbackToSQLite(t *testing.T) {
	db, err := OpenFromURL("")
	if err != nil {
		t.Fatalf("open from url failed: %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}
