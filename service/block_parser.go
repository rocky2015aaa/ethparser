package service

import (
	"log"
	"sync"
	"time"

	"github.com/rocky2015aaa/ethparser/models"
	"github.com/rocky2015aaa/ethparser/repository"
)

// Store interface for accessing transactions and blocks.
type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []*models.Transaction
}

type BlockService struct {
	store *repository.InMemoryStore
}

func NewBlockService() *BlockService {
	return &BlockService{
		store: repository.NewInMemoryStore(),
	}
}

// GetCurrentBlockNumber fetches the current block number from Ethereum.
func (s *BlockService) GetCurrentBlock() int {
	return s.store.GetCurrentBlock()
}

func (s *BlockService) Subscribe(address string) bool {
	return s.store.Subscribe(address)
}

func (s *BlockService) GetTransactions(address string) []*models.Transaction {
	return s.store.GetTransactions(address)
}

// LoadBlocks fetches and parses new blocks periodically.
func (s *BlockService) LoadBlocks(pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	for range ticker.C {
		currentBlock, err := getCurrentBlockNumber()
		if err != nil {
			log.Println("Error fetching current block:", err)
			continue
		}
		lastBlock := s.store.GetCurrentBlock()
		if lastBlock == 0 { // Load block data from current block, not load all the block data
			s.store.Mu.Lock()
			s.store.CurrentBlock = currentBlock
			s.store.Mu.Unlock()
		} else if currentBlock > lastBlock {
			go func() {
				var wg sync.WaitGroup
				for i := lastBlock + 1; i <= currentBlock; i++ {
					wg.Add(1)
					go func(blockNum int) {
						defer wg.Done()
						parseBlock(blockNum, s.store)
					}(i)
				}
				wg.Wait()
			}()
			s.store.Mu.Lock()
			s.store.CurrentBlock = currentBlock
			s.store.Mu.Unlock()
		}
	}
}
