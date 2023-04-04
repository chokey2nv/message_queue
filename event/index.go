package event

import (
	"context"
	"errors"
)

// type Transaction struct {
// 	ClientID      int64  `json:"clientId"`
// 	WalletAddress string `json:"walletAddress"`
// 	CurrencyType  string `json:"currencyType"`
// 	TransactionID string `json:"transactionId"`
// }
type Transaction struct {
	ID            uint64 `gorm:"primaryKey" json:"id"`
	ClientID      uint64 `gorm:"not null" json:"client_id"`
	FromAddress   string `gorm:"not null" json:"from_address"`
	ToAddress     string `gorm:"not null" json:"to_address"`
	CurrencyType  string `gorm:"not null" json:"currency_type"`
	TransactionID string `gorm:"not null" json:"transaction_id"`
	Amount        uint64 `gorm:"not null" json:"amount"`
	Timestamp     uint64 `gorm:"not null" json:"timestamp"`
}

type EventBusClient interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Close() error
}

func NewEventBusClient() (EventBusClient, error) {
	//TODO: Initialize and return the appropriate event bus client based on the chosen technology

	return nil, errors.New("event bus client not implemented")
}
