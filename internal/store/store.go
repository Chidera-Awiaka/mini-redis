package store

import (
	"fmt"
	"sync"
	"time"
)

const maxCapacity = 5

type Store struct {
	mu       sync.Mutex
	data     map[string]string
	expireAt map[string]time.Time
	lru      *LRU
}

func New() *Store {
	s := &Store{
		data:     make(map[string]string),
		expireAt: make(map[string]time.Time),
		lru:      NewLRU(maxCapacity),
	}

	go s.cleanupExpired()

	return s
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	delete(s.expireAt, key)

	s.lru.Touch(key)
	s.evictIfNeeded()
}

func (s *Store) SetWithTTL(key, value string, ttlSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	s.expireAt[key] = time.Now().Add(time.Duration(ttlSeconds) * time.Second)

	s.lru.Touch(key)
	s.evictIfNeeded()
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if exp, ok := s.expireAt[key]; ok && time.Now().After(exp) {
		delete(s.data, key)
		delete(s.expireAt, key)
		s.lru.Remove(key)
		return "", false
	}

	val, ok := s.data[key]
	if !ok {
		return "", false
	}

	s.lru.Touch(key)
	return val, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	delete(s.expireAt, key)
	s.lru.Remove(key)
}

func (s *Store) Stats() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalKeys := len(s.data)
	ttlKeys := len(s.expireAt)
	lruSize := s.lru.Size()

	// Rough memory estimate
	var approxBytes int
	for k, v := range s.data {
		approxBytes += len(k) + len(v)
	}

	return fmt.Sprintf(
		"keys=%d ttl_keys=%d lru_size=%d max_capacity=%d approx_bytes=%d",
		totalKeys,
		ttlKeys,
		lruSize,
		maxCapacity,
		approxBytes,
	)
}

func (s *Store) evictIfNeeded() {
	for s.lru.Size() > maxCapacity {
		oldest := s.lru.RemoveOldest()
		if oldest == "" {
			return
		}
		delete(s.data, oldest)
		delete(s.expireAt, oldest)
	}
}

func (s *Store) cleanupExpired() {
	for {
		time.Sleep(1 * time.Second)

		s.mu.Lock()
		now := time.Now()

		for key, exp := range s.expireAt {
			if now.After(exp) {
				delete(s.data, key)
				delete(s.expireAt, key)
				s.lru.Remove(key)
			}
		}

		s.mu.Unlock()
	}
}
