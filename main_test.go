package main

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()
	if cache == nil || cache.Data == nil || cache.Expiry == nil {
		t.Fatal("NewCache function is not initializing cache correctly.")
	}
}

func TestSetKey(t *testing.T) {
	cache := NewCache()
	args := []string{"SET", "key", "value"}

	result, err := setKey(args, cache)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if result != "OK" {
		t.Fatalf("Expected OK, got %s", result)
	}

	value, exists := cache.Data["key"]

	if !exists {
		t.Fatal("Key was not set in cache")
	}

	if value != "value" {
		t.Fatalf("Expected value 'value', got '%s'", value)
	}
}

func TestGetKey(t *testing.T) {
	cache := NewCache()
	cache.Data["key"] = "value"

	args := []string{"GET", "key"}

	result := getKey(args, cache)

	if result != "value" {
		t.Fatalf("Expected 'value', got '%s'", result)
	}

	cache.Expiry["key"] = time.Now().Add(-1 * time.Second)

	result = getKey(args, cache)

	if result != "(nil)" {
		t.Fatalf("Expected '(nil)', got '%s'", result)
	}
}

func TestKeyExpiry(t *testing.T) {
	cache := NewCache()

	args := []string{"SET", "key", "value", "EX", "1"}
	_, err := setKey(args, cache)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	time.Sleep(1 * time.Second)

	result := getKey([]string{"GET", "key"}, cache)

	if result != "(nil)" {
		t.Fatalf("Expected '(nil)', got '%s'", result)
	}
}
