package newrelic

import (
	"fmt"
	"net/http"
	"time"

	"crypto/tls"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"encoding/json"
)

const (
	dateFormat  = "2006-01-02 15:04:05"
	newRelicAPI = "https://insights-api.newrelic.com/v1/accounts/1159282/query?nrql=%v"
)

func chargeBackHandler(c *gin.Context) {
	// Todo get me from gin
	start := time.Date(2017, time.April, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2017, time.April, 30, 0, 0, 0, 0, time.Local)

	// Define prices
	// Todo from Variables
	//unitPrices := UnitPrices{
	//	quotaCpu:        10.0,
	//	quotaMemory:     2.5,
	//	requestedCpu:    40.0,
	//	requestedMemory: 10,
	//	usedCpu:         40,
	//	usedMemory:      10,
	//	storage:         1.0}

	// Todo from Variables
	//mgmtFee := 1.0625

	quotas, err := getQuotasFromNewRelic(start, end, "nova-prod")
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"Error": err.Error(),
		})
	}
	assignments, err := getAssignmentsFromNewRelic("nova-prod")
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"Error": err.Error(),
		})
	}

	usages, err := getUsagesFromNewRelic(start, end, "nova-prod")
	if err != nil {
		c.HTML(http.StatusOK, chargeBackURL, gin.H{
			"Error": err.Error(),
		})
	}

	log.Println(quotas, usages, assignments)

	c.HTML(http.StatusOK, chargeBackURL, gin.H{})
}

func getQuotasFromNewRelic(start time.Time, end time.Time, project string) (*Quota, error) {
	quotaQuery := fmt.Sprintf(`
	   SELECT average(cpuHard) AS CpuQuota,
	   		  average(cpuUsed) AS CpuRequests,
	   		  average(memoryHard) AS MemoryQuota,
	   		  average(memoryUsed) AS MemoryRequests,
	   		  average(storage) AS Storage
	   FROM %v
	   FACET project WHERE project like '%v'
	   SINCE '%s' UNTIL '%s'
	   WITH TIMEZONE 'Europe/Zurich' LIMIT 1000`, "OpenshiftViasQuota", project, start.Format(dateFormat), end.Format(dateFormat))

	client, req := getNewRelicClient(quotaQuery)

	resp, err :=  client.Do(req)
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
		FACET project WHERE project like '%v'
		LIMIT 1000`, "OpenshiftViasQuota", project)

	client, req := getNewRelicClient(assignmentQuery)

	resp, err :=  client.Do(req)
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
	duration, _ := time.ParseDuration("240h")
	current := start

	// Loop in 240h steps
	usageQueries := []string{}
	for current.Before(end) {
		timeFilter := fmt.Sprintf("SINCE '%s'", current.Format(dateFormat))
		current = current.Add(duration)

		if current.Before(end) {
			timeFilter += fmt.Sprintf(" UNTIL '%s'", current.Format(dateFormat))
		} else {
			// At the end
			timeFilter += fmt.Sprintf(" UNTIL '%s'", end.Format(dateFormat))
		}

		usageQueries = append(usageQueries, fmt.Sprintf(`
			SELECT rate(sum(cpuPercent), 54 minutes)/100 as CPU,
				   rate(sum(memoryResidentSizeBytes), 54 minutes)/(1000*1000*1000) as GB
			FROM ProcessSample
			FACET `+ "`containerLabel_io.kubernetes.pod.namespace`"+ `
			WHERE %v AND `+ "`containerLabel_io.kubernetes.pod.namespace`"+ `LIKE '%v'
			%v
			WITH TIMEZONE 'Europe/Zurich' LIMIT 1000`,
			"fullHostname like '%.sbb.ch'", project, timeFilter))
	}

	usages := []Usage{}
	for _, q := range usageQueries {
		client, req := getNewRelicClient(q)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		usage := Usage{}
		if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
			return nil, err
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
		proxyUrl, _ := url.Parse(proxy)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: http.ProxyURL(proxyUrl),
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

