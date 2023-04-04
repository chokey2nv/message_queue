package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chokey2nv/obiex.finance/api"
	"github.com/chokey2nv/obiex.finance/config"
	"github.com/chokey2nv/obiex.finance/control"
	"github.com/chokey2nv/obiex.finance/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error setting up config file: %v", err)
	}
	controller := control.NewController(config)
	apiClient := api.NewAPIClient(config)

	rabbitUrl := os.Getenv("RABBITMQ_URL")
	if rabbitUrl == "" {
		rabbitUrl = "amqp://guest:guest@localhost:5672/"
	}
	maxAttempts := 10
	attempts := 0
	conn, err := amqp.Dial(config.RMQClient.RabbitMQURL())
	for err != nil && attempts < maxAttempts {
		log.Printf("Failed to connect to RabbitMQ. Retrying in 5 seconds... (attempt %d/%d)", attempts+1, maxAttempts)
		attempts++
		time.Sleep(5 * time.Second)
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"UpdateTransactions", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Handle message
			err := controller.HandleMessage(d)
			if err != nil {
				logger.Log(err)
				continue
			}
			d.Ack(false) //only main msg should be acknowledged
		}
	}()

	// Set up the Gin router
	router := gin.Default()
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.GET("transactions", apiClient.GetAllTransactions)
	router.GET("addresstxs", apiClient.GetAllTransactionsByAddress)
	router.GET("clienttxs", apiClient.GetAllTransactionsByClientId)

	addr := fmt.Sprintf(":%s", port)
	router.Run(addr)
}
