package biz

import (
	"context"
	"errors"
	"owl/common/logger"
	"owl/common/orm"
	"owl/model"
	"strings"
)

func (b *Biz) RegisterAgent(
	ctx context.Context,
	hostId, ip, hostname, version string,
	uptime, idlePct float64,
	metadata map[string]string,
) error {
	// 准备更新或创建主机
	b.logger.InfoWithFields(logger.Fields{
		"agent_host_id":  hostId,
		"agent_ip":       ip,
		"agent_hostname": hostname,
		"agent_version":  version,
		"agent_uptime":   uptime,
		"agent_idle_pct": idlePct,
		"agent_metadata": metadata,
	}, "Biz.RegisterAgent prepare execute dao.SetOrNewHostById.")
	hostObj, err := b.dao.SetOrNewHostById(ctx, hostId, ip, hostname, version, uptime, idlePct)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  hostId,
			"agent_ip":       ip,
			"agent_hostname": hostname,
			"error":          err,
		}, "An error occurred while calling dao.SetOrNewHostById.")
		return err
	}

	if hostObj == nil || len(hostObj.Id) < 1 {
		err = errors.New("nil host object")
		b.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  hostId,
			"agent_ip":       ip,
			"agent_hostname": hostname,
			"error":          err,
		}, "Biz.RegisterAgent Got an nil host object from dao.SetOrNewHostById.")
		return err
	}

	// 产品线相关处理
	agentMetadataProductName := metadata["product"]
	productObj, err := b.productProcess(ctx, hostObj.Id, agentMetadataProductName)
	if err != nil {
		return err
	}
	// 如果产品线对象是空，无法执行后续操作，直接返回即可
	if productObj == nil {
		return nil
	}

	// 主机组相关处理
	agentMetadataGroupName := metadata["group"]
	if err = b.hostGroupProcess(ctx, hostObj.Id, agentMetadataGroupName, productObj.Id); err != nil {
		return err
	}

	return nil
}

func (b *Biz) productProcess(ctx context.Context, hostId, productName string) (*model.Product, error) {
	// 如果没有产品线，那么需要的操作已经完成，可以直接返回
	if productName == "" {
		b.logger.WarnWithFields(logger.Fields{
			"host_id": hostId,
		}, "Biz.productProcess got an nil product object, skipped it.")
		return nil, nil
	}

	// 修剪处理
	trimmedProductName := strings.TrimSpace(productName)

	var (
		productObj *model.Product
		err        error
	)
	// 根据配置文件，在数据库中创建或获取产品线
	if b.conf.AllowCreateProductAuto {
		// 由于配置文件需要自动新建不存在的产品线，则在数据库中获取产品线，如果不存在则新建
		b.logger.InfoWithFields(logger.Fields{
			"host_id":                   hostId,
			"trimmed_product_name":      trimmedProductName,
			"allow_create_product_auto": b.conf.AllowCreateProductAuto,
		}, "Biz.RegisterAgent prepare execute dao.GetOrNewProduct.")
		productObj, err = b.dao.GetOrNewProduct(
			ctx, trimmedProductName, "it's auto created by cfc service v6", "cfc_v6",
		)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"host_id":              hostId,
				"trimmed_product_name": trimmedProductName,
				"error":                err,
			}, "An error occurred while calling dao.GetOrNewProduct.")
			return nil, err
		}
	} else {
		// 由于配置文件不允许自动新建产品线，则只在数据库中查找产品线
		b.logger.InfoWithFields(logger.Fields{
			"host_id":                   hostId,
			"trimmed_product_name":      trimmedProductName,
			"allow_create_product_auto": b.conf.AllowCreateProductAuto,
		}, "Biz.RegisterAgent prepare execute dao.GetProduct.")
		productObj, err = b.dao.GetProduct(ctx, orm.Query{"name": trimmedProductName})
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"host_id":              hostId,
				"trimmed_product_name": trimmedProductName,
				"error":                err,
			}, "An error occurred while calling dao.GetProduct.")
			return nil, err
		}
	}

	// 如果产品线仍然是空，则返回
	if productObj == nil || productObj.Id < 1 {
		b.logger.WarnWithFields(logger.Fields{
			"host_id":                   hostId,
			"trimmed_product_name":      trimmedProductName,
			"allow_create_product_auto": b.conf.AllowCreateProductAuto,
		}, "Product object still nil, skipped it.")
		return nil, nil
	}

	b.logger.InfoWithFields(logger.Fields{
		"host_id":              hostId,
		"product_id":           productObj.Id,
		"trimmed_product_name": trimmedProductName,
	}, "Biz.RegisterAgent prepare execute dao.IsHostInProduct.")
	isHostInProd, err := b.dao.IsHostInProduct(ctx, productObj.Id, hostId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id":              hostId,
			"product_id":           productObj.Id,
			"trimmed_product_name": trimmedProductName,
			"error":                err,
		}, "An error occurred while calling dao.IsHostInProduct.")
	}
	if isHostInProd {
		return productObj, nil
	}

	b.logger.InfoWithFields(logger.Fields{
		"host_id":              hostId,
		"product_id":           productObj.Id,
		"trimmed_product_name": trimmedProductName,
	}, "Biz.RegisterAgent host not in product, prepare execute dao.NewProductHost.")
	if _, err = b.dao.NewProductHost(ctx, productObj.Id, hostId); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id":              hostId,
			"product_id":           productObj.Id,
			"trimmed_product_name": trimmedProductName,
			"error":                err,
		}, "An error occurred while calling dao.NewProductHost.")
		return nil, err
	}

	return productObj, nil
}

func (b *Biz) hostGroupProcess(ctx context.Context, hostId, groupName string, productId uint32) error {
	// 如果没有主机组，那么需要的操作已经完成，可以直接返回
	if groupName == "" {
		b.logger.WarnWithFields(logger.Fields{
			"host_id":    hostId,
			"product_id": productId,
		}, "Got nil group name, skipped it.")
		return nil
	}

	// 修剪处理
	trimmedGroupName := strings.TrimSpace(groupName)

	// 在数据库查询主机组
	b.logger.InfoWithFields(logger.Fields{
		"host_id":            hostId,
		"product_id":         productId,
		"trimmed_group_name": trimmedGroupName,
	}, "Biz.RegisterAgent prepare execute dao.GetOrNewHostGroup.")

	groupObj, err := b.dao.GetOrNewHostGroup(
		ctx, productId, groupName, "it's auto created by cfc service v6", "cfc_v6",
	)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id":            hostId,
			"product_id":         productId,
			"trimmed_group_name": trimmedGroupName,
			"error":              err,
		}, "An error occurred while calling dao.GetOrNewHostGroup.")
		return err
	}

	// 如果主机组仍然是空
	if groupObj == nil || groupObj.Id < 1 {
		b.logger.WarnWithFields(logger.Fields{
			"host_id":            hostId,
			"product_id":         productId,
			"trimmed_group_name": trimmedGroupName,
		}, "Group object still empty, RegisterAgentBiz.hostGroupProcess skipped by failure.")
		return nil
	}

	// 添加主机到主机组
	b.logger.InfoWithFields(logger.Fields{
		"host_id":            hostId,
		"product_id":         productId,
		"trimmed_group_name": trimmedGroupName,
	}, "Biz.RegisterAgent prepare execute dao.IsHostInHostGroup.")
	isHostInGroup, err := b.dao.IsHostInHostGroup(ctx, groupObj.Id, hostId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id":            hostId,
			"product_id":         productId,
			"trimmed_group_name": trimmedGroupName,
			"error":              err,
		}, "An error occurred while calling dao.IsHostInHostGroup.")
	}
	if isHostInGroup {
		return nil
	}

	b.logger.InfoWithFields(logger.Fields{
		"host_id":            hostId,
		"product_id":         productId,
		"trimmed_group_name": trimmedGroupName,
	}, "Biz.RegisterAgent host not in host group, prepare execute dao.NewHostGroupHost.")
	if _, err = b.dao.NewHostGroupHost(ctx, groupObj.Id, hostId); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id":            hostId,
			"product_id":         productId,
			"trimmed_group_name": trimmedGroupName,
			"error":              err,
		}, "An error occurred while calling dao.NewHostGroupHost.")
		return err
	}

	return nil
}
