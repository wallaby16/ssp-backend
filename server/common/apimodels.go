package common

type ProjectName struct {
	Project string `json:"project"`
}

type NewVolumeCommand struct {
	ProjectName
	Size    string `json:"size"`
	PvcName string `json:"pvcName"`
	Mode    string `json:"mode"`
}

type FixVolumeCommand struct {
	ProjectName
}

type GrowVolumeCommand struct {
	ProjectName
	NewSize string `json:"newSize"`
	PvName  string `json:"pvName"`
}

type NewProjectCommand struct {
	ProjectName
	Billing string `json:"billing"`
	MegaId  string `json:"megaId"`
}

type NewTestProjectCommand struct {
	ProjectName
}

type EditBillingDataCommand struct {
	ProjectName
	Billing string `json:"billing"`
}

type EditQuotasCommand struct {
	ProjectName
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type NewServiceAccountCommand struct {
	ProjectName
	ServiceAccount string `json:"serviceAccount"`
}

type FeatureToggleResponse struct {
	Gluster bool `json:"gluster"`
	DDC     bool `json:"ddc"`
}

type ApiResponse struct {
	Message string `json:"message"`
}

type DDCBilling struct {
	Rows []DDCBillingRow `json:"rows"`
	CSV  string `json:"csv"`
}

type DDCBillingRow struct {
	Sender       string `json:"sender"`
	Art          string `json:"art"`
	Project      string `json:"project"`
	Host         string `json:"host"`
	Assignment1   string `json:"assignment1"`
	Assignment2   string `json:"assignment2"`
	Assignment3   string `json:"assignment3"`
	TotalCPU     float64 `json:"totalCpu"`
	TotalMemory  float64 `json:"totalMemory"`
	TotalStorage float64 `json:"totalStorage"`
	Total        float64 `json:"total"`
}
