// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"context"

	"github.com/mattermost/mattermost-server/model"
)

func (s *RedisSupplier) GroupSave(ctx context.Context, group *model.Group, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupSave(ctx, group, hints...)
}

func (s *RedisSupplier) GroupGet(ctx context.Context, groupId string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGet(ctx, groupId, hints...)
}

func (s *RedisSupplier) GroupGetAllPage(ctx context.Context, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGetAllPage(ctx, offset, limit, hints...)
}

func (s *RedisSupplier) GroupDelete(ctx context.Context, groupId string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupDelete(ctx, groupId, hints...)
}

func (s *RedisSupplier) GroupCreateMember(ctx context.Context, groupID string, userID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupCreateMember(ctx, groupID, userID, hints...)
}

func (s *RedisSupplier) GroupDeleteMember(ctx context.Context, groupID string, userID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupDeleteMember(ctx, groupID, userID, hints...)
}

func (s *RedisSupplier) GroupGetGroupTeam(ctx context.Context, groupID string, teamID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGetGroupTeam(ctx, groupID, teamID, hints...)
}

func (s *RedisSupplier) GroupGetAllGroupTeamsByGroupPage(ctx context.Context, groupID string, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGetAllGroupTeamsByGroupPage(ctx, groupID, offset, limit, hints...)
}

func (s *RedisSupplier) GroupSaveGroupTeam(ctx context.Context, groupTeam *model.GroupTeam, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupSaveGroupTeam(ctx, groupTeam, hints...)
}

func (s *RedisSupplier) GroupDeleteGroupTeam(ctx context.Context, groupID string, teamID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupDeleteGroupTeam(ctx, groupID, teamID, hints...)
}

func (s *RedisSupplier) GroupGetGroupChannel(ctx context.Context, groupID string, channelID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGetGroupChannel(ctx, groupID, channelID, hints...)
}

func (s *RedisSupplier) GroupGetAllGroupChannelsByGroupPage(ctx context.Context, groupID string, offset int, limit int, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupGetAllGroupChannelsByGroupPage(ctx, groupID, offset, limit, hints...)
}

func (s *RedisSupplier) GroupSaveGroupChannel(ctx context.Context, groupChannel *model.GroupChannel, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupSaveGroupChannel(ctx, groupChannel, hints...)
}

func (s *RedisSupplier) GroupDeleteGroupChannel(ctx context.Context, groupID string, channelID string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	// TODO: Redis caching.
	return s.Next().GroupDeleteGroupChannel(ctx, groupID, channelID, hints...)
}
