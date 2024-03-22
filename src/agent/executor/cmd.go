package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"owl/common/logger"
	"owl/dto"
)

func (e *Executor) ExecCollectCmd(ctx context.Context, ts int64, command string, args ...string) dto.TsDataArray {
	e.logger.InfoWithFields(logger.Fields{
		"command": command,
		"args":    args,
	}, "Executor.ExecCollectCmd called.")
	defer e.logger.Info("Executor.ExecCollectCmd end.")

	out := bytes.Buffer{}
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"command": command,
			"args":    args,
			"error":   err,
		}, "An error occurred while calling Executor.ExecCollectCmd.")
		return nil
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		err = cmd.Process.Kill()
		if err != nil {
			e.logger.WarnWithFields(logger.Fields{
				"command": command,
				"args":    args,
				"error":   err,
			}, "An error occurred while calling Executor.ExecCollectCmd.")
			return nil
		}

		e.logger.WarnWithFields(logger.Fields{
			"command":       command,
			"args":          args,
			"context_error": ctx.Err(),
		}, "Executor.ExecCollectCmd end by context.Done.")
		return nil

	case err = <-done:
		if err != nil {
			return nil
		}

		res := dto.TsDataArray{}
		if err = json.Unmarshal(out.Bytes(), &res); err != nil {
			e.logger.ErrorWithFields(logger.Fields{
				"command": command,
				"args":    args,
				"out":     out.String(),
				"error":   err,
			}, "An error occurred while calling Executor.ExecCollectCmd.")
			return nil
		}

		// 为每个TsData赋Timestamp
		for _, r := range res {
			r.Timestamp = ts
		}

		return res
	}
}
