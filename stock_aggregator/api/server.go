package api

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rohanchavan1918/stock_aggregator/conf"
	"github.com/rohanchavan1918/stock_aggregator/stocks"
	"github.com/rohanchavan1918/stock_aggregator/utils"
)

func RunServer(config *conf.Config) {
	r := gin.Default()
	SetupRoutes(r)
	dbConn := conf.GetDBConnection(&config.DB)
	err := dbConn.Ping()

	// Once DB Connection is validated, add it to the global connections
	conf.AppConnections.DB = dbConn

	// Check for Kafka connections
	writer, err := conf.AppConfig.KafkaConfig.GetProducer()
	if err != nil {
		utils.AlertAndPanic(err)
	}

	conf.AppConnections.KafkaWriter = writer
	var wg sync.WaitGroup

	utils.LogInfo("Starting KafkaStockReaderWorker goroutines")
	// Spawn KafkaStockReaderWorker goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go stocks.KafkaStockReaderWorker(stocks.StockChannel, &wg)
	}

	utils.LogInfo("Starting ConsumeFromKafka goroutines")
	// Spawn ConsumeFromKafka goroutines
	for i := 0; i < 5; i++ {
		go stocks.ConsumeFromKafka(stocks.StockChannel, &wg)
	}

	defer close(stocks.StockChannel)
	defer wg.Wait()

	port := fmt.Sprintf(":%s", strconv.Itoa(int(conf.AppConfig.Port)))
	r.Run(port)
}
