package repository

import (
	"sync"

	"github.com/rocky2015aaa/ethparser/models"
)

// InMemoryStore struct to hold data in memory.
type InMemoryStore struct {
	Subscriptions map[string]bool
	Transactions  map[string][]*models.Transaction
	CurrentBlock  int
	Mu            sync.RWMutex
}

// NewInMemoryStore initializes a new in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Subscriptions: make(map[string]bool),
		Transactions:  make(map[string][]*models.Transaction),
		CurrentBlock:  0,
	}
}

// GetCurrentBlock returns the last parsed block.
func (s *InMemoryStore) GetCurrentBlock() int {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.CurrentBlock
}

// SetCurrentBlock sets the current block.
func (s *InMemoryStore) SetCurrentBlock(blockNum int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.CurrentBlock = blockNum
}

// Subscribe adds an address to the subscription list.
func (s *InMemoryStore) Subscribe(address string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if _, exists := s.Subscriptions[address]; exists {
		return false
	}
	s.Subscriptions[address] = true
	return true
}

// GetTransactions retrieves transactions for an address.
func (s *InMemoryStore) GetTransactions(address string) []*models.Transaction {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Transactions[address]
}
