package dao

import "owl/model"

// ListActionUserGroupsUsers 根据ActionId列出被ActionUserGroup对应的所有用户
func (d *Dao) ListActionUserGroupsUsers(actionId uint64) (users []*model.User, err error) {
	groupIdSubQuery := d.db.Table("action_user_group").
		Select("user_group_id").
		Where("action_id=?", actionId)

	usersIdSubQuery := d.db.Table("user_group_user").
		Select("user_id").
		Where("user_group_id IN (?)", groupIdSubQuery)

	res := d.db.Where("id IN (?)", usersIdSubQuery).
		Find(&users)

	return users, res.Error
}
