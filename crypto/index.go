package crypto

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

type Tx struct {
	Hash             string   `json:"hash"`
	Nonce            uint64   `json:"nonce"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionIndex uint     `json:"transactionIndex"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Value            *big.Int `json:"value"`
	GasPrice         *big.Int `json:"gasPrice"`
	GasLimit         uint64   `json:"gas"`
	Input            []byte   `json:"input"`
	Timestamp        int64    `json:"timestamp"`
}

func QueryCryptoAPI(walletAddress string, currencyType string, fromTimestamp int64) ([]Tx, error) {
	// Assuming there is a third-party API that returns a list of transactions for a given wallet address, and time range of txs
	// The request body for the API call could be similar to the UpdateTransactions command
	// We can use a HTTP client to make the API call
	httpClient := &http.Client{}
	apiUrl := fmt.Sprintf(
		"https://example.com/api/transactions?address=%s&currency=%s&fromtime=%d",
		walletAddress,
		currencyType,
		fromTimestamp,
	)
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error from crypto API: %s", resp.Status)
	}

	var transactions []Tx
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
