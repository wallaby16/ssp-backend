package billing

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coreos/etcd/client"
)

const (
	jsonDecodingError = "Error decoding json from etcd. %v"
	jsonEncodingError = "Error encoding json to save in etcd. %v"
	readError         = "Error reading %v from etcd. Msg: %v"
	saveError         = "Error while saving %v to etcd. Msg: %v"
)

type Project struct {
	Name              string
	BillingNumber     string
	IsActive          bool
	BillingDatapoints []Datapoint
}

func getProject(name string) *Project {
	data, err := Api.Get(context.Background(), "projects/"+name, nil)
	if err != nil {
		if !client.IsKeyNotFound(err) {
			log.Fatalf(readError, "project", err.Error())
		}
	}
	var p Project
	err = json.Unmarshal([]byte(data.Node.Value), &p)
	if err != nil {
		log.Fatalf(jsonDecodingError, err.Error())
	}
	return &p
}

func getAllProjects() *[]Project {
	data, err := Api.Get(context.Background(), "projects", nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			// First run return empty list
			return &[]Project{}
		}

		log.Fatalf(readError, "projects", err.Error())
	}

	var list []Project
	err = json.Unmarshal([]byte(data.Node.Value), &list)
	if err != nil {
		log.Fatalf(jsonDecodingError, err.Error())
	}
	return &list
}

func (p Project) Save() {
	json, err := json.Marshal(p)
	if err != nil {
		log.Fatalf(jsonEncodingError, err.Error())
	}
	_, err = Api.Set(context.Background(), "projects/"+p.Name, string(json), nil)
	if err != nil {
		log.Fatalf(saveError, "project", err.Error())
	}
}

type Datapoint struct {
	Time                    time.Time
	QuotaCPU                float64
	QuotaMemory             float64
	RequestedCPU            float64
	RequestedMemory         float64
	UsedCPU                 float64
	UsedMemory              float64
	NewrelicApmCU           float64
	NewrelicSyntheticsCount float64
	NewrelicBrowser         float64
	NewrelicMobile          float64
	SematextPlan            string
	SematextMonthlyCost     float64
}


// TODO: WEG?
type UnitPrices struct {
	quotaCPU        float64
	quotaMemory     float64
	storage         float64
	requestedCPU    float64
	requestedMemory float64
	usedCPU         float64
	usedMemory      float64
}

type Resources struct {
	Project           string
	Start             time.Time
	End               time.Time
	AccountAssignment string
	QuotaCPU          float64
	QuotaMemory       float64
	RequestedCPU      float64
	RequestedMemory   float64
	Storage           float64
	TotalUsedCPU      float64
	TotalUsedMemory   float64
	UsageDataPoints   []UsageDataPoint
	Costs             Costs
}
type UsageDataPoint struct {
	UsedCPU    float64
	UsedMemory float64
	End        time.Time
}

type Costs struct {
	QuotaCPU        float64
	QuotaMemory     float64
	Storage         float64
	RequestedCPU    float64
	RequestedMemory float64
	UsedCPU         float64
	UsedMemory      float64
	Total           float64
}

// NewRelic-API Types
type Quota struct {
	Results []QuotaResult
}
type QuotaResult struct {
	Average float64
}
type Assignment struct {
	Results []AssignmentResult
}
type AssignmentResult struct {
	Latest string
}
type Usage struct {
	Results []UsageResult
	End     time.Time
}
type UsageResult struct {
	Result float64
}
