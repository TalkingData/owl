package main

import (
	"sync"
	"time"

	"owl/common/tcp"
	"owl/common/types"
	"owl/controller/cache"
)

var controller *Controller

type Controller struct {
	*tcp.Server
	taskPool         *TaskPool
	resultPool       *ResultPool
	nodePool         *NodePool
	eventQueues      map[int]*Queue
	eventQueuesMutex *sync.RWMutex
	statusCache      *cache.Cache
	taskCache        *cache.Cache
}

func InitController() error {
	controllerServer := tcp.NewServer(GlobalConfig.TCP_BIND, &ControllerHandle{})
	controllerServer.SetMaxPacketSize(uint32(GlobalConfig.MAX_PACKET_SIZE))
	if err := controllerServer.ListenAndServe(); err != nil {
		return err
	}
	lg.Info("Start listen: %v", GlobalConfig.TCP_BIND)

	controller = &Controller{controllerServer,
		NewTaskPool(GlobalConfig.TASK_POOL_SIZE),
		NewResultPool(GlobalConfig.RESULT_POOL_SIZE),
		NewNodePool(),
		generateQueues(),
		new(sync.RWMutex),
		cache.New(time.Duration(GlobalConfig.LOAD_STRATEGIES_INTERVAL)*time.Second, 30*time.Second),
		cache.New(10*time.Minute, 10*time.Minute)}

	go controller.loadStrategiesForever()
	go controller.processStrategyResultForever()
	go controller.processStrategyEventForever()
	go controller.checkNodesForever()
	go controller.startHttpServer()

	return nil
}

// generateQueues 生成产品线报警队列
func generateQueues() map[int]*Queue {
	pq := make(map[int]*Queue)
	for _, product := range mydb.GetProducts() {
		pq[product.ID] = NewQueue(0)
	}
	return pq
}

//checkNodesForever 持续运行检查节点函数，并维护节点数组
func (c *Controller) checkNodesForever() {
	for {
		time_now := time.Now()
		for ip, node := range c.nodePool.Nodes {
			if time_now.Sub(node.Update) > time.Duration(time.Second*10) {
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
			c.eventQueues[product.ID] = NewQueue(0)
		}
		c.eventQueues[product.ID].update_time = now
	}
	for p, q := range c.eventQueues {
		if now.Sub(q.update_time).Minutes() > 10 {
			delete(c.eventQueues, p)
		}
	}
}
