// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import "net/http"

type GroupSyncable struct {
	GroupId  string `json:"group_id"`
	CanLeave bool   `json:"can_leave"`
	AutoAdd  bool   `json:"auto_add"`
	CreateAt int64  `json:"create_at"`
	DeleteAt int64  `json:"delete_at"`
	UpdateAt int64  `json:"update_at"`
}

func (syncable *GroupSyncable) SyncableIsValid() *AppError {
	if len(syncable.GroupId) != 26 {
		return NewAppError("GroupSyncable.SyncableIsValid", "model.group_syncable.group_id.app_error", nil, "", http.StatusBadRequest)
	}
	if syncable.AutoAdd == false && syncable.CanLeave == false {
		return NewAppError("GroupSyncable.SyncableIsValid", "model.group_syncable.invalid_state", nil, "", http.StatusBadRequest)
	}
	return nil
}
