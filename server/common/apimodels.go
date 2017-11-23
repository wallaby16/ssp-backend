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

type BucketListResponse struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Name    string `json:"name"`
	Account string `json:"account"`
}

type S3CredentialsResponse struct {
	Username    string `json:"username"`
	AccessKeyID string `json:"accesskeyid"`
	SecretKey   string `json:"secretkey"`
}

type AdminList struct {
	Admins []string `json:"admins"`
}

type DDCBilling struct {
	Rows []DDCBillingRow `json:"rows"`
	CSV  string          `json:"csv"`
}

type DDCBillingRow struct {
	Sender              string  `json:"sender"`
	Art                 string  `json:"art"`
	Project             string  `json:"project"`
	Host                string  `json:"host"`
	Backup              bool    `json:"backup"`
	ReceptionAssignment string  `json:"receptionAssignment"`
	OrderReception      string  `json:"orderReception"`
	PspElement          string  `json:"pspElement"`
	TotalCPU            float64 `json:"totalCpu"`
	TotalMemory         float64 `json:"totalMemory"`
	TotalStorage        float64 `json:"totalStorage"`
	Total               float64 `json:"total"`
}

type NewS3BucketCommand struct {
	ProjectName
	BucketName string `json:"bucketname"`
	Billing    string `json:"billing"`
	Stage      string `json:"stage"`
}

type NewS3UserCommand struct {
	ProjectName
	UserName   string `json:"username"`
	Stage      string `json:"stage"`
	IsReadonly bool   `json:"isReadonly"`
}
