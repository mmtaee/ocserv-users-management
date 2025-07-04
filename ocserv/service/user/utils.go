package user

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os/exec"
	"path/filepath"
)

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
