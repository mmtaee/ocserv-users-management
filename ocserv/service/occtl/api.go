package occtl

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

type Controller struct{}

type IPBanPoints struct {
	IP    string `json:"IP"`
	Since string `json:"Since"`
	Until string `json:"_Since"`
	Score int    `json:"Score"`
}

type IRoute struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Vhost    string `json:"vhost"`
	Device   string `json:"Device"`
	IP       string `json:"IP"`
	IRoute   string `json:"iRoutes"`
}

const occtlExec = "/usr/bin/occtl"

func New() *Controller {
	return &Controller{}
}

// OnlineUsers returns a list of currently connected usernames.
//
// Executes: occtl -j show users | jq -r '.[].Username'
func (ctrl *Controller) OnlineUsers(c echo.Context) error {
	var users []string

	command := "-j show users | jq -r '.[].Username'"
	cmd := exec.Command("sh", "-c", fmt.Sprintf("/usr/bin/occtl %s", command))
	result, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, line := range strings.Split(string(result), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			users = append(users, trimmed)
		}
	}

	return c.JSON(http.StatusOK, map[string][]string{"users": users})

}

// DisconnectUser disconnects the given user.
//
// Executes: occtl disconnect user <username>
func (ctrl *Controller) DisconnectUser(c echo.Context) error {
	username := c.Param("username")

	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	cmd := exec.Command(occtlExec, "disconnect", "user", username)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to disconnect user: "+string(out))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":  "user disconnected",
		"username": username,
		"output":   string(out),
	})
}

// Reload reloads the ocserv configuration.
//
// Executes: occtl reload
func (ctrl *Controller) Reload(c echo.Context) error {
	cmd := exec.Command(occtlExec, "reload")
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to reload ocserv: "+output)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ocserv reloaded successfully",
		"output":  output,
	})
}

// ShowIPBans returns the current list of IP bans with scores.
//
// Executes: occtl -j show ip bans points
func (ctrl *Controller) ShowIPBans(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "ip", "bans", "points")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get IP bans: "+string(out))
	}

	var ipBans []IPBanPoints
	if err = json.Unmarshal(out, &ipBans); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse JSON: "+err.Error())
	}

	return c.JSON(http.StatusOK, ipBans)
}

// UnbanIP removes an IP ban from the given IP address.
//
// Executes: occtl unban ip <ip>
func (ctrl *Controller) UnbanIP(c echo.Context) error {
	ip := c.Param("ip")
	if ip == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ip parameter is required")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid IP address")
	}

	cmd := exec.Command(occtlExec, "unban", "ip", ip)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		if output == "" {
			output = err.Error()
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": output})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "IP unbanned", "output": output})
}

// ShowStatus returns the current status of ocserv.
//
// Executes: occtl -j show status
func (ctrl *Controller) ShowStatus(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "status")
	out, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get status: "+err.Error())
	}

	var status map[string]interface{}
	if err = json.Unmarshal(out, &status); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse status JSON: "+err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

// ShowIRoutes returns the current iRoutes information.
//
// Executes: occtl -j show iroutes
func (ctrl *Controller) ShowIRoutes(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "iroutes")
	out, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get iroutes: "+err.Error())
	}

	var iRoutes []IRoute
	if err = json.Unmarshal(out, &iRoutes); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse iroutes JSON: "+err.Error())
	}

	return c.JSON(http.StatusOK, iRoutes)
}

// ShowUser returns detailed information about a specific user.
//
// Executes: occtl -j show user <username>
func (ctrl *Controller) ShowUser(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	cmd := exec.Command(occtlExec, "-j", "show", "user", username)
	out, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info: "+err.Error())
	}

	var userInfo map[string]interface{}
	if err = json.Unmarshal(out, &userInfo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse user info JSON: "+err.Error())
	}

	return c.JSON(http.StatusOK, userInfo)
}
