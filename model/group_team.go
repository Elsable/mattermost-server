// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import "net/http"

type GroupTeam struct {
	TeamId string `json:"team_id"`
	GroupSyncable
}

func (groupTeam *GroupTeam) IsValid() *AppError {
	if err := groupTeam.SyncableIsValid(); err != nil {
		return err
	}
	if len(groupTeam.TeamId) != 26 {
		return NewAppError("GroupTeam.IsValid", "model.group_syncable.team_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
