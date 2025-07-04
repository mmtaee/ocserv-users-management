package user

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"ocserv-service/utils"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateUserData holds the input data for creating a new user.
type CreateUserData struct {
	Username string                 `json:"username" validate:"required,min=3"`
	Password string                 `json:"password" validate:"required,min=2"`
	Group    string                 `json:"group" validate:"required"`
	Config   map[string]interface{} `json:"config"`
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

	if data.Config != nil {
		filename := filepath.Join(configUserBaseDir, data.Username)

		// Open file with create, truncate, write-only flags and permission 0640
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer file.Close()

		err = utils.ConfigWriter(file, data.Config)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, nil)
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
