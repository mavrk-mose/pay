package jobs

// import (
// 	"fmt"
// 	repo "github.com/mavrk-mose/pay/internal/ledger/repository"
// 	models "github.com/mavrk-mose/pay/internal/ledger/models"
// 	"github.com/mavrk-mose/pay/pkg/utils"
// 	"github.com/fsnotify/fsnotify"
// 	"sync"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// )

// const (
// 	settlementDir          = "/path/to/settlement/files"
// 	processedDir           = "/path/to/processed/files"
// 	reconciliationInterval = 24 * time.Hour
// )
// type ReconciliationService struct {
// 	db         *sqlx.DB
// 	s3Path     string
// 	watchDirs []string
// }

// func NewReconciliationService(db *sqlx.DB, s3Path string) *ReconciliationService {
// 	return &ReconciliationService{
// 		db:     db,
// 		s3Path: s3Path,
// 		watchDirs: viper.GetStringSlice("settlement_directories"),
// 	}
// }

// // gets triggered the moment a file is dumped in the directory : a file watcher
// func (r *ReconciliationService) StartReconciliationJob() {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer watcher.Close()
// 	done := make(chan bool)
// 	go func() {
// 		for {
// 			select {
// 			case event, ok := <-watcher.Events:
// 				if !ok {
// 					return
// 				}
// 				if event.Op&fsnotify.Create == fsnotify.Create {
// 					fmt.Println("New file detected:", event.Name)
// 					go r.ProcessSettlementFile(event.Name) 
// 				}
// 			case err, ok := <-watcher.Errors:
// 				if !ok {
// 					return
// 				}
// 				fmt.Println("Watcher error:", err)
// 			}
// 		}
// 	}()

// 	for _, dir := range r.watchDirs {
// 		if err := watcher.Add(dir); err != nil {
// 			panic(err)
// 		}
// 	}
// 	<-done
// }

// func (r *ReconciliationService) ProcessSettlementFile(filePath string) {
// 	paymentProvider := extractPaymentProvider(filePath)
	
// 	settlementTransactions, err := utils.ReadCSV(filePath, parseSettlementRow, true)
// 	if err != nil {
// 		fmt.Println("Error reading settlement file:", err)
// 		return
// 	}

// 	dbTransactions, err := repo.FetchTransactionsWithChecksum(r.db, time.Now().Format("2006-01-02"),paymentProvider)

// 	var wg sync.WaitGroup
// 	report := [][]string{{"Transaction ID", "Status"}}

// 	jobs := make(chan models.Transaction, len(settlementTransactions))
// 	results := make(chan []string, len(settlementTransactions))

// 	worker := func() {
// 		for txn := range jobs {
// 			status := "Matched"
// 			if dbChecksum, exists := dbTransactions[txn.ID.String()]; exists {
// 				if dbChecksum != txn.Checksum {
// 					status = "Mismatch"
// 				}
// 			} else {
// 				status = "Missing in DB"
// 			}
// 			results <- []string{txn.ID.String(), status}
// 		}
// 		wg.Done()
// 	}

// 	workerCount := len(watchDirs) // Number of concurrent workers
// 	wg.Add(workerCount)
// 	for i := 0; i < workerCount; i++ {
// 		go worker()
// 	}

// 	for _, txn := range settlementTransactions {
// 		jobs <- txn
// 	}
// 	close(jobs)
// 	wg.Wait()
// 	close(results)

// 	for res := range results {
// 		report = append(report, res)
// 	}

// 	timestamp := time.Now().Format("20060102_150405")
// 	reportFileName := fmt.Sprintf("reconciliation_report_%s.csv", timestamp)
// 	utils.WriteCSV(reportFileName, report, nil, func(row []string) []string { return row })
// }

// func parseSettlementRow(row []string) (Transaction, error) {
// 	if len(row) < 4 {
// 		return Transaction{}, fmt.Errorf("invalid row format")
// 	}
// 	return Transaction{ID: row[0], Checksum: row[3]}, nil
// }

// // Extracts the provider from the path format: provider/{{current_date}}/file.csv
// // paypal/ , adyen/ & stripe/ directories
// func extractPaymentProvider(filePath string) string {
// 	dirPath := filepath.Dir(filePath)
// 	return filepath.Base(filepath.Dir(dirPath)) 
// }
