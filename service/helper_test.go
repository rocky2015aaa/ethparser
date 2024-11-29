package service

import (
	"os"
	"testing"
)

func TestRpcCall(t *testing.T) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}
	os.Setenv(EnvRpcUrl, DefaultRpcUrl)
	response, err := rpcCall(request)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response["result"] == "" {
		t.Errorf("result is empty")
	}
}
