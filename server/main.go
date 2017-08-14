package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
	"github.com/oscp/cloud-selfservice-portal/server/gluster"
	"github.com/oscp/cloud-selfservice-portal/server/openshift"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	// Public routes
	authMiddleware := common.GetAuthMiddleware()
	router.POST("/login", authMiddleware.LoginHandler)

	// Protected routes
	auth := router.Group("/api/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Openshift routes
		openshift.RegisterRoutes(auth)

		// Gluster routes
		gluster.RegisterRoutes(auth)
	}

	log.Println("Cloud SSP is running")
	router.Run()
}
