package sematext

import (
	"github.com/gin-gonic/gin"
	"github.com/Jeffail/gabs"
	"log"
	"errors"
	"github.com/oscp/cloud-selfservice-portal-backend/server/common"
	"net/http"
	"strings"
)

const (
	genericAPIError    = "Fehler beim Aufruf der Sematext-API. Bitte erstelle ein Ticket."
	sematextRoleActive = "ACTIVE"
)

func getLogseneAppsHandler(c *gin.Context) {
	mail := common.GetUserMail(c)

	if appList, err := getAllLogseneAppsForUser(mail); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, appList)
	}
}

func getAllLogseneAppsForUser(userMail string) ([]common.SematextAppList, error) {
	appData, err := getAllLogseneApps()
	if err != nil {
		return nil, err
	}

	// Filter apps where user has an active role
	allApps, err := appData.Path("data.apps").Children()
	if err != nil {
		log.Println("error getting data inside json", err.Error())
		return nil, errors.New(genericAPIError)
	}

	userApps := []common.SematextAppList{}
	for _, app := range allApps {
		appName := app.Path("name").Data().(string)
		userRoles, err := app.Path("userRoles").Children()
		if err != nil {
			log.Println("userRoles not found for current app: ", appName)
			continue
		}

		for _, userRole := range userRoles {
			mail := userRole.S("userEmail").Data().(string)
			role := userRole.S("role").Data().(string)
			status := userRole.S("roleStatus").Data().(string)

			if strings.ToLower(mail) == strings.ToLower(userMail) &&
				status == sematextRoleActive {

				userApps = append(userApps, common.SematextAppList{
					Name: appName,
					PlanName: app.Path("plan.name").Data().(string),
					UserRole: role,
					IsFree: app.Path("plan.free").Data().(bool),
					MaxDailyEvents: app.Path("plan.maxDailyEvents").Data().(float64),
					PricePerDay: app.Path("plan.pricePerDay").Data().(float64),
					BillingInfo: app.Path("description").Data().(string),
				})
			}
		}
	}

	return userApps, nil
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
