package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rohanchavan1918/stock_ingestor/stocks"
)

func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func AddStock(c *gin.Context) {
	// Api endpoint to add stock
	// curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "name": "test", "price": 1.0}' http://localhost:8080/api/v1/stock
	// Get post body
	var stock stocks.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := stock.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// err := stocks.AddToKafka(&stock)
	// if err != nil {
	// 	utils.LogInfo("Failed to add stock to kafka : %s", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Failed to add stock.",
	// 	})
	// 	return
	// }

	stocks.StockChannel <- stock

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully added stock.",
	})
	return

}
