package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"owl/common/types"
)

type Action struct {
	types.Action
	Users      []types.ActionUser      `form:"users" json:"users"`
	UserGroups []types.ActionUserGroup `form:"user_groups" json:"user_groups" binding:"required"`
}

type Strategy struct {
	types.Strategy
	Hosts    []types.StrategyHost  `form:"hosts" json:"hosts"`
	Groups   []types.StrategyGroup `form:"groups" json:"groups"`
	Triggers []types.Trigger       `form:"triggers" json:"triggers" binding:"required"`
	Actions  []Action              `form:"actions" json:"actions" binding:"required"`
}

type ActionInfo struct {
	types.Action
	Users      []types.User      `json:"users"`
	UserGroups []types.UserGroup `json:"user_groups"`
}

type StrategyInfo struct {
	types.Strategy
	Hosts    []types.Host    `json:"hosts"`
	Groups   []types.Group   `json:"groups"`
	Triggers []types.Trigger `json:"triggers"`
	Actions  []*ActionInfo   `json:"actions"`
}

type StrategyWithName struct {
	types.Strategy
	UserName   string `json:"user_name"`
	ParentName string `json:"parent_name"`
}

func strategiesList(c *gin.Context) {
	var strategies []StrategyWithName
	var total int
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
	}
	key := c.DefaultQuery("key", "")
	my := c.DefaultQuery("my", "false")

	offset := (page - 1) * page_size

	where := ""
	if my == "false" {
		where = fmt.Sprintf("s.name LIKE '%%%s%%'", key)
	} else {
		where = fmt.Sprintf("s.name LIKE '%%%s%%' and s.user_id = %d", key, GetUser(c).ID)
	}

	mydb.Where(where).Table("strategy s").Count(&total)

	mydb.Table("strategy s").Select("s.*, user.username user_name, ps.name parent_name").
		Joins("Left Join user ON s.user_id = user.id").
		Joins("Left Join strategy ps ON s.pid = ps.id").Where(where).Offset(offset).Limit(page_size).Find(&strategies)

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategies": &strategies, "total": total})
}

func strategyCreate(c *gin.Context) {
	var strategy Strategy
	if err := c.BindJSON(&strategy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	strategy.UserID = GetUser(c).ID

	if strategy.Pid != 0 {
		var parent_strategy types.Strategy
		if err := mydb.Where("id = ?", strategy.Pid).First(&parent_strategy).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
			return
		}

		if parent_strategy.Pid != 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{"code": http.StatusForbidden, "message": "not allow"})
			return
		}
	}

	if err := mydb.Create(&strategy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "form": &strategy})
}

func strategyDelete(c *gin.Context) {
	strategy_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
	}

	var strategy types.Strategy
	if err := mydb.Where("id = ?", strategy_id).First(&strategy).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	if !GetUser(c).IsAdmin() && GetUser(c).ID != strategy.UserID {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusForbidden, "message": "not allow"})
		return
	}

	tx := mydb.Begin()

	if err := tx.Where("id = ?", strategy_id).Delete(Strategy{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	if err := tx.Where("pid = ?", strategy_id).Delete(Strategy{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "delete"})
}

func strategyInfo(c *gin.Context) {
	strategy_id := c.Param("id")

	var strategy StrategyInfo
	var hosts []types.Host
	var groups []types.Group
	var triggers []types.Trigger
	var actions []*ActionInfo
	var users []types.User
	var user_groups []types.UserGroup

	if len(strategy_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "strategy_id should be applied"})
		return
	}

	if err := mydb.Where("id = ?", strategy_id).First(&strategy).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	mydb.Select("`host`.*").Joins("Join strategy_host ON strategy_host.host_id = host.id").Where("strategy_id = ?", strategy_id).Find(&hosts)
	mydb.Select("`group`.*").Joins("Join strategy_group ON strategy_group.group_id = group.id").Where("strategy_id = ?", strategy_id).Find(&groups)
	mydb.Where("strategy_id = ?", strategy_id).Order("`index` asc").Find(&triggers)
	mydb.Where("strategy_id = ?", strategy_id).Find(&actions)
	for _, action := range actions {
		mydb.Select("`user`.*").Joins("Join action_user ON action_user.user_id = user.id").Where("action_id = ?", action.ID).Find(&users)
		mydb.Select("`user_group`.*").Joins("Join action_user_group ON action_user_group.user_group_id = user_group.id").Where("action_id = ?", action.ID).Find(&user_groups)
		action.Users = users
		action.UserGroups = user_groups
	}
	strategy.Hosts = hosts
	strategy.Groups = groups
	strategy.Triggers = triggers
	strategy.Actions = actions

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategy": &strategy})
}

func strategyUpdate(c *gin.Context) {
	var strategy Strategy
	var strategy_db types.Strategy
	if err := c.BindJSON(&strategy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := mydb.Where("id = ?", strategy.ID).First(&strategy_db).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	if !GetUser(c).IsAdmin() && GetUser(c).ID != strategy_db.UserID {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusForbidden, "message": "not allow"})
		return
	}

	tx := mydb.Begin()

	if err := tx.Where("strategy_id = ?", strategy.ID).Delete(&strategy.Hosts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	if err := mydb.Where("strategy_id = ?", strategy.ID).Delete(&strategy.Groups).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	if err := mydb.Where("strategy_id = ?", strategy.ID).Delete(&strategy.Triggers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	if err := mydb.Where("strategy_id = ?", strategy.ID).Delete(&strategy.Actions).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	if err := tx.Save(&strategy).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategy": &strategy})
}

func strategySwitch(c *gin.Context) {
	strategy_id := c.Param("id")
	enable, err := strconv.Atoi(c.DefaultQuery("enable", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	var strategy StrategyWithName
	if len(strategy_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "strategy_id should be applied"})
		return
	}

	if err := mydb.Table("strategy s").Select("s.*, user.username user_name, ps.name parent_name").
		Joins("Left Join user ON s.user_id = user.id").
		Joins("Left Join strategy ps ON s.pid = ps.id").Where("s.id = ?", strategy_id).Find(&strategy).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	if !GetUser(c).IsAdmin() && GetUser(c).ID != strategy.UserID {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusForbidden, "message": "not allow"})
		return
	}

	strategy.Enable = enable != 0
	if err := mydb.Select("enable").Save(&strategy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategy": &strategy})
}
