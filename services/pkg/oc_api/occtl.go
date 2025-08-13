package oc_api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type OcctlApiRepository struct {
	url string
}

type OcOcctlApiRepositoryInterface interface {
	ShowEvent(ctx context.Context) ([]byte, error)
	UnbanIP(ctx context.Context, ip string) ([]byte, error)
	ShowSessionBySID(ctx context.Context, sid string) ([]byte, error)
	ShowSessionsAllValid(ctx context.Context) ([]byte, error)
	ShowSessionsAll(ctx context.Context) ([]byte, error)
	IPBans(ctx context.Context) ([]byte, error)
	IRoutes(ctx context.Context) ([]byte, error)
	Version(ctx context.Context) ([]byte, error)
	ShowUserByID(ctx context.Context, id string) ([]byte, error)
	OnlineUsersInfo(ctx context.Context) ([]byte, error)
	ShowUserByUsername(ctx context.Context, username string) ([]byte, error)
	Status(ctx context.Context) ([]byte, error)
	OnlineUsers(ctx context.Context) ([]byte, error)
	Disconnect(ctx context.Context, username string) ([]byte, error)
	Reload(ctx context.Context) error
}

func NewOcctlApiRepository(url string) *OcctlApiRepository {
	return &OcctlApiRepository{url: url}
}

func (o *OcctlApiRepository) ShowEvent(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/events/", http.MethodGet)
}

func (o *OcctlApiRepository) UnbanIP(ctx context.Context, ip string) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/unban-ip/"+ip, http.MethodDelete)
}

func (o *OcctlApiRepository) ShowSessionBySID(ctx context.Context, sid string) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/sessions/sid/"+sid, http.MethodGet)
}

func (o *OcctlApiRepository) ShowSessionsAllValid(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/sessions/valid", http.MethodGet)
}

func (o *OcctlApiRepository) ShowSessionsAll(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/sessions/", http.MethodGet)
}

func (o *OcctlApiRepository) ShowUserByUsername(ctx context.Context, username string) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/users/"+username, http.MethodGet)
}

func (o *OcctlApiRepository) ShowUserByID(ctx context.Context, id string) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/users/id/"+id, http.MethodGet)
}

func (o *OcctlApiRepository) Version(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/version", http.MethodGet)
}

func (o *OcctlApiRepository) OnlineUsersInfo(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/online-users/info", http.MethodGet)
}

func (o *OcctlApiRepository) IPBans(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/ip-bans", http.MethodGet)
}

func (o *OcctlApiRepository) IRoutes(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/iroutes", http.MethodGet)
}

func (o *OcctlApiRepository) Disconnect(ctx context.Context, username string) ([]byte, error) {
	return o.doRequestBytes(ctx, fmt.Sprintf("%s/api/occtl/disconnect/%s", o.url, username), http.MethodPost)
}

func (o *OcctlApiRepository) Reload(ctx context.Context) error {
	_, err := o.doRequestBytes(ctx, o.url+"/api/occtl/reload", http.MethodPost)
	return err
}

func (o *OcctlApiRepository) Status(ctx context.Context) ([]byte, error) {
	url := o.url + "/api/occtl/status"
	
	format, _ := ctx.Value("format").(string)
	if format == "json" {
		if strings.Contains(url, "?") {
			url += "&format=json"
		} else {
			url += "?format=json"
		}
	}
	return o.doRequestBytes(ctx, url, http.MethodGet)
}

func (o *OcctlApiRepository) OnlineUsers(ctx context.Context) ([]byte, error) {
	return o.doRequestBytes(ctx, o.url+"/api/occtl/online-users", http.MethodGet)
}
