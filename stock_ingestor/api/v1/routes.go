package v1

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	v1Group := r.Group("/v1")
	v1Group.GET("/healthcheck", Healthcheck)
	v1Group.POST("/stock", AddStock)
}
