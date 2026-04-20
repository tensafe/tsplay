package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	defaultScreenRecordInput  = "Capture screen 0:none"
	defaultScreenRecordOutput = "artifacts/recordings/tsplay-tutorial.mp4"
)

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)
var avfoundationDevicePattern = regexp.MustCompile(`\[[0-9]+\]\s+(.*)$`)

type screenRecordOptions struct {
	InputSpec     string
	OutputPath    string
	Command       string
	Shell         string
	FrameRate     int
	VideoSize     string
	CaptureCursor bool
	Warmup        time.Duration
	Cooldown      time.Duration
	MaxDuration   time.Duration
	CRF           int
	Preset        string
}

type screenRecordResult struct {
	FFmpegPath        string `json:"ffmpeg_path"`
	InputSpec         string `json:"input_spec"`
	OutputPath        string `json:"output_path"`
	Command           string `json:"command,omitempty"`
	CommandExitCode   int    `json:"command_exit_code,omitempty"`
	CommandSucceeded  bool   `json:"command_succeeded"`
	FrameRate         int    `json:"frame_rate"`
	VideoSize         string `json:"video_size,omitempty"`
	CaptureCursor     bool   `json:"capture_cursor"`
	WarmupMs          int64  `json:"warmup_ms"`
	CooldownMs        int64  `json:"cooldown_ms"`
	MaxDurationMs     int64  `json:"max_duration_ms,omitempty"`
	RecordingStarted  string `json:"recording_started_at"`
	RecordingFinished string `json:"recording_finished_at"`
}

type screenRecordDeviceProbe struct {
	FFmpegPath      string   `json:"ffmpeg_path"`
	VideoDevices    []string `json:"video_devices,omitempty"`
	AudioDevices    []string `json:"audio_devices,omitempty"`
	PermissionHint  string   `json:"permission_hint,omitempty"`
	RawOutput       string   `json:"raw_output,omitempty"`
	ProbeExitStatus string   `json:"probe_exit_status,omitempty"`
}

func listScreenRecordDevices() (*screenRecordDeviceProbe, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("screen recording probe currently only supports macOS")
	}
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	cmd := exec.Command(ffmpegPath, "-f", "avfoundation", "-list_devices", "true", "-i", "")
	output, runErr := cmd.CombinedOutput()
	raw := strings.TrimSpace(stripANSI(string(output)))
	videoDevices, audioDevices := parseAVFoundationDevices(raw)

	probe := &screenRecordDeviceProbe{
		FFmpegPath:   ffmpegPath,
		VideoDevices: videoDevices,
		AudioDevices: audioDevices,
	}
	if raw != "" {
		probe.RawOutput = raw
	}
	if runErr != nil {
		probe.ProbeExitStatus = runErr.Error()
	}
	if len(videoDevices) == 0 {
		probe.PermissionHint = "没有探测到可用的视频采集设备。请先确认 Terminal/Codex 已被授予系统的屏幕录制权限，然后重试。"
	}

	if runErr != nil && !strings.Contains(raw, "AVFoundation video devices:") {
		return probe, fmt.Errorf("probe avfoundation devices: %w", runErr)
	}
	return probe, nil
}

func runScreenRecordAction(options screenRecordOptions) (*screenRecordResult, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("screen recording currently only supports macOS")
	}
	if strings.TrimSpace(options.InputSpec) == "" {
		return nil, fmt.Errorf("-record-input cannot be empty")
	}
	if strings.TrimSpace(options.OutputPath) == "" {
		return nil, fmt.Errorf("-record-output cannot be empty")
	}
	if options.FrameRate <= 0 {
		return nil, fmt.Errorf("-record-fps must be greater than 0")
	}
	if options.CRF < 0 {
		return nil, fmt.Errorf("-record-crf must be greater than or equal to 0")
	}
	if strings.TrimSpace(options.Preset) == "" {
		return nil, fmt.Errorf("-record-preset cannot be empty")
	}
	if strings.TrimSpace(options.Shell) == "" {
		options.Shell = "/bin/zsh"
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	outputPath, err := filepath.Abs(options.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("resolve record output %q: %w", options.OutputPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return nil, fmt.Errorf("create record output dir: %w", err)
	}

	recorderArgs := buildScreenRecordFFmpegArgs(options, outputPath)
	recorderCmd := exec.Command(ffmpegPath, recorderArgs...)
	recorderCmd.Stdout = os.Stdout
	recorderCmd.Stderr = os.Stderr
	recorderStdin, err := recorderCmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open ffmpeg stdin: %w", err)
	}
	if err := recorderCmd.Start(); err != nil {
		return nil, fmt.Errorf("start ffmpeg recorder: %w", err)
	}

	recorderDone := make(chan error, 1)
	go func() {
		recorderDone <- recorderCmd.Wait()
	}()

	result := &screenRecordResult{
		FFmpegPath:       ffmpegPath,
		InputSpec:        options.InputSpec,
		OutputPath:       outputPath,
		Command:          strings.TrimSpace(options.Command),
		FrameRate:        options.FrameRate,
		VideoSize:        strings.TrimSpace(options.VideoSize),
		CaptureCursor:    options.CaptureCursor,
		WarmupMs:         options.Warmup.Milliseconds(),
		CooldownMs:       options.Cooldown.Milliseconds(),
		MaxDurationMs:    options.MaxDuration.Milliseconds(),
		RecordingStarted: time.Now().Format(time.RFC3339),
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	if err := waitWithRecorderGuard(options.Warmup, recorderDone); err != nil {
		return result, err
	}

	if strings.TrimSpace(options.Command) == "" {
		fmt.Fprintln(os.Stderr, "screen recording started; press Ctrl-C to stop")
		select {
		case sig := <-signals:
			_ = sig
		case err := <-recorderDone:
			if err != nil {
				return result, fmt.Errorf("ffmpeg recorder exited unexpectedly: %w", err)
			}
			return result, fmt.Errorf("ffmpeg recorder exited unexpectedly")
		}
		if err := stopScreenRecorder(recorderCmd, recorderStdin, recorderDone); err != nil {
			return result, err
		}
		result.CommandSucceeded = true
		result.RecordingFinished = time.Now().Format(time.RFC3339)
		return result, nil
	}

	commandCmd := exec.Command(options.Shell, "-lc", options.Command)
	commandCmd.Stdout = os.Stdout
	commandCmd.Stderr = os.Stderr
	attachProcessGroup(commandCmd)
	if err := commandCmd.Start(); err != nil {
		_ = stopScreenRecorder(recorderCmd, recorderStdin, recorderDone)
		return result, fmt.Errorf("start recorded command: %w", err)
	}

	commandDone := make(chan error, 1)
	go func() {
		commandDone <- commandCmd.Wait()
	}()

	commandExitCode := 0
	commandErr := error(nil)
	recorderExited := false

	select {
	case sig := <-signals:
		_ = sig
		commandErr = fmt.Errorf("recording interrupted")
		terminateProcessGroup(commandCmd)
	case err := <-commandDone:
		commandErr = err
	case err := <-recorderDone:
		recorderExited = true
		commandErr = fmt.Errorf("ffmpeg recorder exited unexpectedly: %w", err)
		terminateProcessGroup(commandCmd)
	}

	if options.Cooldown > 0 {
		if err := waitWithRecorderGuard(options.Cooldown, recorderDone); err != nil {
			recorderExited = true
			if commandErr == nil {
				commandErr = err
			}
		}
	}
	if !recorderExited {
		if err := stopScreenRecorder(recorderCmd, recorderStdin, recorderDone); err != nil && commandErr == nil {
			commandErr = err
		}
	}

	if commandCmd.ProcessState == nil {
		select {
		case <-time.After(500 * time.Millisecond):
		case <-commandDone:
		}
	}
	if commandCmd.ProcessState != nil {
		commandExitCode = commandCmd.ProcessState.ExitCode()
	}
	result.CommandExitCode = commandExitCode
	result.CommandSucceeded = commandErr == nil
	result.RecordingFinished = time.Now().Format(time.RFC3339)

	if commandErr != nil {
		return result, commandErr
	}
	return result, nil
}

func buildScreenRecordFFmpegArgs(options screenRecordOptions, outputPath string) []string {
	args := []string{
		"-y",
		"-hide_banner",
		"-f", "avfoundation",
		"-capture_cursor", boolToFFmpegFlag(options.CaptureCursor),
		"-framerate", strconv.Itoa(options.FrameRate),
	}
	if strings.TrimSpace(options.VideoSize) != "" {
		args = append(args, "-video_size", strings.TrimSpace(options.VideoSize))
	}
	if options.MaxDuration > 0 {
		args = append(args, "-t", formatFFmpegSeconds(options.MaxDuration))
	}
	args = append(args,
		"-i", strings.TrimSpace(options.InputSpec),
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-preset", strings.TrimSpace(options.Preset),
		"-crf", strconv.Itoa(options.CRF),
		"-movflags", "+faststart",
		outputPath,
	)
	return args
}

func waitWithRecorderGuard(delay time.Duration, recorderDone <-chan error) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case err := <-recorderDone:
		if err != nil {
			return fmt.Errorf("ffmpeg recorder exited unexpectedly: %w", err)
		}
		return fmt.Errorf("ffmpeg recorder exited unexpectedly")
	}
}

func stopScreenRecorder(recorderCmd *exec.Cmd, recorderStdin io.WriteCloser, recorderDone <-chan error) error {
	if recorderCmd == nil {
		return nil
	}
	if recorderStdin != nil {
		_, _ = io.WriteString(recorderStdin, "q\n")
		_ = recorderStdin.Close()
	}

	timer := time.NewTimer(8 * time.Second)
	defer timer.Stop()

	select {
	case err := <-recorderDone:
		return normalizeScreenRecorderExit(err)
	case <-timer.C:
		if recorderCmd.Process != nil {
			_ = recorderCmd.Process.Signal(os.Interrupt)
		}
	}

	timer.Reset(3 * time.Second)
	select {
	case err := <-recorderDone:
		return normalizeScreenRecorderExit(err)
	case <-timer.C:
		if recorderCmd.Process != nil {
			_ = recorderCmd.Process.Kill()
		}
	}

	return normalizeScreenRecorderExit(<-recorderDone)
}

func normalizeScreenRecorderExit(err error) error {
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if !strings.Contains(err.Error(), "exit status") {
		return err
	}
	if ok := errorAsExit(err, &exitErr); ok && exitErr.ProcessState != nil && exitErr.ProcessState.ExitCode() == 0 {
		return nil
	}
	return err
}

func errorAsExit(err error, target **exec.ExitError) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	*target = exitErr
	return true
}

func attachProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || runtime.GOOS == "windows" {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func terminateProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	if runtime.GOOS == "windows" {
		_ = cmd.Process.Kill()
		return
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		_ = cmd.Process.Kill()
	}
}

func parseAVFoundationDevices(output string) ([]string, []string) {
	cleaned := stripANSI(output)
	lines := strings.Split(cleaned, "\n")
	videoDevices := []string{}
	audioDevices := []string{}
	section := ""
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		switch {
		case strings.Contains(line, "AVFoundation video devices:"):
			section = "video"
		case strings.Contains(line, "AVFoundation audio devices:"):
			section = "audio"
		case strings.Contains(line, "Error opening input"):
			section = ""
		case strings.Contains(line, "Error opening input file"):
			section = ""
		case strings.Contains(line, "Error opening input files"):
			section = ""
		default:
			name := parseAVFoundationDeviceName(line)
			if name == "" {
				continue
			}
			if section == "video" {
				videoDevices = append(videoDevices, name)
			}
			if section == "audio" {
				audioDevices = append(audioDevices, name)
			}
		}
	}
	return videoDevices, audioDevices
}

func parseAVFoundationDeviceName(line string) string {
	if line == "" {
		return ""
	}
	matches := avfoundationDevicePattern.FindStringSubmatch(line)
	if len(matches) != 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func stripANSI(value string) string {
	return ansiEscapePattern.ReplaceAllString(value, "")
}

func boolToFFmpegFlag(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func formatFFmpegSeconds(value time.Duration) string {
	seconds := value.Seconds()
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", seconds), "0"), ".")
}
