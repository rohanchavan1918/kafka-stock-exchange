package stocks

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/rohanchavan1918/stock_ingestor/conf"
	"github.com/rohanchavan1918/stock_ingestor/utils"
	"github.com/segmentio/kafka-go"
)

type Stock struct {
	// Stock model
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var StockChannel = make(chan Stock)

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
