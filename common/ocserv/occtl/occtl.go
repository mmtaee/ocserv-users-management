package occtl

import (
	"common/pkg"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

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
type OcservOcctl struct{}

type OcservOcctlInterface interface {
	OnlineUsers() (*[]string, error)
	OnlineSessions() (*[]OnlineUserSession, error)
	DisconnectUser(username string) (string, error)
	ReloadConfigs() (string, error)
	ShowIPBans() (*[]IPBanPoints, error)
}

const occtlExec = "/usr/bin/occtl"

func NewOcservOcctl() *OcservOcctl {
	return &OcservOcctl{}
}

// OnlineUsers returns a list of currently connected usernames.
// Executes: occtl -j show users | jq -r '.[].Username'
func (o *OcservOcctl) OnlineUsers() (*[]string, error) {
	var users []string

	command := "-j show users | jq -r '.[].Username'"
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", occtlExec, command))
	result, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	for _, user := range strings.Split(string(result), "\n") {
		trimmed := strings.TrimSpace(user)
		if trimmed != "" {
			users = append(users, trimmed)
		}
	}

	return &users, nil
}

// OnlineSessions returns a list of currently connected user info.
// Executes: occtl -j show users
func (o *OcservOcctl) OnlineSessions() (*[]OnlineUserSession, error) {
	command := "-j show users"
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", occtlExec, command))
	result, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sessions []OnlineUserSession
	if err = json.Unmarshal(result, &sessions); err != nil {
		return nil, err
	}
	return &sessions, nil
}

// DisconnectUser disconnects the given user.
// Executes: occtl disconnect user <username>
func (o *OcservOcctl) DisconnectUser(username string) (string, error) {
	cmd := exec.Command(occtlExec, "disconnect", "user", username)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ReloadConfigs reloads the ocserv configuration.
// Executes: occtl reload
func (o *OcservOcctl) ReloadConfigs() (string, error) {
	cmd := exec.Command(occtlExec, "reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ShowIPBans returns the current list of IP bans with scores.
// Executes: occtl -j show ip bans points
func (o *OcservOcctl) ShowIPBans() (*[]IPBanPoints, error) {
	cmd := exec.Command(occtlExec, "-j", "show", "ip", "bans", "points")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var ipBans []IPBanPoints
	if output := string(out); output != "" {
		if err = json.Unmarshal(out, &ipBans); err != nil {
			return nil, err
		}
	}

	return &ipBans, nil

}

// UnbanIP removes an IP ban from the given IP address.
// Executes: occtl unban ip <ip>
func (o *OcservOcctl) UnbanIP(ip string) (string, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "", fmt.Errorf("invalid IP: %s", ip)
	}

	cmd := exec.Command(occtlExec, "unban", "ip", ip)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ShowStatus returns the current status of ocserv.
// Executes: occtl -j show status
// Executes: occtl show status
func (o *OcservOcctl) ShowStatus(raw bool) (interface{}, error) {
	cmd := exec.Command(occtlExec, "-j", "show", "status")
	if raw {
		cmd = exec.Command(occtlExec, "show", "status")
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	if raw {
		return string(out), nil
	}

	var status map[string]interface{}
	if err = json.Unmarshal(out, &status); err != nil {
		return nil, err
	}
	return status, nil
}

// ShowIRoutes returns the current iRoutes information.
// Executes: occtl -j show iroutes
func (o *OcservOcctl) ShowIRoutes() (*[]IRoute, error) {
	var routes []IRoute
	version := pkg.GetOcservVersion()
	if version == "1.2.4" { // has bug on IRoute Command
		return &routes, nil
	}
	cmd := exec.Command(occtlExec, "-j", "show", "iroutes")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &routes); err != nil {
		return nil, err
	}
	return &routes, nil
}

func (o *OcservOcctl) ShowUser(username string) (OnlineUserSession, error) {
	var session OnlineUserSession

	cmd := exec.Command(occtlExec, "-j", "show", "user", username)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return session, err
	}
	if err = json.Unmarshal(out, &session); err != nil {
		return session, err
	}
	return session, nil
}
