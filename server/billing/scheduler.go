package billing

import (
	"context"
	"log"

	"fmt"
	"os"
	"strings"

	"github.com/oscp/cloud-selfservice-portal/server/openshift"
)

const etcdError = "error accessing etcd db. Msg: "

func StartBillingScheduler() {
	// Do every hour
	fetchProjectList()

	fetchQuotas()
	fetchRequests()
	fetchEffectiveUsage()
	fetchNewrelicUsage()
	fetchSematextUsage()
}

func fetchProjectList() {
	// Get project list from OpenShift and add to etcd
	projects, err := openshift.GetProjectList()
	if err != nil {
		log.Fatal(err.Error())
	}
	ignoreConfig := os.Getenv("BILLING_IGNORE_PROJECTS")
	var ignoreProjects []string
	if ignoreConfig != "" {
		ignoreProjects = strings.Split(ignoreConfig, ",")
	}

	children, err := projects.S("items").Children()
	if err != nil {
		log.Fatal("Error getting project-children in json: " + err.Error())
	}
	// Loop project list and add to etcd if necessary
	for _, p := range children {
		name := p.Path("metadata.name").String()
		if !contains(ignoreProjects, name) {
			// Get existing project from etcd
			_, err := Api.Get(context.Background(), "projects/"+name, nil)
			if err != nil {

			}
		} else {
			fmt.Sprintf("Project %v was ignored becase it is on the ignore list", name)
		}

		log.Println(p.Path("metadata.name").String())
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()

	if err != nil {
		log.Fatal(etcdError + err.Error())
	}

	log.Println(resp)
}

func fetchQuotas() {
	// For each project in etcd:
	// Check last entry, interpolate if necessary
	// Get current quota, add to etcd
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
