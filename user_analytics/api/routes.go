package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/rohanchavan1918/user_analytics/api/v1"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1.SetupRoutes(api)
}
