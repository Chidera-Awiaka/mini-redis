package store

import (
	"testing"
	"time"
)

func TestSetGetDelete(t *testing.T) {
	s := New()

	s.Set("name", "Chidera")

	v, ok := s.Get("name")
	if !ok || v != "Chidera" {
		t.Fatalf("expected (Chidera,true), got (%q,%v)", v, ok)
	}

	s.Delete("name")
	_, ok = s.Get("name")
	if ok {
		t.Fatalf("expected key to be deleted")
	}
}

func TestTTLExpires(t *testing.T) {
	s := New()

	s.SetWithTTL("temp", "x", 1)

	v, ok := s.Get("temp")
	if !ok || v != "x" {
		t.Fatalf("expected temp to exist immediately")
	}

	time.Sleep(2 * time.Second)

	_, ok = s.Get("temp")
	if ok {
		t.Fatalf("expected temp to expire")
	}
}

func TestLRUEviction(t *testing.T) {
	s := New()

	// Fill up to capacity
	s.Set("a", "1")
	s.Set("b", "2")
	s.Set("c", "3")
	s.Set("d", "4")
	s.Set("e", "5")

	// This should evict "a"
	s.Set("f", "6")

	_, ok := s.Get("a")
	if ok {
		t.Fatalf("expected 'a' to be evicted by LRU")
	}

	// "f" must exist
	v, ok := s.Get("f")
	if !ok || v != "6" {
		t.Fatalf("expected 'f' to exist")
	}
}

func TestLRUIsUpdatedOnGet(t *testing.T) {
	s := New()

	s.Set("a", "1")
	s.Set("b", "2")
	s.Set("c", "3")
	s.Set("d", "4")
	s.Set("e", "5")

	// Touch "a" so it becomes most recently used
	_, _ = s.Get("a")

	// Insert new key, should evict "b" now (since a was recently used)
	s.Set("f", "6")

	_, ok := s.Get("b")
	if ok {
		t.Fatalf("expected 'b' to be evicted after touching 'a'")
	}

	_, ok = s.Get("a")
	if !ok {
		t.Fatalf("expected 'a' to remain after being touched")
	}
}
