// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/mattermost-server/model"
)

func (a *App) GetGroup(id string) (*model.Group, *model.AppError) {
	result := <-a.Srv.Store.Group().Get(id)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Group), nil
}

func (a *App) GetGroupsPage(page int, perPage int) ([]*model.Group, *model.AppError) {
	result := <-a.Srv.Store.Group().GetAllPage(page*perPage, perPage)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.Group), nil
}

func (a *App) CreateGroup(group *model.Group) (*model.Group, *model.AppError) {
	// Clear any user-provided values for trusted properties.
	group.Id = ""
	group.CreateAt = 0
	group.UpdateAt = 0
	group.DeleteAt = 0
	result := <-a.Srv.Store.Group().Save(group)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Group), nil
}

func (a *App) UpdateGroup(group *model.Group) (*model.Group, *model.AppError) {
	result := <-a.Srv.Store.Group().Save(group)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Group), nil
}

func (a *App) DeleteGroup(groupID string) (*model.Group, *model.AppError) {
	result := <-a.Srv.Store.Group().Delete(groupID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Group), nil
}

func (a *App) CreateGroupMember(groupID string, userID string) (*model.GroupMember, *model.AppError) {
	result := <-a.Srv.Store.Group().CreateMember(groupID, userID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupMember), nil
}

func (a *App) DeleteGroupMember(groupID string, userID string) (*model.GroupMember, *model.AppError) {
	result := <-a.Srv.Store.Group().DeleteMember(groupID, userID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupMember), nil
}

func (a *App) GetGroupTeam(groupID string, teamID string) (*model.GroupTeam, *model.AppError) {
	result := <-a.Srv.Store.Group().GetGroupTeam(groupID, teamID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupTeam), nil
}

func (a *App) GetGroupTeamsPage(groupID string, page int, perPage int) ([]*model.GroupTeam, *model.AppError) {
	result := <-a.Srv.Store.Group().GetAllGroupTeamsByGroupPage(groupID, page*perPage, perPage)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.GroupTeam), nil
}

func (a *App) CreateGroupTeam(groupTeam *model.GroupTeam) (*model.GroupTeam, *model.AppError) {
	result := <-a.Srv.Store.Group().SaveGroupTeam(groupTeam)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupTeam), nil
}

func (a *App) DeleteGroupTeam(groupID string, teamID string) (*model.GroupTeam, *model.AppError) {
	result := <-a.Srv.Store.Group().DeleteGroupTeam(groupID, teamID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupTeam), nil
}

func (a *App) GetGroupChannel(groupID string, channelID string) (*model.GroupChannel, *model.AppError) {
	result := <-a.Srv.Store.Group().GetGroupChannel(groupID, channelID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupChannel), nil
}

func (a *App) GetGroupChannelsPage(groupID string, page int, perPage int) ([]*model.GroupChannel, *model.AppError) {
	result := <-a.Srv.Store.Group().GetAllGroupChannelsByGroupPage(groupID, page*perPage, perPage)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.GroupChannel), nil
}

func (a *App) CreateGroupChannel(groupChannel *model.GroupChannel) (*model.GroupChannel, *model.AppError) {
	result := <-a.Srv.Store.Group().SaveGroupChannel(groupChannel)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupChannel), nil
}

func (a *App) DeleteGroupChannel(groupID string, channelID string) (*model.GroupChannel, *model.AppError) {
	result := <-a.Srv.Store.Group().DeleteGroupChannel(groupID, channelID)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.GroupChannel), nil
}
