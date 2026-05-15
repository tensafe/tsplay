//go:build windows

package tsplay_core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func configureLocalCDPBrowserCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func terminateLocalCDPBrowserCommand(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return taskkillLocalCDPBrowserCommand(cmd, false)
}

func killLocalCDPBrowserCommand(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return taskkillLocalCDPBrowserCommand(cmd, true)
}

func taskkillLocalCDPBrowserCommand(cmd *exec.Cmd, force bool) error {
	args := []string{"/PID", strconv.Itoa(cmd.Process.Pid), "/T"}
	if force {
		args = append(args, "/F")
	}
	output, err := exec.Command("taskkill", args...).CombinedOutput()
	if err == nil {
		return nil
	}
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		return nil
	}
	if killErr := cmd.Process.Kill(); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
		return fmt.Errorf("taskkill browser process tree: %w; output=%s; fallback kill: %v", err, strings.TrimSpace(string(output)), killErr)
	}
	return nil
}
