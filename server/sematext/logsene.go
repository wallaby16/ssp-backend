package sematext

import (
	"github.com/gin-gonic/gin"
	"github.com/Jeffail/gabs"
	"log"
	"errors"
	"github.com/oscp/cloud-selfservice-portal-backend/server/common"
)

const (
	genericAPIError    = "Fehler beim Aufruf der Sematext-API. Bitte erstelle ein Ticket."
)

func getLogseneAppsHandler(c *gin.Context) {
	mail := common.GetUserMail(c)
	getAllLogseneAppsForUser(mail)
}

func getAllLogseneAppsForUser(username string) (*gabs.Container, error) {
	// Get mail for user


	getAllLogseneApps()

	// Filter apps where user has an active role


	return nil, nil
}

func getAllLogseneApps() (*gabs.Container, error) {
	client, req := getSematextHTTPClient("GET", "users-web/api/v3/apps/users", nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error from Sematext API: ", err.Error())
		return nil, errors.New(genericAPIError)
	}

	defer resp.Body.Close()

	json, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		log.Println("error parsing body of response:", err)
		return nil, errors.New(genericAPIError)
	}

	return json, nil
}
