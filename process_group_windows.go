//go:build windows

package main

import "os/exec"

func attachProcessGroup(cmd *exec.Cmd) {
}

func terminateProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}
