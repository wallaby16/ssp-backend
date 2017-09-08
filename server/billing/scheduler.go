package billing

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


