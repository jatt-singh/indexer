package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pg_password"
	dbname   = "pd_data"
)

func ConnectDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
	return db
}

const schema = `
CREATE TABLE IF NOT EXISTS blocks (
    block_number BIGINT PRIMARY KEY,
    block_hash TEXT NOT NULL,
    transaction_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_hash TEXT NOT NULL UNIQUE,
    block_number BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (block_number) REFERENCES blocks(block_number) ON DELETE CASCADE
);
`

func InitSchema(db *sql.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Error initializing schema: %v", err)
	}
	log.Println("Database schema initialized successfully.")
}

func InsertBlock(db *sql.DB, blockNumber int64, blockHash string, transactionCount int) error {
	query := `
		INSERT INTO blocks (block_number, block_hash, transaction_count)
		VALUES ($1, $2, $3)
		ON CONFLICT (block_number) DO NOTHING;
	`
	_, err := db.Exec(query, blockNumber, blockHash, transactionCount)
	if err != nil {
		return fmt.Errorf("failed to insert block: %v", err)
	}
	return nil
}

func InsertTransaction(db *sql.DB, transactionHash string, blockNumber int64) error {
	query := `
		INSERT INTO transactions (transaction_hash, block_number)
		VALUES ($1, $2)
		ON CONFLICT (transaction_hash) DO NOTHING;
	`
	_, err := db.Exec(query, transactionHash, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}
	return nil
}
