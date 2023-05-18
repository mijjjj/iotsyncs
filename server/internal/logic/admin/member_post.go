// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

type sAdminMemberPost struct{}

func NewAdminMemberPost() *sAdminMemberPost {
	return &sAdminMemberPost{}
}

func init() {
	service.RegisterAdminMemberPost(NewAdminMemberPost())
}

func (s *sAdminMemberPost) UpdatePostIds(ctx context.Context, memberId int64, postIds []int64) (err error) {
	_, err = dao.AdminMemberPost.Ctx(ctx).Where("member_id", memberId).Delete()
	if err != nil {
		err = gerror.Wrap(err, "删除失败")
		return
	}

	for i := 0; i < len(postIds); i++ {
		_, err = dao.AdminMemberPost.Ctx(ctx).
			Insert(entity.AdminMemberPost{
				MemberId: memberId,
				PostId:   postIds[i],
			})
		if err != nil {
			err = gerror.Wrap(err, "插入用户岗位失败")
			return err
		}
	}

	return
}

// GetMemberByIds 获取指定用户的岗位ids
func (s *sAdminMemberPost) GetMemberByIds(ctx context.Context, memberId int64) (postIds []int64, err error) {
	var list []*entity.AdminMemberPost
	err = dao.AdminMemberPost.Ctx(ctx).Fields("post_id").Where("member_id", memberId).Scan(&list)
	if err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	for i := 0; i < len(list); i++ {
		postIds = append(postIds, list[i].PostId)
	}

	return
}
