package common

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ConfigHandler(c *gin.Context) {
	glusterApi := os.Getenv("GLUSTER_API_URL")
	ddcApi := os.Getenv("DDC_API")

	c.JSON(http.StatusOK, FeatureToggleResponse{
		DDC:     ddcApi != "",
		Gluster: glusterApi != "",
	})
}
