package openshift

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/SchweizerischeBundesbahnen/ssp-backend/server/common"
)

func newServiceAccountHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.NewServiceAccountCommand
	if c.BindJSON(&data) == nil {
		if err := validateNewServiceAccount(username, data.Project, data.ServiceAccount); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		if err := createNewServiceAccount(username, data.Project, data.ServiceAccount); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{
				Message: fmt.Sprintf("Der Service Account %v wurde angelegt", data.ServiceAccount),
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func validateNewServiceAccount(username string, project string, serviceAccountName string) error {
	if len(serviceAccountName) == 0 {
		return errors.New("Service Account muss angegeben werden")
	}

	// Validate permissions
	if err := checkAdminPermissions(username, project); err != nil {
		return err
	}

	return nil
}

func createNewServiceAccount(username string, project string, serviceaccount string) error {
	p := newObjectRequest("ServiceAccount", serviceaccount)

	client, req := getOseHTTPClient("POST",
		"api/v1/namespaces/"+project+"/serviceaccounts",
		bytes.NewReader(p.Bytes()))

	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusCreated {
		resp.Body.Close()
		log.Print(username + " created a new service account: " + serviceaccount + " on project " + project)
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating new project:", err, resp.StatusCode, string(errMsg))
	return errors.New(genericAPIError)
}
