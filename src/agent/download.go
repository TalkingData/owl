package main

import (
	"context"
	"google.golang.org/grpc/status"
	"os"
	"owl/common/logger"
	"owl/common/utils"
	proxyProto "owl/proxy/proto"
	"path"
)

// downloadPluginFile 下载指定插件文件
func (agent *agent) downloadPluginFile(relPath string, pathname string) error {
	agent.logger.Info("agent.downloadPluginFile called.")
	defer agent.logger.Info("agent.downloadPluginFile end.")

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.DownloadPluginTimeoutSecs)
	defer cancel()

	stream, err := agent.proxyCli.DownloadPluginFile(ctx, &proxyProto.DownloadPluginReq{RelPath: relPath})
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.proxyCli.DownloadPluginFile in agent.downloadPluginFile.")
		return err
	}

	_ = os.MkdirAll(path.Dir(pathname), 0755)
	fp, err := os.OpenFile(pathname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"rel_path": relPath,
			"pathname": pathname,
			"error":    err,
		}, "An error occurred while os.OpenFile in agent.downloadPluginFile.")
		return err
	}
	defer func() {
		_ = fp.Close()
	}()

	// 开始分块循环接收文件
	for {
		rsp, err := stream.Recv()
		if err != nil {
			sts := status.Convert(err)
			if sts.Code() == utils.DefaultDownloaderEndOfFileExitCode {
				agent.logger.InfoWithFields(logger.Fields{
					"rel_path":       relPath,
					"pathname":       pathname,
					"status_code":    sts.Code(),
					"status_message": sts.Message(),
				}, "agent.downloadPluginFile success by EOF status.Code.")
				return nil
			}
			agent.logger.ErrorWithFields(logger.Fields{
				"rel_path": relPath,
				"pathname": pathname,
				"error":    err,
			}, "An error occurred while stream.Recv in agent.downloadPluginFile.")
			return err
		}

		br, err := fp.Write(rsp.Buffer)
		if err != nil {
			agent.logger.ErrorWithFields(logger.Fields{
				"rel_path": relPath,
				"pathname": pathname,
				"error":    err,
			}, "An error occurred while fp.Write in agent.downloadPluginFile.")
			return err
		}
		agent.logger.DebugWithFields(logger.Fields{
			"rel_path":       relPath,
			"pathname":       pathname,
			"bytes_received": br,
		}, "agent.downloadPluginFile received some data.")
	}
}
