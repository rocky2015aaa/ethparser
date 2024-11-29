package util

import (
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IsHexAddress checks if a given string is a valid Ethereum address.
func IsHexAddress(address string) bool {
	// An Ethereum address must start with "0x" and be 42 characters long
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {
		return false
	}
	// Ensure the address contains only valid hexadecimal characters after "0x"
	matched, _ := regexp.MatchString("^[0-9a-fA-F]{40}$", address[2:])
	return matched
}

// HexToTime converts hex string time to time
func HexToTime(hexTimestamp string) time.Time {
	timestamp, _ := strconv.ParseInt(hexTimestamp, 0, 64)
	return time.Unix(timestamp, 0)
}

// HexWeiToEther converts hex string value to ETH amount
func HexWeiToEther(hexWei string) string {
	wei := new(big.Int)
	_, ok := wei.SetString(hexWei[2:], 16) // Remove "0x" prefix and parse as hex
	if !ok {
		log.Printf("Invalid hex value: %s", hexWei)
		return "0"
	}
	ether := new(big.Float).SetInt(wei)
	ether = ether.Quo(ether, big.NewFloat(1e18)) // Convert wei to Ether
	return ether.Text('g', -1)                   // Use 'g' formatting to avoid trailing zeros
}
