package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rocky2015aaa/ethparser/models"
	"github.com/rocky2015aaa/ethparser/service"
	"github.com/rocky2015aaa/ethparser/util"
)

type TransactionHandler struct {
	service service.Parser
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(service service.Parser) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleSubscribe godoc
//
// @Description	Allows a user to subscribe to an Ethereum address to track transactions
// @Tags		subscription
// @Param		address body string true "Address to subscribe"
// @Success 200 {string} string "Subscribed to address: {address}"
// @Failure 400 {string} string "Bad Request Error Message"
// @Failure 405 {string} string "Method Not Allowed Error Message"
// @Router /subscribe [post]
func (h *TransactionHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		respondWithError(w, "Address is required", http.StatusBadRequest)
		return
	}
	if !util.IsHexAddress(address) {
		respondWithError(w, "Address is invalid", http.StatusBadRequest)
		return
	}
	if h.service.Subscribe(strings.ToLower(address)) {
		fmt.Fprintf(w, "Subscribed to address: %s\n", address)
	} else {
		fmt.Fprintf(w, "Address already subscribed: %s\n", address)
	}
}

// HandleGetTransactions godoc
//
// @Description	Retrieves a list of transactions for an Ethereum address
// @Tags		transactions
// @Param		address query string true "Address to get transactions for"
// @Success 200 {array} models.Transaction "List of transactions"
// @Failure 400 {string} string "Bad Request Error Message"
// @Failure 405 {string} string "Method Not Allowed Error Message"
// @Router /transactions [get]
func (h *TransactionHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	address := r.URL.Query().Get("address")
	if address == "" {
		respondWithError(w, "Address is required", http.StatusBadRequest)
		return
	}
	if !util.IsHexAddress(address) {
		respondWithError(w, "Address is invalid", http.StatusBadRequest)
		return
	}
	transactions := h.service.GetTransactions(strings.ToLower(address))
	w.Header().Set("Content-Type", "application/json")
	if transactions == nil {
		transactions = []*models.Transaction{}
	}
	json.NewEncoder(w).Encode(transactions)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
	log.Printf("HTTP %d: %s", code, message)
}
