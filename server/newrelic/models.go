package newrelic

type UnitPrices struct {
	quotaCpu        float64
	quotaMemory     float64
	storage         float64
	requestedCpu    float64
	requestedMemory float64
	usedCpu         float64
	usedMemory      float64
}
// NewRelic-API Types
type Quota struct {
	Facets   []QuotaFacet
	Metadata Metadata
}
type QuotaFacet struct {
	Name    string
	Results []QuotaResult
}
type QuotaResult struct {
	Average float64
}
type Metadata struct {
	Contents Contents
}
type Contents struct {
	Contents []Content
}
type Content struct {
	Alias string
}
type Assignment struct {
	Facets   []AssignmentFacet
	Metadata Metadata
}
type AssignmentFacet struct {
	Name    string
	Results []AssignmentResult
}
type AssignmentResult struct {
	Latest string
}
type Usage struct {
	Facets   []UsageFacet
	Metadata Metadata
}
type UsageFacet struct {
	Name    string
	Results []UsageResult
}

type UsageResult struct {
	Result float64
}
