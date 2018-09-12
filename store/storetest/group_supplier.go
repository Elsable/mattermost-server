// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

func TestGroupStore(t *testing.T, ss store.Store) {
	t.Run("Save", func(t *testing.T) { testGroupStoreSave(t, ss) })
	t.Run("Get", func(t *testing.T) { testGroupStoreGet(t, ss) })
	t.Run("GetAllPage", func(t *testing.T) { testGroupStoreGetAllPage(t, ss) })
	t.Run("Delete", func(t *testing.T) { testGroupStoreDelete(t, ss) })

	t.Run("CreateMember", func(t *testing.T) { testGroupCreateMember(t, ss) })
	t.Run("DeleteMember", func(t *testing.T) { testGroupDeleteMember(t, ss) })

	t.Run("GetGroupTeam", func(t *testing.T) { testGetGroupTeam(t, ss) })
	t.Run("GetAllGroupTeamsByGroupPage", func(t *testing.T) { testGetAllGroupTeamsByGroupPage(t, ss) })
	t.Run("SaveGroupTeam", func(t *testing.T) { testSaveGroupTeam(t, ss) })
	t.Run("DeleteGroupTeam", func(t *testing.T) { testDeleteGroupTeam(t, ss) })

	t.Run("GetGroupChannel", func(t *testing.T) { testGetGroupChannel(t, ss) })
	t.Run("GetAllGroupChannelsByGroupPage", func(t *testing.T) { testGetAllGroupChannelsByGroupPage(t, ss) })
	t.Run("SaveGroupChannel", func(t *testing.T) { testSaveGroupChannel(t, ss) })
	t.Run("DeleteGroupChannel", func(t *testing.T) { testDeleteGroupChannel(t, ss) })
}

func testGroupStoreSave(t *testing.T, ss store.Store) {
	// Save a new group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}

	// Happy path
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)
	assert.Equal(t, g1.Name, d1.Name)
	assert.Equal(t, g1.DisplayName, d1.DisplayName)
	assert.Equal(t, g1.Description, d1.Description)
	assert.Equal(t, g1.RemoteId, d1.RemoteId)
	assert.NotZero(t, d1.CreateAt)
	assert.NotZero(t, d1.UpdateAt)
	assert.Zero(t, d1.DeleteAt)

	// Requires name and display name
	g2 := &model.Group{
		Name:        "",
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res2 := <-ss.Group().Save(g2)
	assert.Nil(t, res2.Data)
	assert.NotNil(t, res2.Err)
	assert.Equal(t, res2.Err.Id, "model.group.name.app_error")

	g2.Name = model.NewId()
	g2.DisplayName = ""
	res3 := <-ss.Group().Save(g2)
	assert.Nil(t, res3.Data)
	assert.NotNil(t, res3.Err)
	assert.Equal(t, res3.Err.Id, "model.group.display_name.app_error")

	// Can't invent an ID and save it
	g3 := &model.Group{
		Id:          model.NewId(),
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
		CreateAt:    1,
		UpdateAt:    1,
	}
	res4 := <-ss.Group().Save(g3)
	assert.Nil(t, res4.Data)
	assert.Equal(t, res4.Err.Id, "store.sql_group.save.missing.app_error")

	// Won't accept a duplicate name
	g4 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res5 := <-ss.Group().Save(g4)
	assert.Nil(t, res5.Err)
	g4b := &model.Group{
		Name:        g4.Name,
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res5b := <-ss.Group().Save(g4b)
	assert.Nil(t, res5b.Data)
	assert.Equal(t, res5b.Err.Id, "store.sql_group.save.insert.app_error")

	// Fields cannot be greater than max values
	g5 := &model.Group{
		Name:        strings.Repeat("x", model.GroupNameMaxLength),
		DisplayName: strings.Repeat("x", model.GroupDisplayNameMaxLength),
		Description: strings.Repeat("x", model.GroupDescriptionMaxLength),
		Type:        model.GroupTypeLdap,
	}
	assert.Nil(t, g5.IsValidForCreate())

	g5.Name = g5.Name + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.name.app_error")
	g5.Name = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	g5.DisplayName = g5.DisplayName + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.display_name.app_error")
	g5.DisplayName = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	g5.Description = g5.Description + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.description.app_error")
	g5.Description = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	// Must use a valid type
	g6 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        "fake",
	}
	assert.Equal(t, g6.IsValidForCreate().Id, "model.group.type.app_error")

	// Cannot update CreateAt or DeleteAt
	origCreateAt := d1.CreateAt
	origDeleteAt := d1.DeleteAt
	d1.CreateAt = model.GetMillis()
	d1.DeleteAt = model.GetMillis()
	res6 := <-ss.Group().Save(d1)
	d2 := res6.Data.(*model.Group)
	assert.Equal(t, origCreateAt, d2.CreateAt)
	assert.Equal(t, origDeleteAt, d2.DeleteAt)
}

func testGroupStoreGet(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
		RemoteId:    model.NewId(),
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)

	// Get the group
	res2 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res2.Err)
	d2 := res2.Data.(*model.Group)
	assert.Equal(t, d1.Id, d2.Id)
	assert.Equal(t, d1.Name, d2.Name)
	assert.Equal(t, d1.DisplayName, d2.DisplayName)
	assert.Equal(t, d1.Description, d2.Description)
	assert.Equal(t, d1.RemoteId, d2.RemoteId)
	assert.Equal(t, d1.CreateAt, d2.CreateAt)
	assert.Equal(t, d1.UpdateAt, d2.UpdateAt)
	assert.Equal(t, d1.DeleteAt, d2.DeleteAt)

	// Get an invalid group
	res3 := <-ss.Group().Get(model.NewId())
	assert.NotNil(t, res3.Err)
	assert.Equal(t, res3.Err.Id, "store.sql_group.get.app_error")
}

func testGroupStoreGetAllPage(t *testing.T, ss store.Store) {
	numGroups := 10

	groups := []*model.Group{}

	// Create groups
	for i := 0; i < numGroups; i++ {
		g := &model.Group{
			Name:        model.NewId(),
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Type:        model.GroupTypeLdap,
		}
		groups = append(groups, g)
		res := <-ss.Group().Save(g)
		assert.Nil(t, res.Err)
	}

	// Returns all the groups
	res1 := <-ss.Group().GetAllPage(0, 999)
	d1 := res1.Data.([]*model.Group)
	assert.Condition(t, func() bool { return len(d1) >= numGroups })
	for _, expectedGroup := range groups {
		present := false
		for _, dbGroup := range d1 {
			if dbGroup.Id == expectedGroup.Id {
				present = true
				break
			}
		}
		assert.True(t, present)
	}

	// Returns the correct number based on limit
	res2 := <-ss.Group().GetAllPage(0, 2)
	d2 := res2.Data.([]*model.Group)
	assert.Len(t, d2, 2)

	// Check that result sets are different using an offset
	res3 := <-ss.Group().GetAllPage(0, 5)
	d3 := res3.Data.([]*model.Group)
	res4 := <-ss.Group().GetAllPage(5, 5)
	d4 := res4.Data.([]*model.Group)
	for _, d3i := range d3 {
		for _, d4i := range d4 {
			if d4i.Id == d3i.Id {
				t.Error("Expected results to be unique.")
			}
		}
	}
}

func testGroupStoreDelete(t *testing.T, ss store.Store) {
	// Save a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
	}

	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)

	// Check the group is retrievable
	res2 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res2.Err)

	// Get the before count
	res7 := <-ss.Group().GetAllPage(0, 999)
	d7 := res7.Data.([]*model.Group)
	beforeCount := len(d7)

	// Delete the group
	res3 := <-ss.Group().Delete(d1.Id)
	assert.Nil(t, res3.Err)

	// Check the group is deleted
	res4 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res4.Err)
	d2 := res4.Data.(*model.Group)
	assert.NotZero(t, d2.DeleteAt)

	// Check the after count
	res5 := <-ss.Group().GetAllPage(0, 999)
	d5 := res5.Data.([]*model.Group)
	afterCount := len(d5)
	assert.Condition(t, func() bool { return beforeCount == afterCount+1 })

	// Try and delete a nonexistent group
	res6 := <-ss.Group().Delete(model.NewId())
	assert.NotNil(t, res6.Err)
	assert.Equal(t, res6.Err.Id, "store.sql_group.get.app_error")

	// Cannot delete again
	res8 := <-ss.Group().Delete(d1.Id)
	assert.Equal(t, res8.Err.Id, "store.sql_group.delete.already_deleted")
}

func testGroupCreateMember(t *testing.T, ss store.Store) {
	// Create group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create user
	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	res2 := <-ss.User().Save(u1)
	assert.Nil(t, res2.Err)
	user := res2.Data.(*model.User)

	// Happy path
	res3 := <-ss.Group().CreateMember(group.Id, user.Id)
	assert.Nil(t, res3.Err)
	d2 := res3.Data.(*model.GroupMember)
	assert.Equal(t, d2.GroupId, group.Id)
	assert.Equal(t, d2.UserId, user.Id)
	assert.NotZero(t, d2.CreateAt)
	assert.Zero(t, d2.DeleteAt)

	// Duplicate composite key (GroupId, UserId)
	res4 := <-ss.Group().CreateMember(group.Id, user.Id)
	assert.Equal(t, res4.Err.Id, "store.sql_group.save_member.exists.app_error")

	// Invalid UserId
	res5 := <-ss.Group().CreateMember(group.Id, model.NewId())
	assert.Equal(t, res5.Err.Id, "store.sql_group.save_member.save.app_error")

	// Invalid GroupId
	res6 := <-ss.Group().CreateMember(model.NewId(), user.Id)
	assert.Equal(t, res6.Err.Id, "store.sql_group.save_member.save.app_error")
}

func testGroupDeleteMember(t *testing.T, ss store.Store) {
	// Create group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create user
	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	res2 := <-ss.User().Save(u1)
	assert.Nil(t, res2.Err)
	user := res2.Data.(*model.User)

	// Create member
	res3 := <-ss.Group().CreateMember(group.Id, user.Id)
	assert.Nil(t, res3.Err)
	d1 := res3.Data.(*model.GroupMember)

	// Happy path
	res4 := <-ss.Group().DeleteMember(group.Id, user.Id)
	assert.Nil(t, res4.Err)
	d2 := res4.Data.(*model.GroupMember)
	assert.Equal(t, d2.GroupId, group.Id)
	assert.Equal(t, d2.UserId, user.Id)
	assert.Equal(t, d2.CreateAt, d1.CreateAt)
	assert.NotZero(t, d2.DeleteAt)

	// Delete an already deleted member
	res5 := <-ss.Group().DeleteMember(group.Id, user.Id)
	assert.Equal(t, res5.Err.Id, "store.sql_group.delete_member.already_deleted")

	// Delete with invalid UserId
	res6 := <-ss.Group().DeleteMember(group.Id, strings.Repeat("x", 27))
	assert.Equal(t, res6.Err.Id, "model.group_member.user_id.app_error")

	// Delete with invalid GroupId
	res7 := <-ss.Group().DeleteMember(strings.Repeat("x", 27), user.Id)
	assert.Equal(t, res7.Err.Id, "model.group_member.group_id.app_error")

	// Delete with non-existent User
	res8 := <-ss.Group().DeleteMember(group.Id, model.NewId())
	assert.Equal(t, res8.Err.Id, "store.sql_group.get_member.missing.app_error")

	// Delete non-existent Group
	res9 := <-ss.Group().DeleteMember(model.NewId(), group.Id)
	assert.Equal(t, res9.Err.Id, "store.sql_group.get_member.missing.app_error")
}

func testGetGroupTeam(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res2 := <-ss.Team().Save(t1)
	assert.Nil(t, res2.Err)
	team := res2.Data.(*model.Team)

	// Create GroupTeam
	gt1 := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: team.Id,
	}
	res3 := <-ss.Group().SaveGroupTeam(gt1)
	assert.Nil(t, res3.Err)
	groupTeam := res3.Data.(*model.GroupTeam)

	// Get GroupTeam
	res4 := <-ss.Group().GetGroupTeam(groupTeam.GroupId, groupTeam.TeamId)
	assert.Nil(t, res4.Err)
	dgt := res4.Data.(*model.GroupTeam)
	assert.Equal(t, gt1.GroupId, dgt.GroupId)
	assert.Equal(t, gt1.TeamId, dgt.TeamId)
	assert.Equal(t, gt1.CanLeave, dgt.CanLeave)
	assert.Equal(t, gt1.AutoAdd, dgt.AutoAdd)
	assert.NotZero(t, gt1.CreateAt)
	assert.NotZero(t, gt1.UpdateAt)
	assert.Zero(t, gt1.DeleteAt)
}

func testGetAllGroupTeamsByGroupPage(t *testing.T, ss store.Store) {
	numGroupTeams := 10

	// Create group
	g := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	groupTeams := []*model.GroupTeam{}

	// Create groupTeams
	for i := 0; i < numGroupTeams; i++ {
		// Create Team
		t1 := &model.Team{
			DisplayName:     "Name",
			Description:     "Some description",
			CompanyName:     "Some company name",
			AllowOpenInvite: false,
			InviteId:        "inviteid0",
			Name:            "z-z-" + model.NewId() + "a",
			Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
			Type:            model.TEAM_OPEN,
		}
		res2 := <-ss.Team().Save(t1)
		assert.Nil(t, res2.Err)
		team := res2.Data.(*model.Team)

		// create groupteam
		res3 := <-ss.Group().SaveGroupTeam(&model.GroupTeam{
			GroupSyncable: model.GroupSyncable{
				GroupId:  group.Id,
				CanLeave: true,
			},
			TeamId: team.Id,
		})
		assert.Nil(t, res3.Err)
		groupTeam := res3.Data.(*model.GroupTeam)
		groupTeams = append(groupTeams, groupTeam)
	}

	// Returns all the group teams
	res4 := <-ss.Group().GetAllGroupTeamsByGroupPage(group.Id, 0, 999)
	d1 := res4.Data.([]*model.GroupTeam)
	assert.Condition(t, func() bool { return len(d1) >= numGroupTeams })
	for _, expectedGroupTeam := range groupTeams {
		present := false
		for _, dbGroupTeam := range d1 {
			if dbGroupTeam.GroupId == expectedGroupTeam.GroupId && dbGroupTeam.TeamId == expectedGroupTeam.TeamId {
				present = true
				break
			}
		}
		assert.True(t, present)
	}

	// Returns the correct number based on limit
	res5 := <-ss.Group().GetAllGroupTeamsByGroupPage(group.Id, 0, 2)
	d2 := res5.Data.([]*model.GroupTeam)
	assert.Len(t, d2, 2)

	// Check that result sets are different using an offset
	res6 := <-ss.Group().GetAllGroupTeamsByGroupPage(group.Id, 0, 5)
	d3 := res6.Data.([]*model.GroupTeam)
	res7 := <-ss.Group().GetAllGroupTeamsByGroupPage(group.Id, 5, 5)
	d4 := res7.Data.([]*model.GroupTeam)
	for _, d3i := range d3 {
		for _, d4i := range d4 {
			if d4i.GroupId == d3i.GroupId && d4i.TeamId == d3i.TeamId {
				t.Error("Expected results to be unique.")
			}
		}
	}
}

func testSaveGroupTeam(t *testing.T, ss store.Store) {
	// Invalid TeamID
	res1 := <-ss.Group().SaveGroupTeam(&model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: true,
		},
		TeamId: "x",
	})
	assert.Equal(t, res1.Err.Id, "model.group_syncable.team_id.app_error")

	// Invalid GroupID
	res2 := <-ss.Group().SaveGroupTeam(&model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  "x",
			CanLeave: true,
		},
		TeamId: model.NewId(),
	})
	assert.Equal(t, res2.Err.Id, "model.group_syncable.group_id.app_error")

	// Invalid CanLeave/AutoAdd combo (both false)
	res3 := <-ss.Group().SaveGroupTeam(&model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: false,
			AutoAdd:  false,
		},
		TeamId: model.NewId(),
	})
	assert.Equal(t, res3.Err.Id, "model.group_syncable.invalid_state")

	// Create Group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res4 := <-ss.Group().Save(g1)
	assert.Nil(t, res4.Err)
	group := res4.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res5 := <-ss.Team().Save(t1)
	assert.Nil(t, res5.Err)
	team := res5.Data.(*model.Team)

	// New GroupTeam, happy path
	gt1 := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: team.Id,
	}
	res6 := <-ss.Group().SaveGroupTeam(gt1)
	assert.Nil(t, res6.Err)
	d1 := res6.Data.(*model.GroupTeam)
	assert.Equal(t, gt1.TeamId, d1.TeamId)
	assert.Equal(t, gt1.GroupId, d1.GroupId)
	assert.Equal(t, gt1.CanLeave, d1.CanLeave)
	assert.Equal(t, gt1.AutoAdd, d1.AutoAdd)
	assert.NotZero(t, d1.CreateAt)
	assert.Zero(t, d1.DeleteAt)

	// Update existing group team
	gt1.CanLeave = false
	gt1.AutoAdd = true
	res7 := <-ss.Group().SaveGroupTeam(gt1)
	assert.Nil(t, res7.Err)
	d2 := res7.Data.(*model.GroupTeam)
	assert.False(t, d2.CanLeave)
	assert.True(t, d2.AutoAdd)

	// Update to invalid state
	gt1.AutoAdd = false
	gt1.CanLeave = false
	res8 := <-ss.Group().SaveGroupTeam(gt1)
	assert.Equal(t, res8.Err.Id, "model.group_syncable.invalid_state")

	// Non-existent Group
	gt2 := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: team.Id,
	}
	res9 := <-ss.Group().SaveGroupTeam(gt2)
	assert.Equal(t, res9.Err.Id, "store.sql_group.save_group_team.save.app_error")

	// Non-existent Team
	gt3 := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: model.NewId(),
	}
	res10 := <-ss.Group().SaveGroupTeam(gt3)
	assert.Equal(t, res10.Err.Id, "store.sql_group.save_group_team.save.app_error")

	// Cannot update CreateAt or DeleteAt
	origCreateAt := d1.CreateAt
	origDeleteAt := d1.DeleteAt
	d1.CreateAt = model.GetMillis()
	d1.DeleteAt = model.GetMillis()
	res11 := <-ss.Group().SaveGroupTeam(d1)
	d3 := res11.Data.(*model.GroupTeam)
	assert.Equal(t, origCreateAt, d3.CreateAt)
	assert.Equal(t, origDeleteAt, d3.DeleteAt)
}

func testDeleteGroupTeam(t *testing.T, ss store.Store) {
	// Create Group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res2 := <-ss.Team().Save(t1)
	assert.Nil(t, res2.Err)
	team := res2.Data.(*model.Team)

	// Create GroupTeam
	gt1 := &model.GroupTeam{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		TeamId: team.Id,
	}
	res7 := <-ss.Group().SaveGroupTeam(gt1)
	assert.Nil(t, res7.Err)
	groupTeam := res7.Data.(*model.GroupTeam)

	// Invalid GroupId
	res3 := <-ss.Group().DeleteGroupTeam("x", groupTeam.TeamId)
	assert.Equal(t, res3.Err.Id, "store.sql_group.delete_group_team.group_id.invalid")

	// Invalid TeamId
	res4 := <-ss.Group().DeleteGroupTeam(groupTeam.GroupId, "x")
	assert.Equal(t, res4.Err.Id, "store.sql_group.delete_group_team.team_id.invalid")

	// Non-existent Group
	res5 := <-ss.Group().DeleteGroupTeam(model.NewId(), groupTeam.TeamId)
	assert.Equal(t, res5.Err.Id, "store.sql_group.delete_group_team.app_error")

	// Non-existent Team
	res6 := <-ss.Group().DeleteGroupTeam(groupTeam.GroupId, model.NewId())
	assert.Equal(t, res6.Err.Id, "store.sql_group.delete_group_team.app_error")

	// Happy path...
	res8 := <-ss.Group().DeleteGroupTeam(groupTeam.GroupId, groupTeam.TeamId)
	assert.Nil(t, res8.Err)
	d1 := res8.Data.(*model.GroupTeam)
	assert.NotZero(t, d1.DeleteAt)
	assert.Equal(t, d1.GroupId, groupTeam.GroupId)
	assert.Equal(t, d1.TeamId, groupTeam.TeamId)
	assert.Equal(t, d1.CanLeave, groupTeam.CanLeave)
	assert.Equal(t, d1.AutoAdd, groupTeam.AutoAdd)
	assert.Equal(t, d1.CreateAt, groupTeam.CreateAt)
	assert.Condition(t, func() bool { return d1.UpdateAt > groupTeam.UpdateAt })

	// Record already deleted
	res9 := <-ss.Group().DeleteGroupTeam(groupTeam.GroupId, groupTeam.TeamId)
	assert.Equal(t, res9.Err.Id, "store.sql_group.delete_group_team.already_deleted")
}

func testGetGroupChannel(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res2 := <-ss.Team().Save(t1)
	assert.Nil(t, res2.Err)
	team := res2.Data.(*model.Team)

	// Create Channel
	id := model.NewId()
	c1 := &model.Channel{
		DisplayName: "Hello World " + id,
		Name:        "hello-world" + id,
		Type:        model.CHANNEL_OPEN,
		TeamId:      team.Id,
	}
	res5 := <-ss.Channel().Save(c1, 9999)
	assert.Nil(t, res5.Err)
	channel := res5.Data.(*model.Channel)

	// Create GroupChannel
	gc1 := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: channel.Id,
	}
	res3 := <-ss.Group().SaveGroupChannel(gc1)
	assert.Nil(t, res3.Err)
	groupChannel := res3.Data.(*model.GroupChannel)

	// Get GroupChannel
	res4 := <-ss.Group().GetGroupChannel(groupChannel.GroupId, groupChannel.ChannelId)
	assert.Nil(t, res4.Err)
	dgt := res4.Data.(*model.GroupChannel)
	assert.Equal(t, gc1.GroupId, dgt.GroupId)
	assert.Equal(t, gc1.ChannelId, dgt.ChannelId)
	assert.Equal(t, gc1.CanLeave, dgt.CanLeave)
	assert.Equal(t, gc1.AutoAdd, dgt.AutoAdd)
	assert.NotZero(t, gc1.CreateAt)
	assert.NotZero(t, gc1.UpdateAt)
	assert.Zero(t, gc1.DeleteAt)
}

func testGetAllGroupChannelsByGroupPage(t *testing.T, ss store.Store) {
	numGroupChannels := 10

	// Create group
	g := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res5 := <-ss.Team().Save(t1)
	assert.Nil(t, res5.Err)
	team := res5.Data.(*model.Team)

	groupChannels := []*model.GroupChannel{}

	// Create groupChannels
	for i := 0; i < numGroupChannels; i++ {
		// Create channel
		id := model.NewId()
		c1 := &model.Channel{
			DisplayName: "Hello World " + id,
			Name:        "hello-world" + id,
			Type:        model.CHANNEL_OPEN,
			TeamId:      team.Id,
		}
		res9 := <-ss.Channel().Save(c1, 9999)
		channel := res9.Data.(*model.Channel)

		// create groupteam
		res3 := <-ss.Group().SaveGroupChannel(&model.GroupChannel{
			GroupSyncable: model.GroupSyncable{
				GroupId:  group.Id,
				CanLeave: true,
			},
			ChannelId: channel.Id,
		})
		assert.Nil(t, res3.Err)
		groupTeam := res3.Data.(*model.GroupChannel)
		groupChannels = append(groupChannels, groupTeam)
	}

	// Returns all the group teams
	res4 := <-ss.Group().GetAllGroupChannelsByGroupPage(group.Id, 0, 999)
	d1 := res4.Data.([]*model.GroupChannel)
	assert.Condition(t, func() bool { return len(d1) >= numGroupChannels })
	for _, expectedGroupTeam := range groupChannels {
		present := false
		for _, dbGroupTeam := range d1 {
			if dbGroupTeam.GroupId == expectedGroupTeam.GroupId && dbGroupTeam.ChannelId == expectedGroupTeam.ChannelId {
				present = true
				break
			}
		}
		assert.True(t, present)
	}

	// Returns the correct number based on limit
	res6 := <-ss.Group().GetAllGroupChannelsByGroupPage(group.Id, 0, 2)
	d2 := res6.Data.([]*model.GroupChannel)
	assert.Len(t, d2, 2)

	// Check that result sets are different using an offset
	res7 := <-ss.Group().GetAllGroupChannelsByGroupPage(group.Id, 0, 5)
	d3 := res7.Data.([]*model.GroupChannel)
	res8 := <-ss.Group().GetAllGroupChannelsByGroupPage(group.Id, 5, 5)
	d4 := res8.Data.([]*model.GroupChannel)
	for _, d3i := range d3 {
		for _, d4i := range d4 {
			if d4i.GroupId == d3i.GroupId && d4i.ChannelId == d3i.ChannelId {
				t.Error("Expected results to be unique.")
			}
		}
	}
}

func testSaveGroupChannel(t *testing.T, ss store.Store) {
	// Invalid ChannelID
	res1 := <-ss.Group().SaveGroupChannel(&model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: true,
		},
		ChannelId: "x",
	})
	assert.Equal(t, res1.Err.Id, "model.group_channel.channel_id.app_error")

	// Invalid GroupID
	res2 := <-ss.Group().SaveGroupChannel(&model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  "x",
			CanLeave: true,
		},
		ChannelId: model.NewId(),
	})
	assert.Equal(t, res2.Err.Id, "model.group_syncable.group_id.app_error")

	// Invalid CanLeave/AutoAdd combo (both false)
	res3 := <-ss.Group().SaveGroupChannel(&model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: false,
			AutoAdd:  false,
		},
		ChannelId: model.NewId(),
	})
	assert.Equal(t, res3.Err.Id, "model.group_syncable.invalid_state")

	// Create Group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res4 := <-ss.Group().Save(g1)
	assert.Nil(t, res4.Err)
	group := res4.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res5 := <-ss.Team().Save(t1)
	assert.Nil(t, res5.Err)
	team := res5.Data.(*model.Team)

	// Create Channel
	id := model.NewId()
	c1 := &model.Channel{
		DisplayName: "Hello World " + id,
		Name:        "hello-world" + id,
		Type:        model.CHANNEL_OPEN,
		TeamId:      team.Id,
	}
	res9 := <-ss.Channel().Save(c1, 9999)
	assert.Nil(t, res9.Err)
	channel := res9.Data.(*model.Channel)

	// New GroupChannel, happy path
	gt1 := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: channel.Id,
	}
	res6 := <-ss.Group().SaveGroupChannel(gt1)
	assert.Nil(t, res6.Err)
	d1 := res6.Data.(*model.GroupChannel)
	assert.Equal(t, gt1.ChannelId, d1.ChannelId)
	assert.Equal(t, gt1.GroupId, d1.GroupId)
	assert.Equal(t, gt1.CanLeave, d1.CanLeave)
	assert.Equal(t, gt1.AutoAdd, d1.AutoAdd)
	assert.NotZero(t, d1.CreateAt)
	assert.Zero(t, d1.DeleteAt)

	// Update existing group team
	res12 := <-ss.Group().GetAllPage(0, 999)
	beforeCount := len(res12.Data.([]*model.Group))

	gt1.CanLeave = false
	gt1.AutoAdd = true
	res7 := <-ss.Group().SaveGroupChannel(gt1)
	assert.Nil(t, res7.Err)
	d2 := res7.Data.(*model.GroupChannel)
	assert.False(t, d2.CanLeave)
	assert.True(t, d2.AutoAdd)

	res13 := <-ss.Group().GetAllPage(0, 999)
	afterCount := len(res13.Data.([]*model.Group))

	assert.Equal(t, afterCount, beforeCount)

	// Update to invalid state
	gt1.AutoAdd = false
	gt1.CanLeave = false
	res8 := <-ss.Group().SaveGroupChannel(gt1)
	assert.Equal(t, res8.Err.Id, "model.group_syncable.invalid_state")

	// Non-existent Group
	gt2 := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  model.NewId(),
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: channel.Id,
	}
	res10 := <-ss.Group().SaveGroupChannel(gt2)
	assert.Equal(t, res10.Err.Id, "store.sql_group.save_group_channel.save.app_error")

	// Non-existent Channel
	gt3 := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: model.NewId(),
	}
	res11 := <-ss.Group().SaveGroupChannel(gt3)
	assert.Equal(t, res11.Err.Id, "store.sql_group.save_group_channel.save.app_error")

	// Cannot update CreateAt or DeleteAt
	origCreateAt := d1.CreateAt
	origDeleteAt := d1.DeleteAt
	d1.CreateAt = model.GetMillis()
	d1.DeleteAt = model.GetMillis()
	res14 := <-ss.Group().SaveGroupChannel(d1)
	d3 := res14.Data.(*model.GroupChannel)
	assert.Equal(t, origCreateAt, d3.CreateAt)
	assert.Equal(t, origDeleteAt, d3.DeleteAt)
}

func testDeleteGroupChannel(t *testing.T, ss store.Store) {
	// Create Group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	group := res1.Data.(*model.Group)

	// Create Team
	t1 := &model.Team{
		DisplayName:     "Name",
		Description:     "Some description",
		CompanyName:     "Some company name",
		AllowOpenInvite: false,
		InviteId:        "inviteid0",
		Name:            "z-z-" + model.NewId() + "a",
		Email:           "success+" + model.NewId() + "@simulator.amazonses.com",
		Type:            model.TEAM_OPEN,
	}
	res2 := <-ss.Team().Save(t1)
	assert.Nil(t, res2.Err)
	team := res2.Data.(*model.Team)

	// Create channel
	id := model.NewId()
	c1 := &model.Channel{
		DisplayName: "Hello World " + id,
		Name:        "hello-world" + id,
		Type:        model.CHANNEL_OPEN,
		TeamId:      team.Id,
	}
	res9 := <-ss.Channel().Save(c1, 9999)
	assert.Nil(t, res9.Err)
	channel := res9.Data.(*model.Channel)

	// Create GroupChannel
	gc1 := &model.GroupChannel{
		GroupSyncable: model.GroupSyncable{
			GroupId:  group.Id,
			CanLeave: true,
			AutoAdd:  false,
		},
		ChannelId: channel.Id,
	}
	res7 := <-ss.Group().SaveGroupChannel(gc1)
	assert.Nil(t, res7.Err)
	groupChannel := res7.Data.(*model.GroupChannel)
	assert.NotNil(t, groupChannel)

	// Invalid GroupId
	res3 := <-ss.Group().DeleteGroupChannel("x", groupChannel.ChannelId)
	assert.Equal(t, res3.Err.Id, "store.sql_group.delete_group_channel.group_id.invalid")

	// Invalid ChannelId
	res4 := <-ss.Group().DeleteGroupChannel(groupChannel.ChannelId, "x")
	assert.Equal(t, res4.Err.Id, "store.sql_group.delete_group_channel.channel_id.invalid")

	// Non-existent Group
	res5 := <-ss.Group().DeleteGroupChannel(model.NewId(), groupChannel.ChannelId)
	assert.Equal(t, res5.Err.Id, "store.sql_group.delete_group_channel.app_error")

	// Non-existent Channel
	res6 := <-ss.Group().DeleteGroupChannel(groupChannel.ChannelId, model.NewId())
	assert.Equal(t, res6.Err.Id, "store.sql_group.delete_group_channel.app_error")

	// Happy path...
	res8 := <-ss.Group().DeleteGroupChannel(groupChannel.GroupId, groupChannel.ChannelId)
	assert.Nil(t, res8.Err)
	d1 := res8.Data.(*model.GroupChannel)
	assert.NotZero(t, d1.DeleteAt)
	assert.Equal(t, d1.GroupId, groupChannel.GroupId)
	assert.Equal(t, d1.ChannelId, groupChannel.ChannelId)
	assert.Equal(t, d1.CanLeave, groupChannel.CanLeave)
	assert.Equal(t, d1.AutoAdd, groupChannel.AutoAdd)
	assert.Equal(t, d1.CreateAt, groupChannel.CreateAt)
	assert.Equal(t, d1.UpdateAt, d1.DeleteAt)
	assert.Condition(t, func() bool { return d1.UpdateAt > groupChannel.UpdateAt })

	// Record already deleted
	res10 := <-ss.Group().DeleteGroupChannel(groupChannel.GroupId, groupChannel.ChannelId)
	assert.Equal(t, res10.Err.Id, "store.sql_group.delete_group_channel.already_deleted")
}
