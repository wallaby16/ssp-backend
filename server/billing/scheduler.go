package billing

import (
	"log"

	"fmt"
	"os"
	"strings"

	"time"

	"github.com/oscp/cloud-selfservice-portal/server/openshift"
	"golang.org/x/tools/go/gcimporter15/testdata"
)

func StartBillingScheduler() {
	// Do every hour
	fetchProjectList()

	getAllMissingDatapoints()

	fetchQuotas()
	fetchRequests()
	fetchEffectiveUsage()
	fetchNewrelicUsage()
	fetchSematextUsage()
}

func fetchProjectList() {
	// Get project list from OpenShift and add to etcd
	oseProjects, err := openshift.GetProjectList()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get projects that should be ignored in billing
	ignoreConfig := os.Getenv("BILLING_IGNORE_PROJECTS")
	var ignoreProjects []string
	if ignoreConfig != "" {
		log.Println("Project to ignore in billing:", ignoreConfig)
		ignoreProjects = strings.Split(ignoreConfig, ",")
	}

	children, err := oseProjects.S("items").Children()
	if err != nil {
		log.Fatal("Error getting project-children inside json object: " + err.Error())
	}

	// Loop project list and add to etcd if necessary
	activeOseProjects := []string{}
	for _, p := range children {
		name := p.Path("metadata.name").String()
		if !contains(ignoreProjects, name) {
			activeOseProjects = append(activeOseProjects, name)
			existingProject := getProject(name)

			// Set active if not already active
			if existingProject != nil && !existingProject.IsActive {
				fmt.Sprintf("Reactivating project because it is back in OSE: %v", name)
				existingProject.IsActive = true
				existingProject.Save()
			}

			if existingProject == nil {
				fmt.Sprintf("Adding new project from ose to etcd: %v", name)
				newProject := Project{
					IsActive:          true,
					Name:              name,
					BillingDatapoints: []Datapoint{},
					BillingNumber:     p.Path("metadata.annotations").S("openshift.io/kontierung-element").String(),
				}
				newProject.Save()
			}
		} else {
			fmt.Sprintf("Project %v was ignored becase it is on the ignore list", name)
		}
	}

	// Remove projects that are active in etcd but no longer in OSE
	dbProjects := getAllProjects()
	for _, p := range *dbProjects {
		if !contains(activeOseProjects, p.Name) {
			fmt.Sprintf("Setting project to inactive becase it's no longer in OSE. %v", p.Name)
			p.IsActive = false
			p.Save()
		}
	}
}

func getLastEntry(dps *[]Datapoint) time.Time {
	var lastEntry time.Time
	for _, dp := range *dps {
		if lastEntry.Before(dp.Time) {
			lastEntry = dp.Time
		}
	}

	return lastEntry
}

func getAllMissingDatapoints() {
	// For each project in etcd:
	projects := getAllProjects()
	for _, p := range *projects {
		last := getLastEntry(&p.BillingDatapoints)

		// Get data for the next hour
		last = last.Add(1 * time.Hour)

		for last.Before(time.Now()) {
			// Get data for this hour
			newDatapoint := Datapoint{
				Time: time.Now(),
			}

			// Get all relevant data
			fetchQuotas(p.Name, &newDatapoint)

			// Save the new datapoint
			p.BillingDatapoints = append(p.BillingDatapoints, newDatapoint);
			p.Save()

			last = last.Add(1 * time.Hour)
		}
	}
}

func fetchQuotas(projectName string, dp *Datapoint) {
	cpuQuota, memQuota := openshift.GetQuotas(projectName)
	dp.QuotaCPU = float64(cpuQuota)
	dp.QuotaMemory = float64(memQuota)
}

func fetchRequests() {
	// For each project in etcd:
	// Check last entry, interpolate if necessary
	// Get current requests, add to etcd
}

func fetchEffectiveUsage() {
	// For each project in etcd:
	// Check last entry, get if necessary
	// Get usage, add to etcd
}

func fetchNewrelicUsage() {
	// For all project in etcd in one request
	// Check last entry, interpolate if necessary
	// Get APM (CU), Synthetics Count, Browser, Mobile Usage
}

func fetchSematextUsage() {
	// For each project in etcd
	// Check last entry, interpolate if necessary
	// Get current plan & dollar per month
}
