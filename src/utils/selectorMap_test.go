package utils

import (
	"testing"
)

func TestLoadSelectorMap_ExactMatch(t *testing.T) {
	sm, err := LoadSelectorMap("1.2.30.1135.g1c8a648e")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sm.Version != "1.2.30" {
		t.Errorf("Expected version 1.2.30, got %s", sm.Version)
	}
}

func TestLoadSelectorMap_FallbackWhenMissing(t *testing.T) {
	sm, err := LoadSelectorMap("1.1.1.99")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sm.Version != "fallback" {
		t.Errorf("Expected fallback version, got %s", sm.Version)
	}
}

func TestLoadSelectorMap_EmptyString(t *testing.T) {
	sm, err := LoadSelectorMap("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sm.Version != "fallback" {
		t.Errorf("Expected fallback version, got %s", sm.Version)
	}
}
