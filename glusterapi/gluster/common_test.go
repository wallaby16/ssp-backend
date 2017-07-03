package gluster

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

type TestRunner struct{}

var commands []string
var output string

func (r TestRunner) Run(command string, args ...string) ([]byte, error) {
	commands = append(commands, command+" "+strings.Join(args, " "))
	return []byte(output), nil
}

func init() {
	ExecRunner = TestRunner{}
}

func TestGetGlusterPeerServers(t *testing.T) {
	ip1 := "192.168.125.236"
	ip2 := "192.168.125.238"

	output = fmt.Sprintf(`Hostname: %v
						  Hostname: %v`, ip1, ip2)

	servers, _ := getGlusterPeerServers()

	if servers[0] != ip1 {
		t.Errorf("Expected %v to be %v", servers[0], ip1)
	}
	if servers[1] != ip2 {
		t.Errorf("Expected %v to be %v", servers[1], ip2)
	}
}

func TestGetLocalServersIP(t *testing.T) {
	localIP, _ := getLocalServersIP()

	// Make sure response is a valid ip
	ip := net.ParseIP(localIP)

	if ip.To4() == nil {
		t.Errorf("Expected to get local ip, but got %v", localIP)
	}
}
