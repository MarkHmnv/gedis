package main

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()
	if cache == nil {
		t.Fail()
	}
}

func TestSetKey(t *testing.T) {
	cache := NewCache()
	args := []string{"SET", "TestKey", "TestValue"}
	res, err := setKey(args, cache)
	if res != "OK" || err != nil {
		t.Fail()
	}
}

func TestGetKey(t *testing.T) {
	cache := NewCache()
	setArgs := []string{"SET", "TestKey", "TestValue"}
	setKey(setArgs, cache)
	getArgs := []string{"GET", "TestKey"}
	value := getKey(getArgs, cache)
	if value != "TestValue" {
		t.Error("Expected TestValue, got ", value)
		t.Fail()
	}
}

func TestKeyExpiry(t *testing.T) {
	cache := NewCache()
	args := []string{"SET", "TestKey", "TestValue", "PX", "500"}
	setKey(args, cache)
	time.Sleep(time.Duration(600) * time.Millisecond)
	getArgs := []string{"GET", "TestKey"}
	value := getKey(getArgs, cache)
	if value != "(nil)" {
		t.Fail()
	}
}
