package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
	"github.com/oscp/cloud-selfservice-portal/server/openshift"
	"github.com/oscp/cloud-selfservice-portal/server/ddc"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Public routes
	authMiddleware := common.GetAuthMiddleware()
	router.POST("/login", authMiddleware.LoginHandler)

	// Feature toggle
	router.GET("/config", common.ConfigHandler)

	// Protected routes
	auth := router.Group("/api/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Openshift routes
		openshift.RegisterRoutes(auth)

		// DDC routes
		auth.GET("/ddc/billing", ddc.GetDDCBillingHandler)
	}

	log.Println("Cloud SSP is running")
	router.Run()
}
