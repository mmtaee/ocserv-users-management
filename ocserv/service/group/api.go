package group

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Controller struct{}

const (
	ocpasswdPath       = "/etc/ocserv/ocpasswd"
	ocpasswdExec       = "/usr/bin/ocpasswd"
	configGroupBaseDir = "/etc/ocserv/groups/"
	defaultGroupFile   = "/etc/ocserv/defaults/group.conf"
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
	type groupRequest struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	}
	var req groupRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request: "+err.Error())
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group name is required")
	}

	filename := filepath.Join(configGroupBaseDir, req.Name)

	// Open file with create, truncate, write-only flags and permission 0640
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	err = groupWriter(file, req.Config)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"group": req.Name, "filename": filename})
}

// Delete deletes the group config file specified by the ":group" path parameter.
//
// After deleting the config file, it finds all users assigned to this group
// in the ocpasswd file and resets their group assignment to default (empty).
//
// Returns HTTP 200 with a summary of updated users on success.
func (ctrl *Controller) Delete(c echo.Context) error {
	group := c.Param("name")
	log.Println("group name in param:", group)
	if group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group name is required")
	}

	// Delete group config file
	filename := filepath.Join(configGroupBaseDir, group)
	log.Println("filename: ", filename)
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

// GetDefaultsGroup returns defaults group config
//
// Returns HTTP 200 with defaults group config.
func (ctrl *Controller) GetDefaultsGroup(c echo.Context) error {
	config, err := parseOcservConfigFile(defaultGroupFile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, config)
}

// UpdateDefaultsGroup update defaults group config
//
// Returns HTTP 200 without response data.
func (ctrl *Controller) UpdateDefaultsGroup(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request: "+err.Error())
	}

	file, err := os.OpenFile(defaultGroupFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	err = groupWriter(file, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}
