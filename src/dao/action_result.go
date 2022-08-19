package dao

import "owl/model"

func (d *Dao) NewActionResult(
	strategyEventID, actionId uint64,
	actionType, actionKind, success bool,
	scriptId, userId uint,
	username, phone, email, wechat string,
	subject, content, response string,
	count int,
) (*model.ActionResult, error) {
	ar := model.ActionResult{
		StrategyEventId: strategyEventID,
		ActionId:        actionId,

		ActionType: actionType,
		ActionKind: actionKind,
		Success:    success,

		ScriptId: scriptId,
		UserId:   userId,

		Username: username,
		Phone:    phone,
		Email:    email,
		Wechat:   wechat,

		Subject:  subject,
		Content:  content,
		Response: response,

		Count: count,
	}

	res := d.db.Create(&ar)
	return &ar, res.Error
}
