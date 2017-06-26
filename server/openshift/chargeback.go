package openshift

import (
	"fmt"
	"net/http"
	"time"

	"crypto/tls"
	"log"
	"net/url"
	"os"

	"encoding/json"

	"math"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
	"strconv"
)

const (
	dateFormat  = "2006-01-02 15:04:05"
	newRelicAPI = "https://insights-api.newrelic.com/v1/accounts/1159282/query?nrql=%v"
)

func chargeBackHandler(c *gin.Context) {
	// Debug
	//dummyResponse := Resources{
	//	Project:           "test-projekt",
	//	Start:             time.Date(2017, time.April, 1, 0, 0, 0, 0, time.Local),
	//	End:               time.Date(2017, time.April, 30, 0, 0, 0, 0, time.Local),
	//	AccountAssignment: "1234",
	//	QuotaCPU:          38,
	//	QuotaMemory:       123,
	//	RequestedCPU:      0.012020203020402,
	//	RequestedMemory:   14.56662356236,
	//	Storage:           100,
	//	TotalUsedCPU:      0.84323,
	//	TotalUsedMemory:   245.214124125,
	//	UsageDataPoints: []UsageDataPoint{
	//		{
	//			End:        time.Now(),
	//			UsedCPU:    0.1234,
	//			UsedMemory: 214.0001234,
	//		},
	//	},
	//	Costs: Costs{
	//		QuotaCPU:        404,
	//		QuotaMemory:     327,
	//		Storage:         10,
	//		RequestedCPU:    1,
	//		RequestedMemory: 157,
	//		UsedCPU:         36,
	//		UsedMemory:      2618,
	//		Total:           3543,
	//	},
	//}
	//bytes, _ := json.Marshal(dummyResponse.UsageDataPoints)
	//
	//c.HTML(http.StatusOK, chargeBackURL, gin.H{
	//	"data":       dummyResponse,
	//	"dataPoints": string(bytes),
	//})
	//return
	project := c.PostForm("project")
	username := common.GetUserName(c)
	year := c.PostForm("year")
	month := c.PostForm("month")

	start, end, err := validateChargeback(project, username, year, month)
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"error": err.Error(),
		})
		return
	}

	quotas, err := getQuotasFromNewRelic(*start, *end, project)
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"error": err.Error(),
		})
		return
	}
	assignments, err := getAssignmentsFromNewRelic(project)
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"error": err.Error(),
		})
		return
	}

	usages, err := getUsagesFromNewRelic(*start, *end, project)
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Got data from newrelic, calculating cost structure")

	// Mangle the results to an array
	resources := calculateResources(quotas, assignments, *usages)

	resources.Project = project
	resources.Start = *start
	resources.End = *end

	// Calculate Costs for the resources
	calculatePrices(&resources)

	bytes, _ := json.Marshal(resources.UsageDataPoints)
	c.HTML(http.StatusOK, chargeBackURL, gin.H{
		"data":       resources,
		"dataPoints": string(bytes),
	})
}

func validateChargeback(project string, username string, year string, month string) (*time.Time, *time.Time, error) {
	// Validate project access
	//if err := checkAdminPermissions(username, project); err != nil {
	//	return nil, nil, err
	//}

	// Parse month & year
	yearI, err := strconv.Atoi(year)
	if err != nil {
		return nil, nil, errors.New("Invalid year: " + year)
	}
	monthI, err := strconv.Atoi(month)
	if err != nil {
		return nil, nil, errors.New("Invalid month: " + month)
	}

	start := time.Date(yearI, time.Month(monthI), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return &start, &end, nil
}

func calculatePrices(resources *Resources) {
	// Define prices
	// Todo from variables
	unitPrices := UnitPrices{
		quotaCPU:        10.0,
		quotaMemory:     2.5,
		requestedCPU:    40.0,
		requestedMemory: 10,
		usedCPU:         40,
		usedMemory:      10,
		storage:         1.0}

	// Todo from variable
	fee := 1.0625

	// Calculate single Costs
	resources.Costs = Costs{
		QuotaCPU:        math.Ceil(resources.QuotaCPU * unitPrices.quotaCPU * fee),
		QuotaMemory:     math.Ceil(resources.QuotaMemory * unitPrices.quotaMemory * fee),
		RequestedCPU:    math.Ceil(resources.RequestedCPU * unitPrices.requestedCPU * fee),
		RequestedMemory: math.Ceil(resources.RequestedMemory * unitPrices.requestedMemory * fee),
		UsedCPU:         math.Ceil(resources.TotalUsedCPU * unitPrices.usedCPU * fee),
		UsedMemory:      math.Ceil(resources.TotalUsedMemory * unitPrices.usedMemory * fee),
		Storage:         math.Ceil(resources.Storage * unitPrices.storage * fee),
	}

	// Calculate Total Costs
	resources.Costs.Total = resources.Costs.QuotaCPU + resources.Costs.QuotaMemory +
		resources.Costs.Storage + resources.Costs.RequestedCPU + resources.Costs.RequestedMemory +
		resources.Costs.UsedCPU + resources.Costs.UsedMemory
}

func calculateResources(quota *Quota, assignment *Assignment, usages []Usage) Resources {
	// Add Quotas and Requests to resources
	resources := Resources{QuotaCPU: quota.Results[0].Average,
		RequestedCPU:    quota.Results[1].Average,
		QuotaMemory:     quota.Results[2].Average,
		RequestedMemory: quota.Results[3].Average,
		Storage:         quota.Results[4].Average}

	// Add Assigned to resources
	resources.AccountAssignment = assignment.Results[0].Latest

	// Add usages to resources
	var usedCPU float64
	var usedMemory float64

	for _, usage := range usages {
		// Calculate average
		usedCPU += usage.Results[0].Result
		usedMemory += usage.Results[1].Result

		// Add data point for entry
		resources.UsageDataPoints = append(resources.UsageDataPoints, UsageDataPoint{
			UsedCPU:    usage.Results[0].Result,
			UsedMemory: usage.Results[1].Result,
			End:        usage.End,
		})
	}

	// Normalize usage
	resources.TotalUsedCPU = usedCPU / float64(len(usages))
	resources.TotalUsedMemory = usedMemory / float64(len(usages))

	return resources
}

func getQuotasFromNewRelic(start time.Time, end time.Time, project string) (*Quota, error) {
	quotaQuery := fmt.Sprintf(`
	   SELECT average(cpuHard) AS CpuQuota,
	   		  average(cpuUsed) AS CpuRequests,
	   		  average(memoryHard) AS MemoryQuota,
	   		  average(memoryUsed) AS MemoryRequests,
	   		  average(Storage) AS Storage
	   FROM %v
	   WHERE project like '%v'
	   SINCE '%s' UNTIL '%s'
	   WITH TIMEZONE 'Europe/Zurich' LIMIT 1000`, "OpenshiftViasQuota", project, start.Format(dateFormat), end.Format(dateFormat))

	client, req := getNewRelicClient(quotaQuery)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	quota := Quota{}
	if err := json.NewDecoder(resp.Body).Decode(&quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

func getAssignmentsFromNewRelic(project string) (*Assignment, error) {
	assignmentQuery := fmt.Sprintf(`
		SELECT latest(accountAssignment), latest(megaId)
		FROM %v
		WHERE project like '%v'
		LIMIT 1000`, "OpenshiftViasQuota", project)

	client, req := getNewRelicClient(assignmentQuery)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	assignment := Assignment{}
	if err := json.NewDecoder(resp.Body).Decode(&assignment); err != nil {
		return nil, err
	}
	return &assignment, nil
}

func getUsagesFromNewRelic(start time.Time, end time.Time, project string) (*[]Usage, error) {
	duration, _ := time.ParseDuration("24h")
	current := start

	usages := []Usage{}
	for current.Before(end) {
		timeFilter := fmt.Sprintf("SINCE '%s'", current.Format(dateFormat))
		current = current.Add(duration)

		if current.Before(end) {
			timeFilter += fmt.Sprintf(" UNTIL '%s'", current.Format(dateFormat))
		} else {
			// At the end
			timeFilter += fmt.Sprintf(" UNTIL '%s'", end.Format(dateFormat))
		}

		q := fmt.Sprintf(`
			SELECT rate(sum(cpuPercent), 54 minutes)/100 as CPU,
				   rate(sum(memoryResidentSizeBytes), 54 minutes)/(1000*1000*1000) as GB
			FROM ProcessSample
			WHERE %v AND `+"`containerLabel_io.kubernetes.pod.namespace`"+`LIKE '%v'
			%v
			WITH TIMEZONE 'Europe/Zurich' LIMIT 1000`,
			"fullHostname like '%.sbb.ch'", project, timeFilter)

		client, req := getNewRelicClient(q)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		usage := Usage{}
		if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
			return nil, err
		}

		if current.Before(end) {
			usage.End = current
		} else {
			usage.End = end
		}

		usages = append(usages, usage)
		resp.Body.Close()
	}

	return &usages, nil
}

func getNewRelicClient(query string) (*http.Client, *http.Request) {
	token := os.Getenv("NEWRELIC_API_TOKEN")
	if len(token) == 0 {
		log.Fatal("Env variable 'NEWRELIC_API_TOKEN' must be specified")
	}

	proxy := os.Getenv("HTTP_PROXY")
	var tr *http.Transport
	if len(proxy) > 0 {
		proxyURL, _ := url.Parse(proxy)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", fmt.Sprintf(newRelicAPI, url.QueryEscape(query)), nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Query-Key", token)

	return client, req
}
