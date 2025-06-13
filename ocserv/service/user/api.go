package user

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateUserData holds the input data for creating a new user.
type CreateUserData struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=2"`
	Group    string `json:"group" validate:"required"`
}

type Controller struct{}

const (
	ocpasswdPath      = "/etc/ocserv/ocpasswd"
	ocpasswdExec      = "/usr/bin/ocpasswd"
	configUserBaseDir = "/etc/ocserv/users/"
)

func New() *Controller {
	return &Controller{}
}

// validateUsername returns an error if the username is empty.
func validateUsername(username string) error {
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}
	return nil
}

// runOcpasswd runs the ocpasswd command with the given args and returns combined output and error.
func runOcpasswd(args ...string) (string, error) {
	cmd := exec.Command(ocpasswdExec, args...)
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		if output == "" {
			output = err.Error()
		}
		return output, err
	}
	return output, nil
}

// configFilePath returns the file path for a user's config file.
func configFilePath(username string) string {
	return filepath.Join(configUserBaseDir, username)
}

// Create handles the creation of a new user with username, password, and group.
func (ctrl *Controller) Create(c echo.Context) error {
	data := new(CreateUserData)

	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := []string{"-c", ocpasswdPath, data.Username}
	if data.Group != "" {
		args = append([]string{"-g", data.Group}, args...)
	}

	cmd := exec.Command(ocpasswdExec, args...)
	cmd.Stdin = bytes.NewBufferString(data.Password + "\n" + data.Password + "\n")

	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		if output == "" {
			output = err.Error()
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": output})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user added", "output": output})
}

// Lock locks a user account specified by the username path parameter.
func (ctrl *Controller) Lock(c echo.Context) error {
	username := c.Param("username")
	if err := validateUsername(username); err != nil {
		return err
	}

	output, err := runOcpasswd("-l", "-c", ocpasswdPath, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": output})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user locked", "output": output})
}

// Unlock unlocks a user account specified by the username path parameter.
func (ctrl *Controller) Unlock(c echo.Context) error {
	username := c.Param("username")
	if err := validateUsername(username); err != nil {
		return err
	}

	output, err := runOcpasswd("-u", "-c", ocpasswdPath, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": output})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user unlocked", "output": output})
}

// Delete deletes a user account specified by the username path parameter.
func (ctrl *Controller) Delete(c echo.Context) error {
	username := c.Param("username")
	if err := validateUsername(username); err != nil {
		return err
	}

	output, err := runOcpasswd("-d", "-c", ocpasswdPath, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": output})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted", "output": output})
}

// CreateConfig creates or overwrites the JSON config file for the specified user.
func (ctrl *Controller) CreateConfig(c echo.Context) error {
	username := c.Param("username")
	if err := validateUsername(username); err != nil {
		return err
	}

	var config map[string]interface{}
	if err := c.Bind(&config); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	filename := configFilePath(username)
	if err = os.WriteFile(filename, jsonBytes, 0640); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"username": username, "filename": filename})
}

// DeleteConfig deletes the config file for the specified user.
func (ctrl *Controller) DeleteConfig(c echo.Context) error {
	username := c.Param("username")
	if err := validateUsername(username); err != nil {
		return err
	}

	if err := os.Remove(configFilePath(username)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
