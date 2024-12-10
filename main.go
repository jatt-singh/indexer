package main

import (
	"blockchain/indexer/db" // Replace with your actual package path
	"blockchain/indexer/fetch"
	"fmt"
	"time"
	"log"
)


func main() {
	// Connect to the database
	database := db.ConnectDatabase()
	defer database.Close()

	// Initialize schema if not already created
	db.InitSchema(database)

	// Define the RPC URL
	const rpcURL = "https://rpc-dorado.fetch.ai/block"

	// Start a continuous loop to fetch data
	for {
		// Fetch block data from the RPC node
		rpcResp, err := fetch.FetchBlockData(rpcURL)
		if err != nil {
			log.Printf("Error fetching block data: %v", err)
			time.Sleep(10 * time.Second) // Retry after 10 seconds
			continue
		}

		// Extract data
		blockNumber := rpcResp.Result.Block.Header.Height
		blockHash := rpcResp.Result.BlockID.Hash
		transactions := rpcResp.Result.Block.Data.Txs

		// Print fetched data
		fmt.Printf("Fetched Block Number: %s, Hash: %s, Transactions: %d\n", blockNumber, blockHash, len(transactions))

		// Convert block number to integer
		var blockNumInt int64
		fmt.Sscanf(blockNumber, "%d", &blockNumInt)

		// Save data to the database
		fetch.SaveData(database, blockNumInt, blockHash, transactions)

		// Sleep for a while before fetching the next block
		time.Sleep(5 * time.Second) // Adjust sleep time as per the node's update frequency
	}
}
