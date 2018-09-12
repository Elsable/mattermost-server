// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import "net/http"

type GroupChannel struct {
	ChannelId string `json:"channel_id"`
	GroupSyncable
}

func (groupChannel *GroupChannel) IsValid() *AppError {
	if err := groupChannel.SyncableIsValid(); err != nil {
		return err
	}
	if len(groupChannel.ChannelId) != 26 {
		return NewAppError("GroupChannel.IsValid", "model.group_channel.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
