package stdlib

import (
	"testing"
)

func TestGetRemoteHostNameLinuxOrg(t *testing.T) {
	hostname, err := GetRemoteHostName("207.97.227.239")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if hostname != "github.com" {
		t.Errorf("Expected %s, got %s", "github.com", hostname)
	}
}

func TestGetRemoteHostNameBroadcasthost(t *testing.T) {
	hostname, err := GetRemoteHostName("255.255.255.255")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if hostname != "broadcasthost" {
		t.Errorf("Expected %s, got %s", "broadcasthost", hostname)
	}
}

func TestGetRemoteHostNameUknownIP(t *testing.T) {
	hostname, err := GetRemoteHostName("1.1.1.1")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if hostname != "1.1.1.1" {
		t.Errorf("Expected %s, got %s", "1.1.1.1", hostname)
	}
}

func TestGetRemoteHostNameUknownHostName(t *testing.T) {
	_, err := GetRemoteHostName("")
	if err == nil {
		t.Fatalf("Expected error: Error: nodename nor servname provided, or not known")
	}
}
