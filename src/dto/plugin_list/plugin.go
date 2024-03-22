package builtin

import (
	"context"
	"fmt"
	"owl/common/utils"
	"strings"
	"sync"
	"time"
)

type PluginTaskFunc func(ctx context.Context, tx int64, cycle int32, command string, args ...string)

type Plugin struct {
	Id        uint32
	Name      string
	LocalPath string
	Checksum  string
	Args      []string
	Interval  int32
	Timeout   int32

	mu sync.RWMutex

	execUntrusted bool

	taskFunc PluginTaskFunc

	parentCtx context.Context

	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewPlugin 新建Plugin
func NewPlugin(
	ctx context.Context,
	id uint32,
	name, localPath, checksum string,
	args []string,
	interval, timeout int32,
	execUntrusted bool,
	taskFunc PluginTaskFunc,
) *Plugin {

	return &Plugin{
		Id:        id,
		Name:      name,
		LocalPath: localPath,
		Checksum:  checksum,
		Args:      args,
		Interval:  interval,
		Timeout:   timeout,

		execUntrusted: execUntrusted,

		parentCtx: ctx,

		taskFunc: taskFunc,
	}
}

// StartTask 运行插件采集任务
func (p *Plugin) StartTask() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// ctx 或 cancelFunc 不是nil的，是已经启动了采集任务的插件
	if p.ctx != nil || p.cancelFunc != nil {
		return
	}

	// 如果设置是不允许运行不信任的插件，并且插件校验和与插件文件校验和不一致的情况，则不会启动采集任务
	if p.execUntrusted && !p.IsValidChecksum() {
		return
	}

	p.ctx, p.cancelFunc = context.WithCancel(p.parentCtx)
	go func() {
		tk := time.Tick(time.Second * time.Duration(p.Interval))
		for {
			select {
			case c := <-tk:
				tkCtx, tkCancelFunc := context.WithTimeout(p.ctx, time.Second*time.Duration(p.Timeout))
				p.taskFunc(tkCtx, c.Unix(), p.Interval, p.LocalPath, p.Args...)
				tkCancelFunc()
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

// StopTask 结束插件采集任务
func (p *Plugin) StopTask() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// ctx 或 cancelFunc 为nil的，是没有启动采集任务的插件
	if p.ctx == nil && p.cancelFunc == nil {
		return
	}

	p.cancelFunc()
	p.ctx = nil
	p.cancelFunc = nil
}

// IsValidChecksum 判断插件校验和是否与本地一致
func (p *Plugin) IsValidChecksum() bool {
	return p.Checksum == p.GetFileChecksum()
}

// GetFileChecksum 获得插件文件的校验
func (p *Plugin) GetFileChecksum() string {
	cs, err := utils.GetFileMD5(p.LocalPath)
	if err != nil {
		return ""
	}

	return cs
}

func (p *Plugin) GetPk() string {
	return fmt.Sprintf(
		"%s.%d.%d.%s",
		p.LocalPath,
		p.Interval,
		p.Timeout,
		strings.Join(p.Args, ","),
	)
}
