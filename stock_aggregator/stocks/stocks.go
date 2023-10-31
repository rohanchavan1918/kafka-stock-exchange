package stocks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/rohanchavan1918/stock_aggregator/conf"
	"github.com/rohanchavan1918/stock_aggregator/utils"
	"github.com/segmentio/kafka-go"
)

type Stock struct {
	// Stock model
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var StockChannel = make(chan Stock, 1)

func KafkaStockWriterWorker(stockChannel <-chan Stock, wg *sync.WaitGroup) {
	// Worker to write stock to kafka
	defer wg.Done()
	for {
		stock := <-stockChannel
		utils.LogInfo("Adding stock to kafka : %v", stock)
		err := AddToKafka(&stock)
		if err != nil {
			utils.LogInfo("Failed to add stock to kafka : %s", err)
		}
	}
}

func (s *Stock) Validate() error {
	// validate stock to database
	if s.Name == "" {
		return errors.New("Name cannot be empty")
	}
	if s.ID == 0 {
		return errors.New("ID cannot be empty")
	}

	if s.Price < 0 || math.IsNaN(s.Price) || math.IsInf(s.Price, 0) {
		return errors.New("Price must be a non-negative number")
	}
	return nil
}

func (s *Stock) GetPriceByte() (string, error) {
	// s.Price is float64 which cant be converted directly to []byte,
	// so we use convert it into string and then to []byte

	price := fmt.Sprintf("%f", s.Price)
	if price == "" {
		return "", errors.New("Price cannot be empty")
	}
	return price, nil
}

func AddToKafka(stock *Stock) error {
	// Add stock to kafka
	price, err := stock.GetPriceByte()
	if err != nil {
		return err
	}

	km := kafka.Message{
		Key:   []byte(stock.Name),
		Value: []byte(price),
	}

	w, err := conf.AppConfig.KafkaConfig.GetProducer()
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(), km)
	if err != nil {
		return err
	}

	return nil
}

func ConsumeFromKafka(channel chan Stock, wg *sync.WaitGroup) {
	// Create a new Kafka consumer
	defer wg.Done()
	consumer, err := conf.AppConfig.KafkaConfig.GetConsumer()

	if err != nil {
		log.Printf("Error while creating consumer: %v", err)
		utils.AlertAndPanic(err)
		return
	}

	// Continuously poll for new messages
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			continue
		}

		// Add the message value to the channel
		fmt.Println("Message received: ", string(msg.Value))
		float64_val, err := utils.StringToFloat64(string(msg.Value))
		if err != nil {
			log.Printf("Error while converting message value to float64: %v", err)
			utils.AlertAndPanic(err)
		}
		stock := Stock{
			Name:  string(msg.Key),
			Price: float64_val,
		}
		channel <- stock
	}
}

func KafkaStockReaderWorker(stockChannel <-chan Stock, wg *sync.WaitGroup) {
	// Worker to write stock to kafka
	defer wg.Done()
	for {
		stock := <-stockChannel
		utils.LogInfo("Added stock to DB : %v", stock)

		// err := AddToKafka(&stock)
		// if err != nil {
		// 	utils.LogInfo("Failed to add stock to kafka : %s", err)
		// }
	}
}
