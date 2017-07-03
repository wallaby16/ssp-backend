package gluster

import (
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGrowVolume_Empty(t *testing.T) {
	err := growVolume("", "")

	if err == nil {
		t.Error("growVolume should throw error if called empty")
	}
}

func TestGrowVolume_WrongSize(t *testing.T) {
	err := growVolume("pv", "101G")

	if err == nil {
		t.Error("growVolume should throw error if called with wrong size")
	}
}

func TestGrowVolume_WrongSizeMB(t *testing.T) {
	err := growVolume("pv", "1025M")

	if err == nil {
		t.Error("growVolume should throw error if called with wrong size")
	}
}

func TestGrowVolume(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://192.168.125.236:0/sec/lv/grow",
		httpmock.NewStringResponder(200, ""))

	commands = nil
	output = "Hostname: 192.168.125.236"
	VgName = "myvg"

	growVolume("pv", "10M")

	// Should call the remote server
	if httpmock.GetTotalCallCount() != 1 {
		t.Errorf("Should have called the remote gluster server")
	}

	// Should execute commands locally
	if commands[1] != "bash -c lvextend -L 10M /dev/myvg/lv_pv" {
		t.Errorf("Command was wrong: %v", commands[1])
	}
	if commands[2] != "bash -c xfs_growfs /dev/myvg/lv_pv" {
		t.Errorf("Command was wrong: %v", commands[2])
	}
}