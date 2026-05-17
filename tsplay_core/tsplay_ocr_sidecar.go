package tsplay_core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	goddddocrModeHTTP    = "http"
	goddddocrModeSidecar = "sidecar"
	goddddocrModeCLI     = "cli"

	defaultGoddddocrStartupTimeoutMS = 15000
	defaultGoddddocrCLITimeoutMS     = 30000
)

type goddddocrSidecar struct {
	BaseURL  string
	cmd      *exec.Cmd
	waitDone chan error
	close    sync.Once
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
}

func isGoddddocrAction(action string) bool {
	switch action {
	case "ocr_ready", "ocr_request", "ocr_detect", "ocr_slide_comparison", "ocr_slide_match":
		return true
	default:
		return false
	}
}

func normalizeGoddddocrMode(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return "", nil
	case "http", "service", "remote", "url":
		return goddddocrModeHTTP, nil
	case "sidecar", "managed", "managed_sidecar", "managed-sidecar", "server":
		return goddddocrModeSidecar, nil
	case "cli", "direct", "command", "process":
		return goddddocrModeCLI, nil
	default:
		return "", fmt.Errorf("goddddocr mode must be one of auto, http, sidecar, or cli")
	}
}

func staticGoddddocrMode(step FlowStep) string {
	value, ok := step.param("mode")
	if !ok || len(flowReferences(value)) > 0 {
		return ""
	}
	mode, err := normalizeGoddddocrMode(fmt.Sprint(value))
	if err != nil {
		return ""
	}
	return mode
}

func goddddocrActionMode(ctx *FlowContext, step FlowStep, allowCLI bool) (string, error) {
	value, ok, err := flowStepResolvedParam(ctx, step, "mode")
	if err != nil {
		return "", err
	}
	modeText := ""
	if ok {
		modeText = fmt.Sprint(value)
	}
	if strings.TrimSpace(modeText) == "" {
		modeText = os.Getenv("GODDDDOCR_MODE")
	}
	mode, err := normalizeGoddddocrMode(modeText)
	if err != nil {
		return "", err
	}
	if mode == goddddocrModeCLI && !allowCLI {
		return "", fmt.Errorf("%s does not support mode=cli; use ocr_request for direct CLI OCR", step.Action)
	}
	if mode != "" {
		return mode, nil
	}

	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(endpoint) != "" || strings.TrimSpace(os.Getenv("GODDDDOCR_URL")) != "" {
		return goddddocrModeHTTP, nil
	}
	if goddddocrProcessAllowed(ctx) {
		return goddddocrModeSidecar, nil
	}
	return goddddocrModeHTTP, nil
}

func goddddocrProcessAllowed(ctx *FlowContext) bool {
	return ctx == nil || ctx.Security == nil || ctx.Security.AllowProcess
}

func requireGoddddocrProcess(ctx *FlowContext, action string) error {
	if goddddocrProcessAllowed(ctx) {
		return nil
	}
	return fmt.Errorf("%s mode requires allow_process=true for trusted local process execution", action)
}

func goddddocrEndpointForAction(ctx *FlowContext, step FlowStep, normalizeEndpoint func(string) string, needDetection bool) (string, *goddddocrSidecar, error) {
	mode, err := goddddocrActionMode(ctx, step, false)
	if err != nil {
		return "", nil, err
	}
	if mode == goddddocrModeSidecar {
		sidecar, err := ctx.ensureGoddddocrSidecar(step, needDetection)
		if err != nil {
			return "", nil, err
		}
		return normalizeEndpoint(sidecar.BaseURL), sidecar, nil
	}
	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return "", nil, err
	}
	return normalizeEndpoint(endpoint), nil, nil
}

func (ctx *FlowContext) ensureGoddddocrSidecar(step FlowStep, needDetection bool) (*goddddocrSidecar, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%s cannot start goddddocr sidecar without a flow context", step.Action)
	}
	if err := requireGoddddocrProcess(ctx, step.Action); err != nil {
		return nil, err
	}
	executable, err := goddddocrSidecarExecutable(ctx, step)
	if err != nil {
		return nil, err
	}
	serverArgs, err := flowStepOptionalStringListParam(ctx, step, "server_args")
	if err != nil {
		return nil, err
	}
	if det, ok, err := goddddocrOptionalBoolParam(ctx, step, "det"); err != nil {
		return nil, err
	} else if ok && det {
		needDetection = true
	}

	key := goddddocrSidecarKey(executable, serverArgs, needDetection)
	if ctx.OCRSidecars == nil {
		ctx.OCRSidecars = map[string]*goddddocrSidecar{}
	}
	if sidecar := ctx.OCRSidecars[key]; sidecar != nil {
		return sidecar, nil
	}

	startupTimeoutMS, err := goddddocrStartupTimeoutMS(ctx, step)
	if err != nil {
		return nil, err
	}
	port, err := freeLocalTCPPort()
	if err != nil {
		return nil, err
	}
	sidecar, err := startGoddddocrSidecar(ctx, executable, serverArgs, needDetection, port, time.Duration(startupTimeoutMS)*time.Millisecond)
	if err != nil {
		return nil, err
	}
	ctx.OCRSidecars[key] = sidecar
	return sidecar, nil
}

func (ctx *FlowContext) closeOCRSidecars() {
	if ctx == nil {
		return
	}
	for _, sidecar := range ctx.OCRSidecars {
		_ = sidecar.Close()
	}
}

func goddddocrSidecarKey(executable string, serverArgs []string, detection bool) string {
	parts := []string{executable, strconv.FormatBool(detection)}
	parts = append(parts, serverArgs...)
	return strings.Join(parts, "\x00")
}

func goddddocrStartupTimeoutMS(ctx *FlowContext, step FlowStep) (int, error) {
	if value, ok, err := flowStepResolvedParam(ctx, step, "startup_timeout"); err != nil {
		return 0, err
	} else if ok {
		timeoutMS, err := intParam(value)
		if err != nil {
			return 0, fmt.Errorf("%s startup_timeout %w", step.Action, err)
		}
		if timeoutMS < 1 {
			return 0, fmt.Errorf("%s startup_timeout must be at least 1", step.Action)
		}
		return timeoutMS, nil
	}
	if value, ok, err := flowStepResolvedParam(ctx, step, "timeout"); err != nil {
		return 0, err
	} else if ok {
		timeoutMS, err := intParam(value)
		if err != nil {
			return 0, fmt.Errorf("%s timeout %w", step.Action, err)
		}
		if timeoutMS > 0 {
			return timeoutMS, nil
		}
	}
	return defaultGoddddocrStartupTimeoutMS, nil
}

func startGoddddocrSidecar(ctx *FlowContext, executable string, serverArgs []string, detection bool, port int, timeout time.Duration) (*goddddocrSidecar, error) {
	args := append([]string(nil), serverArgs...)
	args = append(args, "-addr", fmt.Sprintf("127.0.0.1:%d", port))
	if detection && !stringListContains(args, "-det") {
		args = append(args, "-det")
	}

	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	cmd := exec.CommandContext(runCtx, executable, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	configureLocalCDPBrowserCommand(cmd)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start goddddocr sidecar %q: %w", executable, err)
	}
	sidecar := &goddddocrSidecar{
		BaseURL:  fmt.Sprintf("http://127.0.0.1:%d", port),
		cmd:      cmd,
		waitDone: make(chan error, 1),
		stdout:   stdout,
		stderr:   stderr,
	}
	go func() {
		sidecar.waitDone <- cmd.Wait()
	}()
	if err := waitForGoddddocrReady(sidecar.BaseURL, timeout, sidecar.waitDone); err != nil {
		_ = sidecar.Close()
		return nil, fmt.Errorf("%w; stderr=%s", err, trimDiagnosticText(stderr.String()))
	}
	return sidecar, nil
}

func (sidecar *goddddocrSidecar) Close() error {
	if sidecar == nil || sidecar.cmd == nil {
		return nil
	}
	var closeErr error
	sidecar.close.Do(func() {
		select {
		case <-sidecar.waitDone:
			return
		default:
		}
		if err := terminateLocalCDPBrowserCommand(sidecar.cmd); err != nil {
			closeErr = err
		}
		select {
		case <-sidecar.waitDone:
			return
		case <-time.After(3 * time.Second):
		}
		if err := killLocalCDPBrowserCommand(sidecar.cmd); err != nil && closeErr == nil {
			closeErr = err
		}
		select {
		case <-sidecar.waitDone:
		case <-time.After(2 * time.Second):
		}
	})
	return closeErr
}

func waitForGoddddocrReady(baseURL string, timeout time.Duration, waitDone <-chan error) error {
	if timeout <= 0 {
		timeout = time.Duration(defaultGoddddocrStartupTimeoutMS) * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	readyURL := normalizeGoddddocrReadyEndpoint(baseURL)
	client := &http.Client{Timeout: 750 * time.Millisecond}
	var lastErr error
	for {
		select {
		case err := <-waitDone:
			if err == nil {
				return fmt.Errorf("goddddocr sidecar exited before /ready")
			}
			return fmt.Errorf("goddddocr sidecar exited before /ready: %w", err)
		default:
		}

		resp, err := client.Get(readyURL)
		if err == nil {
			body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
			_ = resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 && goddddocrReadyBody(body) {
				return nil
			} else {
				lastErr = fmt.Errorf("ready status=%d body=%s", resp.StatusCode, trimDiagnosticText(string(body)))
			}
		} else {
			lastErr = err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait for goddddocr sidecar /ready timed out after %s: %v", timeout, lastErr)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func goddddocrReadyBody(body []byte) bool {
	decoded := map[string]any{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return false
	}
	if ready, _ := decoded["ready"].(bool); ready {
		return true
	}
	status := strings.TrimSpace(fmt.Sprint(decoded["status"]))
	return strings.EqualFold(status, "ready")
}

func goddddocrSidecarExecutable(ctx *FlowContext, step FlowStep) (string, error) {
	explicit, err := flowStepOptionalStringParam(ctx, step, "executable")
	if err != nil {
		return "", err
	}
	return resolveGoddddocrExecutable(explicit,
		[]string{"GODDDDOCR_SERVER_BIN", "GODDDDOCR_BIN"},
		goddddocrSidecarExecutableCandidates(),
		"goddddocr-server",
	)
}

func goddddocrCLIExecutable(ctx *FlowContext, step FlowStep) (string, error) {
	explicit, err := flowStepOptionalStringParam(ctx, step, "executable")
	if err != nil {
		return "", err
	}
	return resolveGoddddocrExecutable(explicit,
		[]string{"GODDDDOCR_CLI_BIN", "GODDDDOCR_OCR_BIN", "OCRDOCTOR_BIN", "GODDDDOCR_BIN"},
		goddddocrCLIExecutableCandidates(),
		"goddddocr/ocrdoctor",
	)
}

func resolveGoddddocrExecutable(explicit string, envNames []string, candidates []string, label string) (string, error) {
	searched := []string{}
	if resolved, err := resolveGoddddocrExecutableCandidate(explicit); err == nil && resolved != "" {
		return resolved, nil
	} else if strings.TrimSpace(explicit) != "" {
		return "", fmt.Errorf("resolve %s executable %q: %w", label, explicit, err)
	}
	for _, envName := range envNames {
		value := strings.TrimSpace(os.Getenv(envName))
		if value == "" {
			continue
		}
		searched = append(searched, envName+"="+value)
		if resolved, err := resolveGoddddocrExecutableCandidate(value); err == nil && resolved != "" {
			return resolved, nil
		}
	}
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate) == "" {
			continue
		}
		searched = append(searched, candidate)
		if resolved, err := resolveGoddddocrExecutableCandidate(candidate); err == nil && resolved != "" {
			return resolved, nil
		}
	}
	return "", fmt.Errorf("could not find %s executable; set executable, one of %s, or put it on PATH. searched: %s", label, strings.Join(envNames, "/"), strings.Join(searched, ", "))
}

func resolveGoddddocrExecutableCandidate(candidate string) (string, error) {
	candidate = strings.TrimSpace(expandHomePath(candidate))
	if candidate == "" {
		return "", nil
	}
	if !strings.Contains(candidate, string(os.PathSeparator)) && !strings.Contains(candidate, "/") && !strings.Contains(candidate, "\\") {
		if resolved, err := exec.LookPath(candidate); err == nil {
			return resolved, nil
		}
		return "", fmt.Errorf("not found on PATH")
	}
	abs, err := filepath.Abs(candidate)
	if err != nil {
		return "", err
	}
	if isExecutableFile(abs) {
		return abs, nil
	}
	return "", fmt.Errorf("not an executable file")
}

func goddddocrSidecarExecutableCandidates() []string {
	names := []string{"goddddocr-server"}
	if runtime.GOOS == "windows" {
		names = append(names, "goddddocr-server.exe")
	}
	candidates := append([]string(nil), names...)
	for _, name := range names {
		candidates = append(candidates, filepath.Join(".", name), filepath.Join(".", "bin", name))
	}
	if current, err := os.Executable(); err == nil {
		for _, name := range names {
			candidates = append(candidates, filepath.Join(filepath.Dir(current), name))
		}
	}
	return candidates
}

func goddddocrCLIExecutableCandidates() []string {
	names := []string{"goddddocr", "ocrdoctor"}
	if runtime.GOOS == "windows" {
		names = append(names, "goddddocr.exe", "ocrdoctor.exe")
	}
	candidates := append([]string(nil), names...)
	for _, name := range names {
		candidates = append(candidates, filepath.Join(".", name), filepath.Join(".", "bin", name))
	}
	if current, err := os.Executable(); err == nil {
		for _, name := range names {
			candidates = append(candidates, filepath.Join(filepath.Dir(current), name))
		}
	}
	return candidates
}

func runGoddddocrCLIRequest(ctx *FlowContext, step FlowStep, filePath string, fields map[string]any) (map[string]any, error) {
	if err := requireGoddddocrProcess(ctx, step.Action); err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Security != nil {
		var err error
		filePath, err = resolveRuntimeFilePath(filePath, flowFileInputPath, *ctx.Security)
		if err != nil {
			return nil, fmt.Errorf("ocr_request parameter %q %w", "file_path", err)
		}
	}
	executable, err := goddddocrCLIExecutable(ctx, step)
	if err != nil {
		return nil, err
	}
	args, err := goddddocrCLIArgs(ctx, step, executable, filePath, fields)
	if err != nil {
		return nil, err
	}
	timeoutMS, err := goddddocrCLITimeoutMS(ctx, step)
	if err != nil {
		return nil, err
	}

	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	runCtx, cancel := context.WithTimeout(runCtx, time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(runCtx, executable, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = -1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		if runCtx.Err() != nil {
			return nil, fmt.Errorf("ocr_request cli timed out after %dms: %w", timeoutMS, runCtx.Err())
		}
	}

	body := map[string]any{}
	if parseErr := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &body); parseErr != nil {
		if err != nil {
			return nil, fmt.Errorf("ocr_request cli failed with exit status %d: %v; stderr=%s; stdout=%s", exitCode, err, trimDiagnosticText(stderr.String()), trimDiagnosticText(stdout.String()))
		}
		return nil, fmt.Errorf("ocr_request cli expected stdout JSON object: %w; stdout=%s", parseErr, trimDiagnosticText(stdout.String()))
	}
	if _, ok := body["processing_time_ms"]; !ok {
		if elapsed, ok := body["elapsed_ms"]; ok {
			body["processing_time_ms"] = elapsed
		}
	}
	ok := exitCode == 0
	if okValue, hasOK := body["ok"].(bool); hasOK {
		ok = okValue
	}
	if !ok {
		reason := firstNonEmpty(strings.TrimSpace(fmt.Sprint(body["error"])), trimDiagnosticText(stderr.String()), trimDiagnosticText(stdout.String()))
		return nil, fmt.Errorf("ocr_request cli failed with exit status %d: %s", exitCode, reason)
	}

	savePath, err := saveGoddddocrCLIResponse(ctx, step, stdout.Bytes(), body)
	if err != nil {
		return nil, err
	}
	response := map[string]any{
		"ok":            true,
		"status":        exitCode,
		"body":          body,
		"service_mode":  goddddocrModeCLI,
		"executable":    executable,
		"args":          args,
		"stderr":        trimDiagnosticText(stderr.String()),
		"stdout_length": stdout.Len(),
	}
	if savePath != "" {
		response["save_path"] = savePath
	}
	result, err := buildOCRRequestResult(response)
	if err != nil {
		return nil, err
	}
	annotateGoddddocrResult(result, goddddocrModeCLI, "", nil)
	return result, nil
}

func goddddocrCLIArgs(ctx *FlowContext, step FlowStep, executable string, filePath string, fields map[string]any) ([]string, error) {
	prefixArgs, err := flowStepOptionalStringListParam(ctx, step, "cli_args")
	if err != nil {
		return nil, err
	}
	if len(prefixArgs) > 0 {
		rewritten, hadPlaceholder := substituteGoddddocrCLIFilePlaceholders(prefixArgs, filePath)
		if hadPlaceholder {
			return rewritten, nil
		}
		return append(rewritten, defaultGoddddocrCLIArgs(executable, filePath, fields)...), nil
	}
	return defaultGoddddocrCLIArgs(executable, filePath, fields), nil
}

func defaultGoddddocrCLIArgs(executable string, filePath string, fields map[string]any) []string {
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(executable)), ".exe")
	charsetRange := ""
	if raw, ok := fields["charset_range"]; ok {
		charsetRange = strings.TrimSpace(fmt.Sprint(raw))
	}
	if base == "goddddocr" {
		args := []string{"ocr", "--file", filePath, "--json"}
		if confidence, _ := fields["confidence"].(bool); confidence {
			args = append(args, "--confidence")
		}
		if charsetRange != "" {
			args = append(args, "--charset-range", charsetRange)
		}
		return args
	}
	args := []string{"-image", filePath, "-json"}
	if charsetRange != "" {
		args = append(args, "-charset-range", charsetRange)
	}
	return args
}

func substituteGoddddocrCLIFilePlaceholders(args []string, filePath string) ([]string, bool) {
	rewritten := make([]string, 0, len(args))
	hadPlaceholder := false
	for _, arg := range args {
		next := arg
		for _, placeholder := range []string{"{file}", "{file_path}", "{{file_path}}"} {
			if strings.Contains(next, placeholder) {
				hadPlaceholder = true
				next = strings.ReplaceAll(next, placeholder, filePath)
			}
		}
		rewritten = append(rewritten, next)
	}
	return rewritten, hadPlaceholder
}

func goddddocrCLITimeoutMS(ctx *FlowContext, step FlowStep) (int, error) {
	if value, ok, err := flowStepResolvedParam(ctx, step, "timeout"); err != nil {
		return 0, err
	} else if ok {
		timeoutMS, err := intParam(value)
		if err != nil {
			return 0, fmt.Errorf("ocr_request timeout %w", err)
		}
		if timeoutMS < 1 {
			return 0, fmt.Errorf("ocr_request timeout must be at least 1")
		}
		return timeoutMS, nil
	}
	return defaultGoddddocrCLITimeoutMS, nil
}

func saveGoddddocrCLIResponse(ctx *FlowContext, step FlowStep, stdout []byte, body map[string]any) (string, error) {
	savePath, err := flowStepOptionalStringParam(ctx, step, "save_path")
	if err != nil {
		return "", err
	}
	savePath = strings.TrimSpace(savePath)
	if savePath == "" {
		return "", nil
	}
	if ctx != nil && ctx.Security != nil {
		savePath, err = resolveRuntimeFilePath(savePath, flowFileOutputPath, *ctx.Security)
		if err != nil {
			return "", fmt.Errorf("ocr_request parameter %q %w", "save_path", err)
		}
	}
	content := bytes.TrimSpace(stdout)
	if len(content) == 0 {
		content, err = json.MarshalIndent(body, "", "  ")
		if err != nil {
			return "", fmt.Errorf("encode ocr_request cli response: %w", err)
		}
	}
	if err := os.WriteFile(savePath, append(content, '\n'), 0600); err != nil {
		return "", fmt.Errorf("write ocr_request save_path %q: %w", savePath, err)
	}
	return savePath, nil
}

func goddddocrOptionalBoolParam(ctx *FlowContext, step FlowStep, name string) (bool, bool, error) {
	value, ok, err := flowStepResolvedParam(ctx, step, name)
	if err != nil || !ok {
		return false, ok, err
	}
	parsed, err := ocrBoolParam(value)
	if err != nil {
		return false, true, fmt.Errorf("%s %s %w", step.Action, name, err)
	}
	return parsed, true, nil
}

func annotateGoddddocrResult(result map[string]any, mode string, baseURL string, sidecar *goddddocrSidecar) {
	if result == nil {
		return
	}
	result["service_mode"] = mode
	if baseURL != "" {
		result["service_url"] = baseURL
	}
	if sidecar != nil && sidecar.cmd != nil && sidecar.cmd.Process != nil {
		result["sidecar"] = map[string]any{
			"managed": true,
			"url":     sidecar.BaseURL,
			"pid":     sidecar.cmd.Process.Pid,
		}
	}
}

func stringListContains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func trimDiagnosticText(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 1000 {
		return value
	}
	return value[:1000] + "...(truncated)"
}
