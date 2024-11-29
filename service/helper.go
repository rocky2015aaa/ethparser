package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rocky2015aaa/ethparser/models"
	"github.com/rocky2015aaa/ethparser/repository"
	"github.com/rocky2015aaa/ethparser/util"
)

const (
	EnvRpcUrl     = "RPC_URL"
	DefaultRpcUrl = "https://ethereum-sepolia-rpc.publicnode.com"
)

// Helper function to make JSON-RPC calls.
func rpcCall(request map[string]interface{}) (map[string]interface{}, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(os.Getenv(EnvRpcUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// getCurrentBlockNumber fetches the current block number from Ethereum.
func getCurrentBlockNumber() (int, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}
	response, err := rpcCall(request)
	if err != nil {
		return 0, err
	}
	if _, ok := response["result"]; !ok {
		return 0, fmt.Errorf("invalid RPC response: missing 'result'")
	}
	if _, ok := response["result"].(string); !ok {
		return 0, fmt.Errorf("invalid 'result' type")
	}
	blockNum, err := strconv.ParseInt(response["result"].(string), 0, 64)
	return int(blockNum), err
}

// ParseBlock parses an Ethereum block data and updates the store.
func parseBlock(blockNum int, store *repository.InMemoryStore) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNum), true},
		"id":      1,
	}
	response, err := rpcCall(request)
	if err != nil {
		log.Println("Error fetching block:", err)
		return
	}
	if _, ok := response["result"]; !ok {
		log.Printf("invalid RPC response: missing 'result'")
		return
	}
	block, ok := response["result"].(map[string]interface{})
	if !ok || block == nil {
		log.Printf("Invalid block data for block number %d", blockNum)
		return
	}
	var blockTimestamp string
	if timestamp, ok := block["timestamp"].(string); ok {
		blockTimestamp = util.HexToTime(timestamp).String()
	} else {
		log.Printf("Missing or invalid timestamp for block %d", blockNum)
		return
	}
	transactions, ok := block["transactions"].([]interface{})
	if !ok {
		log.Printf("No transactions found in block %d", blockNum)
		return
	}

	for _, txData := range transactions {
		tx, ok := txData.(map[string]interface{})
		if !ok {
			continue
		}
		txObj := models.Transaction{
			BlockNum:  blockNum,
			Timestamp: blockTimestamp,
		}
		if hash, ok := tx["hash"].(string); ok {
			txObj.Hash = hash
		}
		if from, ok := tx["from"].(string); ok {
			txObj.From = from
		}
		if to, ok := tx["to"].(string); ok {
			txObj.To = to
		}
		if value, ok := tx["value"].(string); ok {
			txObj.Value = util.HexWeiToEther(value)
		}
		if blockHash, ok := tx["blockHash"].(string); ok {
			txObj.BlockHash = blockHash
		}
		store.Mu.Lock()
		if store.Subscriptions[txObj.From] {
			store.Transactions[txObj.From] = append(store.Transactions[txObj.From], &txObj)
		}
		if store.Subscriptions[txObj.To] {
			store.Transactions[txObj.To] = append(store.Transactions[txObj.To], &txObj)
		}
		store.Mu.Unlock()
	}
}
