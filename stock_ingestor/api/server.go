package api

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rohanchavan1918/stock_ingestor/conf"
	"github.com/rohanchavan1918/stock_ingestor/stocks"
	"github.com/rohanchavan1918/stock_ingestor/utils"
)

func RunServer(config *conf.Config) {
	r := gin.Default()
	SetupRoutes(r)
	dbConn := conf.GetDBConnection(&config.DB)
	err := dbConn.Ping()
	if err != nil {
		utils.AlertAndPanic(err)
	}

	// Once DB Connection is validated, add it to the global connections
	conf.AppConnections.DB = dbConn

	// Check for Kafka connections
	writer, err := conf.AppConfig.KafkaConfig.GetProducer()
	if err != nil {
		utils.AlertAndPanic(err)
	}

	conf.AppConnections.KafkaWriter = writer
	var wg sync.WaitGroup

	// Spawn KafkaStockWriterWorker goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go stocks.KafkaStockWriterWorker(stocks.StockChannel, &wg)
	}
	defer close(stocks.StockChannel)
	defer wg.Wait()

	port := fmt.Sprintf(":%s", strconv.Itoa(int(conf.AppConfig.Port)))
	r.Run(port)
}
