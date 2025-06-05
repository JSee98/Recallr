package dragonfly

import (
	"context"
	"testing"
	"time"
)

func TestDragonflyStore_SetGet(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	val := "hello-world"

	err := store.Set(ctx, key, val, 10*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != val {
		t.Errorf("Expected %s, got %s", val, got)
	}

	// Optional cleanup
	_ = store.Delete(ctx, key)
}

func TestDragonflyStore_GetNonExistent(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	key := "non-existent-key"

	got, err := store.Get(ctx, key)
	if err == nil {
		t.Errorf("Expected error when getting non-existent key, got nil")
	}
	if got != "" {
		t.Errorf("Expected empty string for non-existent key, got %s", got)
	}
}

func TestDragonflyStore_Delete(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	key := "delete-test-key"
	val := "delete-me"

	// Set up test data
	err := store.Set(ctx, key, val, 1*time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists
	_, err = store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed after setting key: %v", err)
	}

	// Delete it
	err = store.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	_, err = store.Get(ctx, key)
	if err == nil {
		t.Errorf("Expected error after deleting key, got nil")
	}
}

func TestDragonflyStore_Expiration(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	key := "expiring-key"
	val := "short-lived"

	// Set with very short TTL
	err := store.Set(ctx, key, val, 1*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists immediately
	got, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed immediately after set: %v", err)
	}
	if got != val {
		t.Errorf("Expected %s, got %s", val, got)
	}

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Verify it's gone
	_, err = store.Get(ctx, key)
	if err == nil {
		t.Errorf("Expected error after key expiration, got nil")
	}
}

func TestDragonflyStore_SetMultiple(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// Set multiple keys
	for k, v := range testData {
		err := store.Set(ctx, k, v, 10*time.Second)
		if err != nil {
			t.Fatalf("Set failed for %s: %v", k, err)
		}
	}

	// Verify all keys
	for k, expected := range testData {
		got, err := store.Get(ctx, k)
		if err != nil {
			t.Fatalf("Get failed for %s: %v", k, err)
		}
		if got != expected {
			t.Errorf("For key %s: expected %s, got %s", k, expected, got)
		}
	}

	// Cleanup
	for k := range testData {
		_ = store.Delete(ctx, k)
	}
}

func TestDragonflyStore_UpdateExisting(t *testing.T) {
	cfg := DragonflyConfig{
		Addr: "localhost:6379",
	}
	store := NewDragonflyStore(&cfg)
	defer store.Close()

	ctx := context.Background()
	key := "update-key"
	val1 := "original"
	val2 := "updated"

	// Set initial value
	err := store.Set(ctx, key, val1, 10*time.Second)
	if err != nil {
		t.Fatalf("Initial set failed: %v", err)
	}

	// Update value
	err = store.Set(ctx, key, val2, 10*time.Second)
	if err != nil {
		t.Fatalf("Update set failed: %v", err)
	}

	// Verify updated value
	got, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != val2 {
		t.Errorf("Expected updated value %s, got %s", val2, got)
	}

	// Cleanup
	_ = store.Delete(ctx, key)
}
