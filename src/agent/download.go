package main

import (
	"context"
	"google.golang.org/grpc/status"
	"os"
	"owl/common/logger"
	"owl/common/utils"
	proxypb "owl/proxy/proto"
	"path"
)

// downloadPluginFile 下载指定插件文件
func (a *agent) downloadPluginFile(relPath string, pathname string) error {
	a.logger.Info("agent.downloadPluginFile called.")
	defer a.logger.Info("agent.downloadPluginFile end.")

	ctx, cancel := context.WithTimeout(a.ctx, a.conf.DownloadPluginTimeoutSecs)
	defer cancel()

	stream, err := a.proxyCli.DownloadPluginFile(ctx, &proxypb.DownloadPluginReq{RelPath: relPath})
	if err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling agent.proxyCli.DownloadPluginFile.")
		return err
	}

	_ = os.MkdirAll(path.Dir(pathname), 0755)
	fp, err := os.OpenFile(pathname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"rel_path": relPath,
			"pathname": pathname,
			"error":    err,
		}, "An error occurred while calling os.OpenFile.")
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
				a.logger.InfoWithFields(logger.Fields{
					"rel_path":       relPath,
					"pathname":       pathname,
					"status_code":    sts.Code(),
					"status_message": sts.Message(),
				}, "agent.downloadPluginFile success by EOF status.Code.")
				return nil
			}
			a.logger.ErrorWithFields(logger.Fields{
				"rel_path": relPath,
				"pathname": pathname,
				"error":    err,
			}, "An error occurred while calling stream.Recv.")
			return err
		}

		br, err := fp.Write(rsp.Buffer)
		if err != nil {
			a.logger.ErrorWithFields(logger.Fields{
				"rel_path": relPath,
				"pathname": pathname,
				"error":    err,
			}, "An error occurred while calling fp.Write.")
			return err
		}
		a.logger.DebugWithFields(logger.Fields{
			"rel_path":       relPath,
			"pathname":       pathname,
			"bytes_received": br,
		}, "agent.downloadPluginFile received some data.")
	}
}
