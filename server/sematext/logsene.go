package sematext

import (
	"github.com/gin-gonic/gin"
	"github.com/Jeffail/gabs"
	"log"
	"errors"
	"github.com/oscp/cloud-selfservice-portal-backend/server/common"
	"net/http"
	"strings"
	"bytes"
	"strconv"
	"fmt"
	"io/ioutil"
)

const (
	genericAPIError    = "Fehler beim Aufruf der Sematext-API. Bitte erstelle ein Ticket."
	sematextRoleActive = "ACTIVE"
	noAccessError      = "Du hast keinen Zugriff auf diese Sematext-Anwendung"
)

func getLogseneAppsHandler(c *gin.Context) {
	mail := common.GetUserMail(c)
	username := common.GetUserName(c)

	fmt.Sprintf("User %v listed all his sematext logsene apps", username)

	if appList, err := getAllLogseneAppsForUser(mail); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, appList)
	}
}

func updateLogseneAppHandler(c *gin.Context) {
	username := common.GetUserName(c)
	mail := common.GetUserMail(c)
	appId, err := strconv.Atoi(c.Param("appId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
		return
	}

	var data common.EditBillingDataCommand
	if c.BindJSON(&data) == nil {
		if err := validateLogseneAppEdit(mail, appId, data.Project, data.Billing); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		if err := updateBillingInfo(username, data.Billing, data.Project, appId); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{
				Message: fmt.Sprintf("Die Kontierungsdaten (%v / %v) wurden gespeichert.", data.Billing, data.Project),
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func validateLogseneAppEdit(mail string, appId int, project string, billing string) error {
	// Check permissions
	err := validateLogseneAppPermissions(mail, appId)
	if err != nil {
		return err
	}

	// Check values
	if len(project) == 0 {
		return errors.New("Das Projekt muss angegeben werden!")
	}

	if len(billing) == 0 {
		return errors.New("Die Kontierungsnummer muss angegeben werden!")
	}

	return nil
}

func validateLogseneAppPermissions(mail string, appId int) error {
	userApps, err := getAllLogseneAppsForUser(mail)

	if err != nil {
		return err
	}

	for _, a := range userApps {
		if a.AppId == appId {
			return nil
		}
	}

	return errors.New(noAccessError)
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
					AppId:          int(app.Path("id").Data().(float64)),
					Name:           appName,
					PlanName:       app.Path("plan.name").Data().(string),
					UserRole:       role,
					IsFree:         app.Path("plan.free").Data().(bool),
					MaxDailyEvents: app.Path("plan.maxDailyEvents").Data().(float64),
					PricePerDay:    app.Path("plan.pricePerDay").Data().(float64),
					BillingInfo:    app.Path("description").Data().(string),
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

func updateBillingInfo(username string, billing string, project string, appId int) error {
	fmt.Sprintf("User %v updated logsene app billing to %v / %v.", username, billing, project)

	json := gabs.New()
	json.Set(billing+" / "+project, "description")

	client, req := getSematextHTTPClient("PUT", "users-web/api/v3/apps/"+strconv.Itoa(appId), bytes.NewReader(json.Bytes()))
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error from Sematext API: ", err.Error())
		return errors.New(genericAPIError)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Println("Sematext response status code was: ", resp.StatusCode, string(bodyBytes))

	return errors.New(genericAPIError)
}
