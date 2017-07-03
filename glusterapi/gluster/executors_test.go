package gluster

import (
	"testing"
)

func init() {
	ExecRunner = TestRunner{}
}

func TestExecuteCommandsLocally(t *testing.T) {
	executeCommandsLocally([]string{"test", "test2"})

	equals(t, "bash -c test", commands[0])
	equals(t, "bash -c test2", commands[1])
}
