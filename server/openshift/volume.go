package openshift

import (
	"errors"
	"net/http"

	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"encoding/json"

	"os"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/glusterapi/models"
	"github.com/oscp/cloud-selfservice-portal/server/common"
)

const (
	wrongSizeFormatError = "Ungültige Grösse. Format muss Zahl gefolgt von M/G sein (z.B. 500M)."
	wrongSizeLimitError = "Grösse nicht erlaubt. Mindestgrösse: 500M. Maximale Grössen sind: M: %v, G: %v"
)

func newVolumeHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.NewVolumeCommand
	if c.BindJSON(&data) == nil {
		if err := validateNewVolume(data.Project, data.Size, data.PvcName, data.Mode, username); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		if err := createNewVolume(data.Project, username, data.Size, data.PvcName, data.Mode); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{
				Message: "Das Volume wurde erstellt. Deinem Projekt wurde das PVC, und der Gluster Service & Endpunkte hinzugefügt.",
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func fixVolumeHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.FixVolumeCommand
	if c.BindJSON(&data) == nil {
		if err := validateFixVolume(data.Project, username); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		if err := recreateGlusterObjects(data.Project, username); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{
				Message: "Die Gluster-Objekte wurden in deinem Projekt erzeugt.",
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func growVolumeHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.GrowVolumeCommand
	if c.BindJSON(&data) == nil {
		if err := validateGrowVolume(data.Project, data.NewSize, data.PvName, username); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		if err := growExistingVolume(data.Project, data.NewSize, data.PvName, username); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{Message: "Das Volume wurde vergrössert."})
		}

	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func validateNewVolume(project string, size string, pvcName string, mode string, username string) error {
	// Required fields
	if len(project) == 0 || len(pvcName) == 0 || len(size) == 0 || len(mode) == 0 {
		return errors.New("Es müssen alle Felder ausgefüllt werden")
	}

	if err := validateMaxSize(size); err != nil {
		return err
	}

	// Permissions on project
	if err := checkAdminPermissions(username, project); err != nil {
		return err
	}

	// Check if pvc name already taken
	if err := checkPvcName(project, pvcName); err != nil {
		return err
	}

	return nil
}

func validateGrowVolume(project string, newSize string, pvName string, username string) error {
	// Required fields
	if len(project) == 0 || len(pvName) == 0 || len(newSize) == 0 {
		return errors.New("Es müssen alle Felder ausgefüllt werden")
	}

	if err := validateMaxSize(newSize); err != nil {
		return err
	}

	// Permissions on project
	if err := checkAdminPermissions(username, project); err != nil {
		return err
	}

	return nil
}

func validateFixVolume(project string, username string) error {
	if len(project) == 0 {
		return errors.New("Projekt muss angegeben werden")
	}

	// Permissions on project
	if err := checkAdminPermissions(username, project); err != nil {
		return err
	}

	return nil
}

func validateMaxSize(size string) error {
	minMB := 500
	maxMB := 1024
	maxGB := os.Getenv("MAX_VOLUME_GB")

	maxGBInt, errGB := strconv.Atoi(maxGB)
	if errGB != nil || maxGBInt <= 0 {
		log.Fatal("Env variable 'MAX_VOLUME_GB' must be specified and a valid integer")
	}

	// Size limits
	if strings.HasSuffix(size, "M") {
		sizeInt, err := strconv.Atoi(strings.Replace(size, "M", "", 1))
		if err != nil {
			return errors.New(wrongSizeFormatError)
		}

		if sizeInt < minMB {
			return fmt.Errorf(wrongSizeLimitError, maxMB, maxGB)
		}
		if sizeInt > maxMB {
			return errors.New("Deine Angaben sind zu gross für 'M'. Bitte gib die Grösse als Ganzzahl in 'G' an")
		}
	}
	if strings.HasSuffix(size, "G") {
		sizeInt, err := strconv.Atoi(strings.Replace(size, "G", "", 1))
		if err != nil {
			return errors.New(wrongSizeFormatError)
		}

		if sizeInt > maxGBInt {
			return fmt.Errorf(wrongSizeLimitError, maxMB, maxGB)
		}
	}

	return nil
}

func checkPvcName(project string, pvcName string) error {
	client, req := getOseHTTPClient("GET", "api/v1/namespaces/" + project + "/persistentvolumeclaims", nil)
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error from server while getting pvc-list: ", err.Error())
		return errors.New(genericAPIError)
	}

	defer resp.Body.Close()

	json, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		log.Println("error parsing body of response:", err)
		return errors.New(genericAPIError)
	}

	// Check if pvc name is not already used
	children, err := json.S("items").Children()
	if err != nil {
		log.Println("Unable to parse pvc list", err.Error())
		return errors.New(genericAPIError)
	}
	for _, v := range children {
		if v.Path("metadata.name").Data().(string) == pvcName {
			return fmt.Errorf("Der gewünschte PVC-Name %v existiert bereits.", pvcName)
		}
	}

	return nil
}

func createNewVolume(project string, username string, size string, pvcName string, mode string) error {
	pvName, err := createGlusterVolume(project, size, username)
	if err != nil {
		return err
	}

	if err := createOpenShiftPV(size, pvName, mode, username); err != nil {
		return err
	}

	if err := createOpenShiftPVC(project, size, pvcName, mode, username); err != nil {
		return err
	}

	// Create Gluster Service & Endpoints in user project
	if err := createOpenShiftGlusterService(project, username); err != nil {
		return err
	}
	if err := createOpenShiftGlusterEndpoint(project, username); err != nil {
		return err
	}

	return nil
}

func createGlusterVolume(project string, size string, username string) (string, error) {
	cmd := models.CreateVolumeCommand{
		Project: project,
		Size:    size,
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		log.Println(err.Error())
		return "", errors.New(genericAPIError)
	}

	client, req := getGlusterHTTPClient("sec/volume", b)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling gluster-api", err.Error())
		return "", errors.New(genericAPIError)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("%v created a gluster volume. Project: %v, size: %v", username, project, size)

		respJson, err := gabs.ParseJSONBuffer(resp.Body)
		if err != nil {
			log.Println("Error parsing respJson from gluster-api response", err.Error())
			return "", errors.New(genericAPIError)
		}

		// Add gl_ to pvName because of conflicting PVs on other storage technology
		return fmt.Sprintf("gl_%v", respJson.Path("message").Data().(string)), nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating gluster volume:", err, resp.StatusCode, string(errMsg))

	return "", fmt.Errorf("Fehlerhafte Antwort vom Gluster-API: %v", string(errMsg))
}

func growExistingVolume(project string, newSize string, pvName string, username string) error {
	// Renaming Rules:
	// OpenShift cannot use _ in names. Thus the pvName will be gl-<project>-pv<number>
	// 1. Remove gl-
	// 2. Change -pv to _pv
	cmd := models.GrowVolumeCommand{
		PvName:  strings.Replace(strings.Replace(pvName, "gl-", "", 1), "-pv", "_pv", 1),
		NewSize: newSize,
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		log.Println(err.Error())
		return errors.New(genericAPIError)
	}

	client, req := getGlusterHTTPClient("sec/volume/grow", b)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling gluster-api", err.Error())
		return errors.New(genericAPIError)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("%v grew gluster volume. Project: %v, newSize: %v", username, project, newSize)
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error growing gluster volume:", err, resp.StatusCode, string(errMsg))

	return fmt.Errorf("Fehlerhafte Antwort vom Gluster-API: %v", string(errMsg))
}

func createOpenShiftPV(size string, pvName string, mode string, username string) error {
	// The Volume will use _ in the name, OpenShift can't, so we change it to -
	p := newObjectRequest("PersistentVolume", strings.Replace(pvName, "_", "-", -1))

	p.SetP(size, "spec.capacity.storage")
	p.SetP("glusterfs-cluster", "spec.glusterfs.endpoints")

	// The gluster volume starts with vol_ instead of gl_
	p.SetP(strings.Replace(pvName, "gl_", "vol_", 1), "spec.glusterfs.path")
	p.SetP(false, "spec.glusterfs.readOnly")
	p.SetP("Retain", "spec.persistentVolumeReclaimPolicy")

	p.ArrayP("spec.accessModes")
	p.ArrayAppend(mode, "spec", "accessModes")

	client, req := getOseHTTPClient("POST",
		"api/v1/persistentvolumes",
		bytes.NewReader(p.Bytes()))

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusCreated {
		log.Printf("Created the pv %v based on the request of %v", pvName, username)
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating new PV:", err, resp.StatusCode, string(errMsg))

	return errors.New(genericAPIError)
}

func createOpenShiftPVC(project string, size string, pvcName string, mode string, username string) error {
	p := newObjectRequest("PersistentVolumeClaim", strings.Replace(pvcName, "_", "-", -1))

	p.SetP(size, "spec.resources.requests.storage")
	p.ArrayP("spec.accessModes")
	p.ArrayAppend(mode, "spec", "accessModes")

	client, req := getOseHTTPClient("POST",
		"api/v1/namespaces/"+project+"/persistentvolumeclaims",
		bytes.NewReader(p.Bytes()))

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusCreated {
		log.Printf("Created the pvc %v based on the request of %v", pvcName, username)
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating new PVC:", err, resp.StatusCode, string(errMsg))

	return errors.New(genericAPIError)
}

func recreateGlusterObjects(project string, username string) error {
	if err := createOpenShiftGlusterService(project, username); err != nil {
		return err
	}

	if err := createOpenShiftGlusterEndpoint(project, username); err != nil {
		return err
	}

	return nil
}

func createOpenShiftGlusterService(project string, username string) error {
	p := newObjectRequest("Service", "glusterfs-cluster")

	port := gabs.New()
	port.Set(1, "port")

	p.ArrayP("spec.ports")
	p.ArrayAppendP(port.Data(), "spec.ports")

	client, req := getOseHTTPClient("POST",
		"api/v1/namespaces/"+project+"/services",
		bytes.NewReader(p.Bytes()))

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusCreated {
		log.Printf("Created the gluster service based on the request of %v", username)
		return nil
	}

	if resp.StatusCode == http.StatusConflict {
		log.Println("Gluster service already existed, skipping")
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating gluster service:", err, resp.StatusCode, string(errMsg))

	return errors.New(genericAPIError)
}

func createOpenShiftGlusterEndpoint(project string, username string) error {
	p, err := getGlusterEndpointsContainer()
	if err != nil {
		return err
	}

	client, req := getOseHTTPClient("POST",
		"api/v1/namespaces/"+project+"/endpoints",
		bytes.NewReader(p.Bytes()))

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusCreated {
		log.Printf("Created the gluster endpoints based on the request of %v", username)
		return nil
	}

	if resp.StatusCode == http.StatusConflict {
		log.Println("Gluster endpoints already existed, skipping")
		return nil
	}

	errMsg, _ := ioutil.ReadAll(resp.Body)
	log.Println("Error creating gluster endpoints:", err, resp.StatusCode, string(errMsg))

	return errors.New(genericAPIError)
}

func getGlusterEndpointsContainer() (*gabs.Container, error) {
	p := newObjectRequest("Endpoints", "glusterfs-cluster")
	p.Array("subsets")

	// Add gluster endpoints
	glusterIPs := os.Getenv("GLUSTER_IPS")
	if len(glusterIPs) == 0 {
		log.Println("Wrong configuration. Missing env variable 'GLUSTER_IPS'")
		return nil, errors.New(genericAPIError)
	}

	addresses := gabs.New()
	addresses.Array("addresses")
	addresses.Array("ports")
	for _, ip := range strings.Split(glusterIPs, ",") {
		address := gabs.New()
		address.Set(ip, "ip")

		addresses.ArrayAppend(address.Data(), "addresses")
	}

	port := gabs.New()
	port.Set(1, "port")
	addresses.ArrayAppend(port.Data(), "ports")

	p.ArrayAppend(addresses.Data(), "subsets")

	return p, nil
}
