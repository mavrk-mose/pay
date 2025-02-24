package jobs

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	settlementDir  = "/path/to/settlement/files"
	processedDir   = "/path/to/processed/files"
	reconciliationInterval = 24 * time.Hour
)

type Jobs struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewJobs(db *sqlx.DB) *Jobs {
	return &Jobs{db: db}
}

func StartReconciliationJob(db *sql.DB) {
	ticker := time.NewTicker(reconciliationInterval) // Runs every 24 hours
	go func() {
		for {
			<-ticker.C

			date := time.Now().Format("2006-01-02") // Get today's date
			settlementFile := fmt.Sprintf("%s/settlement_%s.csv", settlementDir, date)

			log.Println("🔄 Running reconciliation job...")
			service.ReconcileTransactions(db, settlementFile, date)
		}
	}()
}

func ReconcileTransactions(db *sql.DB, date string) {
	// Fetch transactions from database
	dbTransactions, err := repository.FetchTransactionsWithChecksum(db, date)
	if err != nil {
		log.Fatalf("Error fetching transactions: %v", err)
	}

	// Read transactions from settlement file
	fileTransactions, err := utils.ReadSettlementFile(settlementDir)
	if err != nil {
		log.Fatalf("Error reading settlement file: %v", err)
	}

	matches := 0
	mismatches := 0
	missingDB := 0
	missingFile := 0

	for id, checksum := range dbTransactions {
		if fileChecksum, found := fileTransactions[id]; found {
			if checksum == fileChecksum {
				matches++
			} else {
				fmt.Printf("⚠️ Mismatch: Transaction %s has different checksums!\n", id)
				mismatches++
			}
		} else {
			fmt.Printf("❌ Missing in Settlement File: %s\n", id)
			missingFile++
		}
	}

	for id := range fileTransactions {
		if _, found := dbTransactions[id]; !found {
			fmt.Printf("❌ Missing in Database: %s\n", id)
			missingDB++
		}
	}

	fmt.Printf("\n✅ Matches: %d\n⚠️ Mismatches: %d\n❌ Missing in Settlement File: %d\n❌ Missing in Database: %d\n",
		matches, mismatches, missingFile, missingDB)
}

func ReadSettlementFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, err
	}

	transactions := make(map[string]string)

	for {
		record, err := reader.Read()
		if err != nil {
			break // EOF
		}

		transactionID := record[0]
		amount, _ := strconv.ParseFloat(record[1], 64)
		currency := record[2]
		timestamp := record[3]

		checksum := utils.GenerateChecksum(transactionID, amount, currency, timestamp)
		transactions[transactionID] = checksum
	}

	return transactions, nil
}