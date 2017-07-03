package gluster

import (
	"fmt"

	"testing"
)

func init() {
	ExecRunner = TestRunner{}
}

func TestExecuteCommandsLocally(t *testing.T) {
	executeCommandsLocally([]string{"test", "test2"})

	if commands[0] != "bash -c test" {
		fmt.Errorf("Expected %v but was %v", "bash -c test", commands[0])
	}
	if commands[1] != "bash -c test2" {
		fmt.Errorf("Expected %v but was %v", "bash -c test2", commands[1])
	}
}
