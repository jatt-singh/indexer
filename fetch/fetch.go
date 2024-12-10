package fetch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
)


// RPCResponse structure for parsing the node response
type RPCResponse struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BlockID struct {
			Hash string `json:"hash"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Height string `json:"height"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	} `json:"result"`
}

// Function to fetch block data from the RPC node
func FetchBlockData(rpcURL string) (*RPCResponse, error) {
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from RPC: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RPC response: %v", err)
	}

	var rpcResp RPCResponse
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RPC response: %v", err)
	}

	return &rpcResp, nil
}

// Function to save block and transactions data into the database
func SaveData(db *sql.DB, blockNumber int64, blockHash string, transactions []string) {
	// Insert block data
	_, err := db.Exec(
		"INSERT INTO blocks (block_number, block_hash, transaction_count) VALUES ($1, $2, $3) ON CONFLICT (block_number) DO NOTHING",
		blockNumber, blockHash, len(transactions),
	)
	if err != nil {
		log.Printf("Error inserting block: %v", err)
	}

	// Insert transactions data
	for _, tx := range transactions {
		_, err := db.Exec(
			"INSERT INTO transactions (transaction_hash, block_number) VALUES ($1, $2) ON CONFLICT (transaction_hash) DO NOTHING",
			tx, blockNumber,
		)
		if err != nil {
			log.Printf("Error inserting transaction: %v", err)
		}
	}
}
