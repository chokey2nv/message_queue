package control

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chokey2nv/obiex.finance/config"
	"github.com/chokey2nv/obiex.finance/crypto"
	"github.com/chokey2nv/obiex.finance/event"
	"github.com/streadway/amqp"
)

type TransactionUpdate struct {
	ClientID      uint64
	WalletAddress string
	CurrencyType  string
}
type AppController struct {
	*config.Config
}

func NewController(cfg *config.Config) *AppController {
	return &AppController{cfg}
}
func (ctler *AppController) HandleMessage(delivery amqp.Delivery) error {
	// Unmarshal the message body into a TransactionUpdate struct
	var update TransactionUpdate
	err := json.Unmarshal(delivery.Body, &update)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message body: %v", err)
	}

	//Get last timestamp queried
	lastTimestamp, err := ctler.DBClient.GetLastTxTimestamp(update.ClientID, update.WalletAddress)
	if err != nil {
		return err
	}
	// Perform the necessary query to the crypto API here using update.WalletAddress and update.CurrencyType
	txs, err := crypto.QueryCryptoAPI(update.WalletAddress, update.CurrencyType, lastTimestamp)
	if err != nil {
		return err
	}

	// Initialize event bus client
	eventBusClient, err := event.NewEventBusClient()
	if err != nil {
		return err
	}
	for _, tx := range txs {
		// Publish NewTransaction event to event bus
		newTransactionEvent := event.Transaction{
			ClientID:      update.ClientID,
			FromAddress:   tx.From,
			ToAddress:     tx.To,
			CurrencyType:  update.CurrencyType,
			TransactionID: tx.Hash,
			Amount:        tx.Value.Uint64(),
			Timestamp:     uint64(tx.Timestamp),
		}
		// Convert the message payload to byte buffer
		msg, err := json.Marshal(newTransactionEvent)
		if err != nil {
			return err
		}
		//Store in database
		err = ctler.DBClient.StoreTransaction(&newTransactionEvent)
		if err != nil {
			return err
		}
		//publish event
		err = eventBusClient.Publish(context.Background(), "new_transaction", msg)
		if err != nil {
			return err
		}
	}
	return nil
}
