// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
)

func TestGetGroup(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()

	if _, err := th.App.GetGroup(group.Id); err != nil {
		t.Log(err)
		t.Fatal("Should get the group")
	}

	if _, err := th.App.GetGroup(model.NewId()); err == nil {
		t.Fatal("Should not have found a group")
	}
}

func TestGetGroupsPage(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	th.CreateGroup()
	th.CreateGroup()
	th.CreateGroup()

	groups, err := th.App.GetGroupsPage(1, 2)
	if err != nil {
		t.Log(err)
		t.Fatal("Should have groups")
	}

	if len(groups) < 1 {
		t.Fatal("Should have retrieved at least one group")
	}

	if groups, _ = th.App.GetGroupsPage(999, 1); len(groups) > 0 {
		t.Fatal("Should not have groups.")
	}
}

func TestCreateGroup(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	id := model.NewId()
	group := &model.Group{
		DisplayName: "dn_" + id,
		Name:        "name" + id,
		Type:        model.GroupTypeLdap,
	}

	if _, err := th.App.CreateGroup(group); err != nil {
		t.Log(err)
		t.Fatal("Should create a new group")
	}

	if _, err := th.App.CreateGroup(group); err == nil {
		t.Fatal("Should not create a new group - group already exist")
	}
}

func TestUpdateGroup(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	group.DisplayName = model.NewId()

	if _, err := th.App.UpdateGroup(group); err != nil {
		t.Log(err)
		t.Fatal("Should update the group")
	}
}

func TestDeleteGroup(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()

	if _, err := th.App.DeleteGroup(group.Id); err != nil {
		t.Log(err)
		t.Fatal("Should delete the group")
	}

	if _, err := th.App.DeleteGroup(group.Id); err == nil {
		t.Fatal("Should not delete the group again - group already deleted")
	}
}

func TestCreateGroupMember(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()

	if _, err := th.App.CreateGroupMember(group.Id, th.BasicUser.Id); err != nil {
		t.Log(err)
		t.Fatal("Should create a group member")
	}

	if _, err := th.App.CreateGroupMember(group.Id, th.BasicUser.Id); err == nil {
		t.Fatal("Should not create a new group member - group member already exist")
	}
}

func TestDeleteGroupMember(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupMember, err := th.App.CreateGroupMember(group.Id, th.BasicUser.Id)
	if err != nil {
		t.Log(err)
		t.Fatal("Should create a group member")
	}

	if _, err := th.App.DeleteGroupMember(groupMember.GroupId, groupMember.UserId); err != nil {
		t.Log(err)
		t.Fatal("Should delete group member")
	}

	if _, err := th.App.DeleteGroupMember(groupMember.GroupId, groupMember.UserId); err == nil {
		t.Fatal("Should not re-delete group member - group member already deleted")
	}
}

func TestCreateGroupTeam(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupTeam := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: th.BasicTeam.Id,
	}

	if _, err := th.App.CreateGroupTeam(groupTeam); err != nil {
		t.Log(err)
		t.Fatal("Should create group team")
	}

	if _, err := th.App.CreateGroupTeam(groupTeam); err == nil {
		t.Fatal("Should not create group team - group team already exists")
	}
}

func TestGetGroupTeam(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupTeam := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: th.BasicTeam.Id,
	}

	// Create GroupTeam
	if _, err := th.App.CreateGroupTeam(groupTeam); err != nil {
		t.Log(err)
		t.Fatal("Should create group team")
	}

	if _, err := th.App.GetGroupTeam(group.Id, th.BasicTeam.Id); err != nil {
		t.Log(err)
		t.Fatal("Should delete group team")
	}
}

func TestGetGroupTeamsPage(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()

	// Create a group team
	groupTeam := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: th.BasicTeam.Id,
	}

	if _, err := th.App.CreateGroupTeam(groupTeam); err != nil {
		t.Log(err)
		t.Fatal("Should create group team")
	}

	groupTeams, err := th.App.GetGroupTeamsPage(group.Id, 0, 99)
	if err != nil {
		t.Log(err)
		t.Fatal("Should have group teams")
	}

	if len(groupTeams) < 1 {
		t.Fatal("Should have retrieved at least one group team")
	}

	if groupTeams, _ = th.App.GetGroupTeamsPage(group.Id, 999, 1); len(groupTeams) > 0 {
		t.Fatal("Should not have group teams")
	}
}

func TestDeleteGroupTeam(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupTeam := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: th.BasicTeam.Id,
	}

	// Create GroupTeam
	if _, err := th.App.CreateGroupTeam(groupTeam); err != nil {
		t.Log(err)
		t.Fatal("Should create group team")
	}

	if _, err := th.App.DeleteGroupTeam(group.Id, th.BasicTeam.Id); err != nil {
		t.Log(err)
		t.Fatal("Should delete group team")
	}

	if _, err := th.App.DeleteGroupTeam(group.Id, th.BasicTeam.Id); err == nil {
		t.Fatal("Should not re-delete group team - group team already deleted")
	}
}

func TestCreateGroupChannel(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupChannel := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: th.BasicChannel.Id,
	}

	if _, err := th.App.CreateGroupChannel(groupChannel); err != nil {
		t.Log(err)
		t.Fatal("Should create group channel")
	}

	if _, err := th.App.CreateGroupChannel(groupChannel); err == nil {
		t.Fatal("Should not create group channel - group channel already exists")
	}
}

func TestGetGroupChannel(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupChannel := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: th.BasicChannel.Id,
	}

	// Create GroupChannel
	if _, err := th.App.CreateGroupChannel(groupChannel); err != nil {
		t.Log(err)
		t.Fatal("Should create group channel")
	}

	if _, err := th.App.GetGroupChannel(group.Id, th.BasicChannel.Id); err != nil {
		t.Log(err)
		t.Fatal("Should delete group channel")
	}
}

func TestGetGroupChannelsPage(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()

	// Create a group channel
	groupChannel := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: th.BasicChannel.Id,
	}

	if _, err := th.App.CreateGroupChannel(groupChannel); err != nil {
		t.Log(err)
		t.Fatal("Should create group channel")
	}

	groupTeams, err := th.App.GetGroupChannelsPage(group.Id, 0, 99)
	if err != nil {
		t.Log(err)
		t.Fatal("Should have group channels")
	}

	if len(groupTeams) < 1 {
		t.Fatal("Should have retrieved at least one group team")
	}

	if groupTeams, _ = th.App.GetGroupChannelsPage(group.Id, 999, 1); len(groupTeams) > 0 {
		t.Fatal("Should not have group channels.")
	}
}

func TestDeleteGroupChannel(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	group := th.CreateGroup()
	groupChannel := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: th.BasicChannel.Id,
	}

	// Create GroupChannel
	if _, err := th.App.CreateGroupChannel(groupChannel); err != nil {
		t.Log(err)
		t.Fatal("Should create group channel")
	}

	if _, err := th.App.DeleteGroupChannel(group.Id, th.BasicChannel.Id); err != nil {
		t.Log(err)
		t.Fatal("Should delete group channel")
	}

	if _, err := th.App.DeleteGroupChannel(group.Id, th.BasicChannel.Id); err == nil {
		t.Fatal("Should not re-delete group channel - group channel already deleted")
	}
}
