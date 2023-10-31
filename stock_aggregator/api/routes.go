package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/rohanchavan1918/stock_aggregator/api/v1"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1.SetupRoutes(api)
}
