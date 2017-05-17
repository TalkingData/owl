package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Host struct {
	types.Host
	Groups     []types.Group `json:"groups"`
	Metrics    int           `json:"metrics"`
	Plugins    int           `json:"plugins"`
	Strategies int           `json:"strategies"`
}
type Metric struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	DT       string    `json:"dt"`
	Cycle    int       `json:"cycle"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func hostList(c *gin.Context) {
	hosts := []*Host{}
	response := gin.H{"code": http.StatusOK}
	db := mydb.Table("host").Order("status asc")

	if id := c.Query("group_id"); len(id) > 0 {
		group_id, _ := strconv.Atoi(id)
		if group_id == 0 {
			db = db.Joins("LEFT JOIN host_group ON host.id = host_group.host_id").Where("host_group.group_id is NULL")
		} else {
			db = db.Joins("JOIN host_group ON host.id = host_group.host_id").Where("host_group.group_id = ?", group_id)
		}
	}
	if status := c.Query("status"); len(status) > 0 {
		db = db.Where("status = ?", status)
	}
	q := c.Query("q")
	if len(q) > 0 {
		q := fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ? or hostname like ? or ip like ?", q, q, q)
	}
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))
	var total int
	db.Count(&total)
	if page != 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}

	//db.Count(&cnt)
	db.Find(&hosts)
	response["total"] = total
	if page != 0 {
		for _, h := range hosts {
			//获取主机关联的组
			mydb.Joins("JOIN host_group ON host_group.group_id = group.id").
				Where("host_group.host_id = ?", h.ID).Find(&h.Groups)

			//获取metric数量
			mydb.Table("metric").Where("host_id = ?", h.ID).Count(&h.Metrics)

			//获取插件数量
			mydb.Table("host_plugin").Where("host_id = ?", h.ID).Count(&h.Plugins)

			//获取主机的策略数量
			strategies := getStrategiesByHostID(h.ID)
			h.Strategies = len(strategies)
		}
	}
	response["hosts"] = hosts
	c.JSON(http.StatusOK, response)
}

func hostInfo(c *gin.Context) {
	var (
		host_cnt, metric_cnt, group_cnt int
	)
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	group_id, _ := strconv.Atoi(c.DefaultQuery("group_id", "-1"))

	if group_id == -1 {
		//all
		mydb.Table("host").Count(&host_cnt)
		mydb.Table("metric").Count(&metric_cnt)
		mydb.Table("group").Count(&group_cnt)
		response["groups"] = group_cnt
	} else if group_id == 0 {
		//未分组
		mydb.Table("host").Joins("LEFT JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id is NULL").Count(&host_cnt).
			Joins("JOIN metric ON host.id = metric.host_id").Count(&metric_cnt)
		//mydb.Table("host").Count(&host_cnt)
		//mydb.Table("metric").Count(&metric_cnt)
		//mydb.Table("group").Count(&group_cnt)
		//response["groups"] = group_cnt
	} else {
		mydb.Table("host").Joins("JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id = ?", group_id).
			Count(&host_cnt)
		mydb.Table("host").Joins("JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id = ?", group_id).
			Joins("JOIN metric ON host.id = metric.host_id").
			Count(&metric_cnt)
	}
	response["hosts"] = host_cnt
	response["metrics"] = metric_cnt
}

func hostStatus(c *gin.Context) {
	var (
		normal  int
		failed  int
		disable int
		pending int
	)
	group_id, _ := strconv.Atoi(c.DefaultQuery("group_id", "-1"))
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	db := mydb.Table("host")
	if group_id == -1 {
		db.Where("status=?", "0").Count(&failed)
		db.Where("status=?", "1").Count(&normal)
		db.Where("status=?", "2").Count(&disable)
		db.Where("status=?", "3").Count(&pending)
	} else if group_id == 0 {
		db = db.Joins("LEFT JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id is NULL")
		db.Where("host.status=?", "0").Count(&failed)
		db.Where("host.status=?", "1").Count(&normal)
		db.Where("host.status=?", "2").Count(&disable)
		db.Where("host.status=?", "3").Count(&pending)
	} else {
		db = db.Joins("JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id = ?", group_id)
		db.Where("host.status = ?", "0").
			Count(&failed)
		db.Where("host.status = ?", "1").
			Count(&normal)
		db.Where("host.status = ?", "3").
			Count(&pending)
		db.Where("host.status = ?", "2").
			Count(&disable)
	}
	response["pending"] = pending
	response["normal"] = normal
	response["failed"] = failed
	response["disable"] = disable
}

func metricsByHostId(c *gin.Context) {
	host_id := c.Param("id")
	var cnt int
	mydb.Table("host").Where("id = ?", host_id).Count(&cnt)
	if cnt == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "message": "the specified host is not found"})
		return
	}

	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": "invalid parameters"})
		return
	}
	db := mydb.Table("metric").Where("host_id = ?", host_id)
	q := c.Query("q")
	if len(q) > 0 {
		q := fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ?", q)
	}
	var total int
	db.Count(&total)
	offset := (page - 1) * pageSize

	metrics := []Metric{}
	db.Offset(offset).Limit(pageSize).Find(&metrics)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "metrics": metrics, "total": total})
}

func pluginByHostID(c *gin.Context) {
	host_id := c.Param("id")
	var cnt int
	mydb.Table("host").Where("id = ?", host_id).Count(&cnt)
	if cnt == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "message": "the specified host is not found"})
		return
	}

	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": "invalid parameters"})
		return
	}
	db := mydb.Table("plugin").
		Joins("LEFT JOIN host_plugin ON plugin.id = host_plugin.plugin_id").
		Where("host_plugin.host_id = ?", host_id)
	q := c.Query("q")
	if len(q) > 0 {
		q := fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ?", q)
	}
	var total int
	db.Count(&total)
	offset := (page - 1) * pageSize

	plugins := []types.Plugin{}
	db.Offset(offset).Limit(pageSize).Find(&plugins)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "plugins": plugins, "total": total})

}

func strategiesByHostId(c *gin.Context) {
	host_id := c.Param("id")

	strategies := getStrategiesByHostID(host_id)
	response := gin.H{"strategies": make([]gin.H, 0)}
	for _, strategy := range strategies {
		strategy_slice := response["strategies"].([]gin.H)
		response["strategies"] = append(strategy_slice, gin.H{"id": strategy.ID, "name": strategy.Name})
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "response": response})
}

func getStrategiesByHostID(host_id string) []types.Strategy {
	var strategies []types.Strategy
	var temp_strategies []types.Strategy
	var unique_strategies []types.Strategy = make([]types.Strategy, 0)

	mydb.Joins("JOIN strategy_group ON strategy_group.strategy_id = strategy.id").
		Joins("JOIN host_group ON host_group.group_id = strategy_group.group_id").
		Where("host_group.host_id = ?", host_id).Find(&temp_strategies)
	strategies = append(strategies, temp_strategies...)
	mydb.Joins("JOIN strategy_host ON strategy_host.strategy_id = strategy.id").
		Where("strategy_host.host_id = ?", host_id).Find(&temp_strategies)
	strategies = append(strategies, temp_strategies...)
	ids := make(map[int]bool)
	for _, strategy := range strategies {
		if strategy.Enable == false {
			continue
		}
		if _, ok := ids[strategy.ID]; ok {
			continue
		}
		ids[strategy.ID] = true
		unique_strategies = append(unique_strategies, strategy)
	}
	return unique_strategies
}

func hostDelete(c *gin.Context) {
	id := c.Param("id")
	host := types.Host{}
	if mydb.Table("host").Where("id =? or sn =? ", id, id).First(&host).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "message": "host not found"})
		return
	}
	mydb.Delete(&host)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "delete success",
	})
}

func hostRename(c *gin.Context) {
	id := c.Param("id")
	host := types.Host{}
	mydb.Table("host").Where("id = ?", id).First(&host)
	if mydb.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": mydb.Error})
		return
	}
	h := types.Host{}
	if err := c.BindJSON(&h); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}
	host.Name = h.Name
	mydb.Save(&host)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "rename the host success"})
}

func hostEnable(c *gin.Context) {
	id := c.Param("id")
	host := types.Host{}
	if mydb.Table("host").Where("id =? or sn =? ", id, id).First(&host).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "message": "host not found"})
		return
	}
	mydb.Model(&host).UpdateColumn("status", "1")
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "enable the host success"})
}

func hostDisable(c *gin.Context) {
	id := c.Param("id")
	host := types.Host{}
	if mydb.Table("host").Where("id =? or sn =? ", id, id).First(&host).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "message": "host not found"})
		return
	}
	mydb.Model(&host).UpdateColumn("status", "2")
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "host": "disable the host success"})
}

func (this Metric) MarshalJSON() ([]byte, error) {
	type Alias Metric
	return json.Marshal(&struct {
		CreateAt string `json:"create_at"`
		UpdateAt string `json:"update_at"`
		Alias
	}{
		CreateAt: this.CreateAt.Format("2006-01-02 15:04:05"),
		UpdateAt: this.UpdateAt.Format("2006-01-02 15:04:05"),
		Alias:    (Alias)(this),
	})
}
