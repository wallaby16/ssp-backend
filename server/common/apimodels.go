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

type CreateLogseneAppCommand struct {
	AppName string `json:"appName"`
	EditSematextPlanCommand
	EditBillingDataCommand
}

type EditSematextPlanCommand struct {
	PlanId int `json:"planId"`
	Limit  int `json:"limit"`
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

type InstanceListResponse struct {
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Name                  string `json:"name"`
	InstanceId            string `json:"id"`
	State                 string `json:"state"`
	StateTransitionReason string `json:"statereason"`
	Account               string `json:"account"`
}

type S3CredentialsResponse struct {
	Username    string `json:"username"`
	AccessKeyID string `json:"accesskeyid"`
	SecretKey   string `json:"secretkey"`
}

type AdminList struct {
	Admins []string `json:"admins"`
}

type SematextAppList struct {
	AppId         int     `json:"appId"`
	Name          string  `json:"name"`
	PlanName      string  `json:"planName"`
	UserRole      string  `json:"userRole"`
	IsFree        bool    `json:"isFree"`
	PricePerMonth float64 `json:"pricePerMonth"`
	BillingInfo   string  `json:"billingInfo"`
}

type SematextLogsenePlan struct {
	PlanId                     int     `json:"planId"`
	Name                       string  `json:"name"`
	IsFree                     bool    `json:"isFree"`
	PricePerMonth              float64 `json:"pricePerMonth"`
	DefaultDailyMaxLimitSizeMb float64 `json:"defaultDailyMaxLimitSizeMb"`
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
	UserName   string `json:"username"`
	IsReadonly bool   `json:"isReadonly"`
}
