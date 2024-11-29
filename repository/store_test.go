package repository

import (
	"fmt"
	"sync"
	"testing"
)

func TestSubscribe(t *testing.T) {
	store := NewInMemoryStore()
	success := store.Subscribe("0x12345")
	if !success {
		t.Errorf("Expected subscription to succeed")
	}
	success = store.Subscribe("0x12345")
	if success {
		t.Errorf("Expected subscription to fail for duplicate address")
	}
}

func TestConcurrentSubscribe(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			store.Subscribe(addr)
		}(fmt.Sprintf("0x12345%d", i))
	}
	wg.Wait()
	if len(store.Subscriptions) != 100 {
		t.Errorf("Expected 100 subscriptions, got %d", len(store.Subscriptions))
	}
}
