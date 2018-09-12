// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
)

func (api *API) InitGroup() {
	api.BaseRoutes.Groups.Handle("", api.ApiSessionRequired(createGroup)).Methods("POST")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}", api.ApiSessionRequiredTrustRequester(getGroup)).Methods("GET")
	api.BaseRoutes.Groups.Handle("", api.ApiSessionRequired(getGroups)).Methods("GET")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}", api.ApiSessionRequired(updateGroup)).Methods("PUT")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}", api.ApiSessionRequired(deleteGroup)).Methods("DELETE")

	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/members", api.ApiSessionRequired(createGroupMember)).Methods("POST")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/members/{user_id:[A-Za-z0-9]+}", api.ApiSessionRequired(deleteGroupMember)).Methods("DELETE")

	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/teams", api.ApiSessionRequired(createGroupTeam)).Methods("POST")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/teams", api.ApiSessionRequired(getGroupTeams)).Methods("GET")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/teams/{team_id:[A-Za-z0-9]+}", api.ApiSessionRequired(getGroupTeam)).Methods("GET")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/teams/{team_id:[A-Za-z0-9]+}/patch", api.ApiSessionRequired(patchGroupTeam)).Methods("PUT")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/teams/{team_id:[A-Za-z0-9]+}", api.ApiSessionRequired(deleteGroupTeam)).Methods("DELETE")

	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/channels", api.ApiSessionRequired(createGroupChannel)).Methods("POST")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/channels", api.ApiSessionRequired(getGroupChannels)).Methods("GET")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}", api.ApiSessionRequired(getGroupChannel)).Methods("GET")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}/patch", api.ApiSessionRequired(patchGroupChannel)).Methods("PUT")
	api.BaseRoutes.Groups.Handle("/{group_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}", api.ApiSessionRequired(deleteGroupChannel)).Methods("DELETE")
}

func createGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	group := model.GroupFromJson(r.Body)
	if group == nil {
		c.SetInvalidParam("group")
		return
	}

	if c.App.License() == nil || !*c.App.License().Features.LDAP {
		c.Err = model.NewAppError("Api4.CreateGroup", "api.group.create_group.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	group, err := c.App.CreateGroup(group)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(group.ToJson()))
}

func getGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireGroupId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	group, err := c.App.GetGroup(c.Params.GroupId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(group.ToJson()))
}

func getGroups(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	groups, err := c.App.GetGroupsPage(c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.GroupsToJson(groups)))
}

func updateGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireGroupId()
	if c.Err != nil {
		return
	}

	update := model.GroupFromJson(r.Body)
	if update == nil {
		c.SetInvalidParam("group")
		return
	}

	if c.App.License() == nil || !*c.App.License().Features.LDAP {
		c.Err = model.NewAppError("Api4.UpdateGroup", "api.group.update_group.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	update.Id = c.Params.GroupId

	group, err := c.App.UpdateGroup(update)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAudit("")
	w.Write([]byte(group.ToJson()))
}

func deleteGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireGroupId()
	if c.Err != nil {
		return
	}

	if c.App.License() == nil || !*c.App.License().Features.LDAP {
		c.Err = model.NewAppError("Api4.DeleteGroup", "api.group.delete_group.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if _, err := c.App.DeleteGroup(c.Params.GroupId); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func createGroupMember(c *Context, w http.ResponseWriter, r *http.Request) {}
func deleteGroupMember(c *Context, w http.ResponseWriter, r *http.Request) {}

func createGroupTeam(c *Context, w http.ResponseWriter, r *http.Request) {}
func getGroupTeams(c *Context, w http.ResponseWriter, r *http.Request)   {}
func getGroupTeam(c *Context, w http.ResponseWriter, r *http.Request)    {}
func patchGroupTeam(c *Context, w http.ResponseWriter, r *http.Request)  {}
func deleteGroupTeam(c *Context, w http.ResponseWriter, r *http.Request) {}

func createGroupChannel(c *Context, w http.ResponseWriter, r *http.Request) {}
func getGroupChannels(c *Context, w http.ResponseWriter, r *http.Request)   {}
func getGroupChannel(c *Context, w http.ResponseWriter, r *http.Request)    {}
func patchGroupChannel(c *Context, w http.ResponseWriter, r *http.Request)  {}
func deleteGroupChannel(c *Context, w http.ResponseWriter, r *http.Request) {}
