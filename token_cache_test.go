package requesterV2

import "testing"

func TestSetGetToken(t *testing.T) {
	cache := NewTokenCache()
	cache.Set("test", "test")
	val, ok := cache.Get("test")

	if !ok || val != "test" {
		t.Errorf("Invalid value or value not found")
	}
}

func TestGetInvalidToken(t *testing.T) {
	cache := NewTokenCache()
	cache.Set("test", "test")
	val, ok := cache.Get("test2")

	if ok || val != "" {
		t.Errorf("Invalid value or value not found")
	}
}

func TestDeleteToken(t *testing.T) {
	cache := NewTokenCache()
	cache.Set("test", "test")
	val, ok := cache.Get("test")

	if !ok {
		t.Error("Error to get value")
	}

	cache.Delete("test")

	val2, ok2 := cache.Get("test")

	if val2 == val || ok2 {
		t.Error("Value was not deleted")
	}
}
