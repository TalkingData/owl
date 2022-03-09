package main

import (
	"sync"
	"time"

	"owl/common/types"
	"owl/controller/cache"

	"github.com/wuyingsong/tcp"
)

var controller *Controller

type Controller struct {
	*tcp.AsyncTCPServer
	taskPool         *TaskPool
	resultPool       *ResultPool
	nodePool         *NodePool
	eventQueues      map[int]*EventPool
	eventQueuesMutex *sync.RWMutex
	statusCache      *cache.Cache
	taskCache        *cache.Cache
}

func InitController() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MAX_PACKET_SIZE))
	server := tcp.NewAsyncTCPServer(GlobalConfig.TCP_BIND, &callback{}, protocol)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	lg.Info("Start listen: %v", GlobalConfig.TCP_BIND)

	controller = &Controller{server,
		NewTaskPool(GlobalConfig.TASK_POOL_SIZE),
		NewResultPool(GlobalConfig.RESULT_POOL_SIZE),
		NewNodePool(),
		make(map[int]*EventPool),
		new(sync.RWMutex),
		cache.New(time.Duration(GlobalConfig.LOAD_STRATEGIES_INTERVAL)*time.Second, 30*time.Second),
		cache.New(10*time.Minute, 10*time.Minute)}

	go controller.processStrategyResultForever()
	go controller.processStrategyEventForever()
	go controller.checkNodesForever()
	go controller.startHttpServer()
	go controller.loadStrategiesForever()

	return nil
}

//checkNodesForever 持续运行检查节点函数，并维护节点数组
func (c *Controller) checkNodesForever() {
	for {
		now := time.Now()
		for ip, node := range c.nodePool.Nodes {
			if now.Sub(node.Update) > time.Duration(time.Second*10) {
				delete(c.nodePool.Nodes, ip)
				lg.Warn("Inspector %v, %v lost from controller", ip, node.Hostname)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

//receiveHearbeat 接收心跳包，并更新节点状态
func (c *Controller) receiveHearbeat(heartbeat *types.HeartBeat) {
	controller.nodePool.Lock.Lock()
	defer controller.nodePool.Lock.Unlock()

	if node, ok := c.nodePool.Nodes[heartbeat.IP]; ok {
		node.IP = heartbeat.IP
		node.Hostname = heartbeat.Hostname
		node.Update = time.Now()
	} else {
		node = &types.Node{}
		node.IP = heartbeat.IP
		node.Hostname = heartbeat.Hostname
		node.Update = time.Now()
		c.nodePool.Nodes[heartbeat.IP] = node
	}
}

//refreshQueue 同步产品线的报警队列
func (c *Controller) refreshQueue(products []*types.Product) {
	c.eventQueuesMutex.Lock()
	defer c.eventQueuesMutex.Unlock()
	now := time.Now()
	for _, product := range products {
		if _, ok := c.eventQueues[product.ID]; !ok {
			lg.Info("create event queue %s", product.Name)
			eventPool := NewEventPool(product.Name, GlobalConfig.EVENT_POOL_SIZE)
			c.eventQueues[product.ID] = eventPool
			go processSingleQueue(eventPool)
		}
		lg.Info("refresh event queue %s", product.Name)
		c.eventQueues[product.ID].update_time = now
	}
	for p, q := range c.eventQueues {
		if now.Sub(q.update_time).Minutes() > 10 {
			lg.Warn("product id %s event queue update time is expires, delete", p)
			delete(c.eventQueues, p)
		}
	}
}
