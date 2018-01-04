package main

import (
	"log"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal-backend/server/aws"
	"github.com/oscp/cloud-selfservice-portal-backend/server/common"
	"github.com/oscp/cloud-selfservice-portal-backend/server/ddc"
	"github.com/oscp/cloud-selfservice-portal-backend/server/openshift"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	// Allow cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("*")
	router.Use(cors.New(corsConfig))

	// Public routes
	authMiddleware := common.GetAuthMiddleware()
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/s/")
	})
	router.StaticFS("/s", http.Dir("static"))
	router.POST("/login", authMiddleware.LoginHandler)
	router.GET("/config", common.ConfigHandler)

	// Protected routes
	auth := router.Group("/api/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Openshift routes
		openshift.RegisterRoutes(auth)

		// DDC routes
		ddc.RegisterRoutes(auth)

		// AWS routes
		aws.RegisterRoutes(auth)
	}

	log.Println("Cloud SSP is running")
	router.Run()
}
