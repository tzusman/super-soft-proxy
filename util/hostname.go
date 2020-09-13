package util

import (
	"os/exec"
	"strings"
)

func GetHostname() (string, error) {
	hostnameRaw, err := exec.Command("hostname").Output()
	if err != nil {
		return "", err
	}
	hostname := strings.ToLower(strings.TrimSpace(string(hostnameRaw)))
	return hostname, nil
}
