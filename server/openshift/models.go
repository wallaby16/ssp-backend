package openshift

import "time"

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
