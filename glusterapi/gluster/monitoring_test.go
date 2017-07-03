package gluster

import "testing"

func init() {
	ExecRunner = TestRunner{}
}

func TestGetVolumeUsage(t *testing.T) {
	output = []string{"    49664    2864 /dev/mapper/vg_mylv_project_pv1"}

	volInfo, _ := getVolumeUsage("pv1")

	equals(t, 49664, volInfo.TotalKiloBytes)
	equals(t, 2864, volInfo.UsedKiloBytes)
}

func TestCheckVolumeUsage_OK(t *testing.T) {
	output = []string{"    49664    2864 /dev/mapper/vg_mylv_project_pv1"}

	err := checkVolumeUsage("pv1", "20")
	ok(t, err)
}

func TestCheckVolumeUsage_Error(t *testing.T) {
	output = []string{"    49664    49555 /dev/mapper/vg_mylv_project_pv1"}

	err := checkVolumeUsage("pv1", "20")
	assert(t, err != nil, "Should return error as bigger than threshold")
}