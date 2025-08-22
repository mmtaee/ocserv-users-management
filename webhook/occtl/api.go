package occtl

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"ocserv-service/utils"
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

type OnlineUserSession struct {
	Username    string `json:"Username"`
	Group       string `json:"Groupname"`
	AverageRX   string `json:"Average RX"`
	AverageTX   string `json:"Average TX"`
	ConnectedAt string `json:"_Connected at"`
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

// OnlineUsersInfo returns a list of currently connected user info.
//
// Executes: occtl -j show users
func (ctrl *Controller) OnlineUsersInfo(c echo.Context) error {
	var users []OnlineUserSession

	command := "-j show users"
	cmd := exec.Command("sh", "-c", fmt.Sprintf("/usr/bin/occtl %s", command))
	result, err := cmd.Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = json.Unmarshal(result, &users); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse occtl output: "+err.Error())
	}
	return c.JSON(http.StatusOK, &users)
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
	if output := string(out); output != "" {
		if err = json.Unmarshal(out, &ipBans); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse JSON: "+err.Error())
		}
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
	var cmd *exec.Cmd
	format := c.QueryParam("format")

	if format == "json" {
		cmd = exec.Command(occtlExec, "-j", "show", "status")
	} else {
		cmd = exec.Command(occtlExec, "show", "status")
	}

	if format == "json" {
		return utils.OcctlResponse(c, cmd, map[string]interface{}{})
	}
	return utils.OcctlResponse(c, cmd, "")
}

// ShowIRoutes returns the current iRoutes information.
//
// Executes: occtl -j show iroutes
func (ctrl *Controller) ShowIRoutes(c echo.Context) error {
	var iRoutes []IRoute

	version := utils.GetOcservVersion()
	if version == "1.2.4" { // has bug on IRoute Command
		return c.JSON(http.StatusOK, iRoutes)
	}

	cmd := exec.Command(occtlExec, "-j", "show", "iroutes", "")
	return utils.OcctlResponse(c, cmd, []interface{}{})
}

// ShowUser returns detailed information about a specific user by username.
//
// Executes: occtl -j show user <username>
func (ctrl *Controller) ShowUser(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	cmd := exec.Command(occtlExec, "-j", "show", "user", username)
	return utils.OcctlResponse(c, cmd, map[string]interface{}{})
}

// Version returns detailed information about ocserv version.
//
// Executes: ocserv -v
func (ctrl *Controller) Version(c echo.Context) error {
	version := utils.GetOcservVersion()
	occtlVersion := utils.GetOCCTLVersion()

	return c.JSON(http.StatusOK, map[string]string{
		"version":       version,
		"occtl_version": occtlVersion,
	})
}

// ShowUserByID returns detailed information about a specific user by ID.
//
// Executes: occtl -j show id <id>
func (ctrl *Controller) ShowUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	cmd := exec.Command(occtlExec, "-j", "show", "id", id)
	return utils.OcctlResponse(c, cmd, map[string]interface{}{})
}

// ShowSession returns detailed information about a specific session by SID.
//
// Executes: occtl -j show session <SID>
func (ctrl *Controller) ShowSession(c echo.Context) error {
	sid := c.Param("sid")
	if sid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "SID is required")
	}

	cmd := exec.Command(occtlExec, "-j", "show", "session", sid)
	return utils.OcctlResponse(c, cmd, map[string]interface{}{})
}

// ShowSessionsALL returns detailed information about all sessions.
//
// Executes: occtl -j show sessions all
func (ctrl *Controller) ShowSessionsALL(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "sessions", "all")
	return utils.OcctlResponse(c, cmd, []interface{}{})
}

// ShowSessionsValid returns detailed information  about all valid sessions.
//
// Executes: occtl -j show sessions valid
func (ctrl *Controller) ShowSessionsValid(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "sessions", "valid")
	return utils.OcctlResponse(c, cmd, []interface{}{})
}

// ShowEvent returns detailed information about events.
//
// Executes: occtl -j show events
func (ctrl *Controller) ShowEvent(c echo.Context) error {
	cmd := exec.Command(occtlExec, "-j", "show", "events")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		_, _ = stdin.Write([]byte("q\n"))
	}()

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	cleaned := strings.Replace(string(output), "Press 'q' or CTRL+C to quit", "", 1)
	cleaned = strings.TrimSpace(cleaned)

	return c.String(http.StatusOK, cleaned)
}
