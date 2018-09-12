// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"context"

	"github.com/mattermost/mattermost-server/model"
)

func (s *LocalCacheSupplier) handleClusterInvalidateGroup(msg *model.ClusterMessage) {
	if msg.Data == CLEAR_CACHE_MESSAGE_DATA {
		s.groupCache.Purge()
	} else {
		s.groupCache.Remove(msg.Data)
	}
}

func (s *LocalCacheSupplier) GroupSave(ctx context.Context, group *model.Group, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	if len(group.Id) != 0 {
		defer s.doInvalidateCacheCluster(s.groupCache, group.Id)
	}
	return s.Next().GroupSave(ctx, group, hints...)
}

func (s *LocalCacheSupplier) GroupGet(ctx context.Context, groupId string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	if result := s.doStandardReadCache(ctx, s.groupCache, groupId, hints...); result != nil {
		return result
	}

	result := s.Next().GroupGet(ctx, groupId, hints...)

	s.doStandardAddToCache(ctx, s.groupCache, groupId, result, hints...)

	return result
}

func (s *LocalCacheSupplier) GroupGetAllPage(ctx context.Context, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupGetAllPage(ctx, offset, limit, hints...)
}

func (s *LocalCacheSupplier) GroupDelete(ctx context.Context, groupId string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	defer s.doInvalidateCacheCluster(s.groupCache, groupId)
	defer s.doClearCacheCluster(s.roleCache)

	return s.Next().GroupDelete(ctx, groupId, hints...)
}

func (s *LocalCacheSupplier) GroupCreateMember(ctx context.Context, groupID string, userID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupCreateMember(ctx, groupID, userID, hints...)
}

func (s *LocalCacheSupplier) GroupDeleteMember(ctx context.Context, groupID string, userID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupDeleteMember(ctx, groupID, userID, hints...)
}

func (s *LocalCacheSupplier) GroupGetGroupTeam(ctx context.Context, groupID string, teamID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupGetGroupTeam(ctx, groupID, teamID, hints...)
}

func (s *LocalCacheSupplier) GroupGetAllGroupTeamsByGroupPage(ctx context.Context, groupID string, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupGetAllGroupTeamsByGroupPage(ctx, groupID, offset, limit, hints...)
}

func (s *LocalCacheSupplier) GroupSaveGroupTeam(ctx context.Context, groupTeam *model.GroupTeam, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupSaveGroupTeam(ctx, groupTeam, hints...)
}

func (s *LocalCacheSupplier) GroupDeleteGroupTeam(ctx context.Context, groupID string, teamID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupDeleteGroupTeam(ctx, groupID, teamID, hints...)
}

func (s *LocalCacheSupplier) GroupGetGroupChannel(ctx context.Context, groupID string, channelID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupGetGroupChannel(ctx, groupID, channelID, hints...)
}

func (s *LocalCacheSupplier) GroupGetAllGroupChannelsByGroupPage(ctx context.Context, groupID string, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupGetAllGroupChannelsByGroupPage(ctx, groupID, offset, limit, hints...)
}

func (s *LocalCacheSupplier) GroupSaveGroupChannel(ctx context.Context, groupChannel *model.GroupChannel, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupSaveGroupChannel(ctx, groupChannel, hints...)
}

func (s *LocalCacheSupplier) GroupDeleteGroupChannel(ctx context.Context, groupID string, channelID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	return s.Next().GroupDeleteGroupChannel(ctx, groupID, channelID, hints...)
}
