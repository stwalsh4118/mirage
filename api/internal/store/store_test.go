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
