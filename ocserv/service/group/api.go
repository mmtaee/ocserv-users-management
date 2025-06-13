package group

import (
	"bufio"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Controller struct{}

const (
	ocpasswdPath       = "/etc/ocserv/ocpasswd"
	ocpasswdExec       = "/usr/bin/ocpasswd"
	configGroupBaseDir = "/etc/ocserv/groups/"
)

func New() *Controller {
	return &Controller{}
}

// Create creates or updates the JSON configuration file for a specific group.
// The group name is taken from the URL path parameter ":group".
// The JSON config is read from the request body.
//
// If the config file exists, it is truncated (cleared) before writing.
// If it does not exist, it is created with permissions 0640.
//
// Returns HTTP 200 with the group name and config file path on success.
func (ctrl *Controller) Create(c echo.Context) error {
	group := c.Param("group")
	if group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group name is required")
	}

	var config map[string]interface{}
	if err := c.Bind(&config); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	filename := filepath.Join(configGroupBaseDir, group)

	// Open file with create, truncate, write-only flags and permission 0640
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"group": group, "filename": filename})
}

// Delete deletes the group config file specified by the ":group" path parameter.
//
// After deleting the config file, it finds all users assigned to this group
// in the ocpasswd file and resets their group assignment to default (empty).
//
// Returns HTTP 200 with a summary of updated users on success.
func (ctrl *Controller) Delete(c echo.Context) error {
	group := c.Param("group")
	if group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group name is required")
	}

	// Delete group config file
	filename := filepath.Join(configGroupBaseDir, group)
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "group config not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Find all users currently assigned to the group
	users, err := getUsersByGroup(group)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users by group: "+err.Error())
	}

	// Reset group assignment for each user to default (empty)
	for _, user := range users {
		out, err := exec.Command(ocpasswdExec, "-g", "", "-c", ocpasswdPath, user).CombinedOutput()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user '"+user+"': "+string(out))
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "group config deleted and users reset to default group",
		"group":         group,
		"updated_users": users,
	})
}

// ListUsers returns a JSON list of all usernames assigned to the group specified
// by the ":group" path parameter.
//
// Returns HTTP 200 with group name and slice of usernames.
func (ctrl *Controller) ListUsers(c echo.Context) error {
	group := c.Param("group")
	if group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group name is required")
	}

	users, err := getUsersByGroup(group)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"group": group,
		"users": users,
	})
}

// getUsersByGroup parses the ocpasswd file and returns a slice of usernames
// that belong to the specified group.
//
// It reads each line of the ocpasswd file, ignoring comments and malformed lines.
// Assumes that the group is stored as the third colon-separated field.
//
// Returns an error if reading the file or scanning fails.
func getUsersByGroup(groupName string) ([]string, error) {
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
