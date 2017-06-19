package newrelic

import (
	"github.com/gin-gonic/gin"
)

const (
	chargeBackURL      = "chargeback.html"
)

// RegisterRoutes registers the routes for NewRelic
func RegisterRoutes(r *gin.RouterGroup) {
	// Quotas
	r.GET("/newrelic/chargeback", func(c *gin.Context) {
		chargeBackHandler(c)
		//c.HTML(http.StatusOK, chargeBackURL, gin.H{})
	})
}


