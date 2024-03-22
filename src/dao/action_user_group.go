package dao

import (
	"context"
	"owl/model"
)

// ListActionUserGroupsUsers 根据ActionId列出被ActionUserGroup对应的所有用户
func (d *Dao) ListActionUserGroupsUsers(ctx context.Context, actionId uint64) (users []*model.User, err error) {
	groupIdSubQuery := d.getDbWithCtx(ctx).Table("action_user_group").
		Select("user_group_id").
		Where("action_id=?", actionId)

	usersIdSubQuery := d.getDbWithCtx(ctx).Table("user_group_user").
		Select("user_id").
		Where("user_group_id IN (?)", groupIdSubQuery)

	res := d.db.Where("id IN (?)", usersIdSubQuery).
		Find(&users)

	return users, res.Error
}
