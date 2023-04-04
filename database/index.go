package database

import (
	"fmt"

	"github.com/chokey2nv/obiex.finance/event"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbClient struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string
}

func NewDBClient(host string, port int, user string, pass string, dbName string) *DbClient {
	return &DbClient{
		Host:   host,
		Port:   port,
		User:   user,
		Pass:   pass,
		DBName: dbName,
	}
}
func (dbClient *DbClient) MySQLURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbClient.User,
		dbClient.Pass,
		dbClient.Host,
		dbClient.Port,
		dbClient.DBName,
	)
}
func (dbClient *DbClient) ConnectToDatabase() (*gorm.DB, error) {
	// dsn := fmt.Sprintf("username:password@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)
	db, err := gorm.Open(mysql.Open(dbClient.MySQLURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// AutoMigrate will automatically create a table named "transactions" in the database, and use the
	// Transaction struct to define the table schema.
	if err = db.AutoMigrate(&event.Transaction{}); err != nil {
		return nil, err
	}
	return db, nil
}
func CloseDB(db *gorm.DB) {
	sqlDb, _ := db.DB()
	sqlDb.Close()
}

func (dbClient *DbClient) StoreTransaction(tx *event.Transaction) error {
	// Connect to the database
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return err
	}
	defer CloseDB(db)

	//check if transactionId exist for specific wallet and client
	var count int64
	db.Model(&event.Transaction{}).Where("transaction_id = ?",
		tx.TransactionID).Count(&count)

	if count > 0 {
		return fmt.Errorf("transaction already exists")
	}

	// Save the transaction to the database
	result := db.Create(tx)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (dbClient *DbClient) GetLastTxTimestamp(clientID uint64, clientAddress string) (int64, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return 0, err
	}
	defer CloseDB(db)
	var lastTimestamp int64
	err = db.Model(&event.Transaction{}).
		Select("Max(timestamp)").
		Where("client_id = ? AND client_address = ?", clientID, clientAddress).
		Scan(&lastTimestamp).Error
	if err != nil {
		return 0, err
	}
	return lastTimestamp, nil
}
func (dbClient *DbClient) GetAllTransactions(page, pageSize int) ([]event.Transaction, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	var transactions []event.Transaction
	offset := (page - 1) * pageSize
	err = db.Order("timestamp desc").Limit(pageSize).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Retrieve the total count of transactions
func (dbClient *DbClient) GetTransactionCount() (int64, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return 0, err
	}
	defer CloseDB(db)
	var count int64
	err = db.Model(&event.Transaction{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (dbClient *DbClient) GetAllTransactionsByAddress(walletAddress string, page, pageSize int) ([]event.Transaction, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	var transactions []event.Transaction
	offset := (page - 1) * pageSize
	err = db.Order("timestamp desc").Limit(pageSize).Offset(offset).Where(
		"from_address = ? OR to_address = ?",
		walletAddress, walletAddress).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func (dbClient *DbClient) GetTransactionCountByWallet(walletAddress string) (int64, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return 0, err
	}
	defer CloseDB(db)
	var count int64

	err = db.Model(&event.Transaction{}).Where(
		"from_address = ? OR to_address = ?",
		walletAddress, walletAddress).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dbClient *DbClient) GetAllTransactionsByClientId(clientId uint64, page, pageSize int) ([]event.Transaction, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	var transactions []event.Transaction
	offset := (page - 1) * pageSize
	err = db.Order("timestamp desc").Limit(pageSize).Offset(offset).Where(
		"client_id = ?",
		clientId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func (dbClient *DbClient) GetTransactionCountByClientId(clientId uint64) (int64, error) {
	db, err := dbClient.ConnectToDatabase()
	if err != nil {
		return 0, err
	}
	defer CloseDB(db)
	var count int64

	err = db.Model(&event.Transaction{}).Where(
		"from_address = ?",
		clientId).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
