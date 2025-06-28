package occtl

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
)

func GetOcservVersion() string {
	cmd := exec.Command("ocserv", "--version")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Command error:", err)
		return ""
	}

	// Combine stdout and stderr for pattern matching
	fullOutput := out.String() + stderr.String()
	log.Println("Raw output:\n" + fullOutput)

	// Regex to find the version number
	re := regexp.MustCompile(`OpenConnect VPN Server\s+([0-9]+\.[0-9]+\.[0-9]+)`)
	match := re.FindStringSubmatch(fullOutput)
	if len(match) >= 2 {
		log.Println("Version:", match[1])
		return match[1]
	}
	log.Println("Version not found")
	return ""
}
