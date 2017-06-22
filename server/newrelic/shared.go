package newrelic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	chargeBackURL      = "chargeback.html"
)

// RegisterRoutes registers the routes for NewRelic
func RegisterRoutes(r *gin.RouterGroup) {
	// Quotas
	r.GET("/newrelic/chargeback", func(c *gin.Context) {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{})
	})
	r.POST("/newrelic/chargeback", chargeBackHandler)
}


