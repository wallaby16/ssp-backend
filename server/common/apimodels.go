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

type ApiResponse struct {
	Message string `json:"message"`
}
