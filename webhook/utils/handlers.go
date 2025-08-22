package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	ocpasswdPath       = "/etc/ocserv/ocpasswd"
	ocpasswdExec       = "/usr/bin/ocpasswd"
	configGroupBaseDir = "/etc/ocserv/groups/"
	defaultGroupFile   = "/etc/ocserv/defaults/group.conf"
)

var listKeys = map[string]bool{
	"dns":       true,
	"no-route":  true,
	"route":     true,
	"split-dns": true,
}

// ConfigWriter a method to write configs in group config file
func ConfigWriter(file *os.File, config map[string]interface{}) error {
	for k, v := range config {
		if b, ok := v.(bool); ok && !b {
			continue
		}
		if v == nil {
			continue
		}

		if k == "dns" {
			for _, dns := range v.([]interface{}) {
				if _, err := file.WriteString(fmt.Sprintf("dns=%s\n", dns)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
			continue
		} else if k == "route" {
			for _, route := range v.([]interface{}) {
				if _, err := file.WriteString(fmt.Sprintf("route=%s\n", route)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
			continue
		} else if k == "no-route" {
			for _, route := range v.([]interface{}) {
				if _, err := file.WriteString(fmt.Sprintf("no-route=%s\n", route)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
			continue
		} else if k == "split-dns" {
			for _, dns := range v.([]interface{}) {
				if _, err := file.WriteString(fmt.Sprintf("split-dns=%s\n", dns)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
			continue
		} else {
			if _, err := file.WriteString(fmt.Sprintf("%s=%v\n", k, v)); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}

// GetUsersByGroup parses the ocpasswd file and returns a slice of usernames
// that belong to the specified group.
//
// It reads each line of the ocpasswd file, ignoring comments and malformed lines.
// Assumes that the group is stored as the third colon-separated field.
//
// Returns an error if reading the file or scanning fails.
func GetUsersByGroup(groupName string) ([]string, error) {
	file, err := os.Open(ocpasswdPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue // skip malformed lines
		}

		username := parts[0]
		group := parts[2]

		if group == groupName {
			users = append(users, username)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func ParseOcservConfigFile(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]interface{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		rawValue := strings.TrimSpace(parts[1])

		if listKeys[key] {
			if existing, exists := config[key]; exists {
				config[key] = append(existing.([]string), rawValue)
			} else {
				config[key] = []string{rawValue}
			}
			continue
		}

		parsedValue := ParseTypedValue(rawValue)

		if existing, exists := config[key]; exists {
			switch v := existing.(type) {
			case []interface{}:
				config[key] = append(v, parsedValue)
			default:
				config[key] = []interface{}{v, parsedValue}
			}
		} else {
			config[key] = parsedValue
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseTypedValue(s string) interface{} {
	s = strings.TrimSpace(s)

	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

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
	// Regex to find the version number
	re := regexp.MustCompile(`OpenConnect VPN Server\s+([0-9]+\.[0-9]+\.[0-9]+)`)
	match := re.FindStringSubmatch(fullOutput)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func GetOCCTLVersion() string {
	cmd := exec.Command("occtl", "--version")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Command error:", err)
		return ""
	}

	fullOutput := out.String() + stderr.String()
	lines := strings.Split(fullOutput, "\n")
	var resultLines []string
	for _, line := range lines {
		if trim := strings.TrimSpace(line); trim != "" {
			if strings.HasPrefix(trim, "Copyright") {
				break
			}
			resultLines = append(resultLines, trim)
		}
	}

	if len(resultLines) == 0 {
		return ""
	}

	finalOutput := strings.Join(resultLines, "\n")
	return finalOutput
}

func fixTrailingComma(jsonBytes []byte) []byte {
	re := regexp.MustCompile(`("in_use"\s*:\s*\d+)\s*,`)
	return re.ReplaceAll(jsonBytes, []byte("$1"))
}

func OcctlResponse(c echo.Context, cmd *exec.Cmd, resultType interface{}) error {
	out, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to info: "+err.Error())
	}

	var info interface{}

	switch resultType.(type) {
	case string:
		return c.String(http.StatusOK, string(out))
	case map[string]interface{}:
		info = map[string]interface{}{}
	case []interface{}:
		info = []interface{}{}
	default:
		info = new(interface{})
	}

	out = fixTrailingComma(out)

	if err = json.Unmarshal(out, &info); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse info JSON: "+err.Error())
	}
	return c.JSON(http.StatusOK, info)
}
