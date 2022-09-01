package service

import (
	"google.golang.org/grpc/status"
	"io"
	"owl/common/logger"
	"owl/common/utils"
	proxyProto "owl/proxy/proto"
	"path"
	"path/filepath"
)

func (proxySrv *OwlProxyService) DownloadPluginFile(req *proxyProto.DownloadPluginReq, stream proxyProto.OwlProxyService_DownloadPluginFileServer) error {
	proxySrv.logger.Debug("proxySrv.DownloadPluginFile called.")
	defer proxySrv.logger.Debug("proxySrv.DownloadPluginFile end.")

	pluginPathname, err := filepath.Abs(path.Join(proxySrv.conf.PluginDir, req.RelPath))
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"plugin_dir": proxySrv.conf.PluginDir,
			"rel_path":   req.RelPath,
			"error":      err,
		}, "An error occurred while filepath.Abs in proxySrv.DownloadPluginFile.")
	}

	proxySrv.logger.InfoWithFields(logger.Fields{
		"plugin_pathname": pluginPathname,
	}, "proxySrv.DownloadPluginFile prepare send plugin file.")
	err = proxySrv.grpcDownloader.Download(pluginPathname, func(buffer []byte) error {
		return stream.Send(&proxyProto.PluginFile{Buffer: buffer})
	})
	if err != nil {
		switch err {
		case utils.ErrEndOfFileExit:
			return status.Error(utils.DefaultDownloaderEndOfFileExitCode, io.EOF.Error())
		case utils.ErrNormallyExit:
			return status.Error(utils.DefaultDownloaderNormallyExitCode, "normally exit")
		}
	}

	return err
}