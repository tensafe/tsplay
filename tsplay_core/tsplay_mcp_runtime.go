package tsplay_core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultTSPlayBrowserRunTimeoutMS    = 120000
	defaultTSPlayFlowRunTimeoutMS       = 180000
	defaultTSPlayBrowserRunQueueTimeout = 15000
	defaultTSPlayBrowserRunGlobalLimit  = 2
	defaultTSPlayBrowserRunSessionLimit = 1
	defaultTSPlayBrowserRunTimeoutMaxMS = 600000
	defaultTSPlayBrowserRunFolderName   = "mcp_runs"
)

type TSPlayMCPCaller struct {
	SessionID     string `json:"session_id,omitempty"`
	ClientName    string `json:"client_name,omitempty"`
	ClientVersion string `json:"client_version,omitempty"`
}

type TSPlayBrowserRun struct {
	ID             string             `json:"id"`
	Tool           string             `json:"tool"`
	Status         string             `json:"status"`
	QueuedAt       string             `json:"queued_at,omitempty"`
	StartedAt      string             `json:"started_at,omitempty"`
	FinishedAt     string             `json:"finished_at,omitempty"`
	QueueWaitMS    int64              `json:"queue_wait_ms,omitempty"`
	DurationMS     int64              `json:"duration_ms,omitempty"`
	TimeoutMS      int                `json:"timeout_ms,omitempty"`
	QueueTimeoutMS int                `json:"queue_timeout_ms,omitempty"`
	ArtifactRoot   string             `json:"artifact_root,omitempty"`
	RunRoot        string             `json:"run_root,omitempty"`
	AuditPath      string             `json:"audit_path,omitempty"`
	Caller         TSPlayMCPCaller    `json:"caller,omitempty"`
	Grants         FlowSecurityPolicy `json:"grants,omitempty"`
	Details        map[string]any     `json:"details,omitempty"`
	Error          string             `json:"error,omitempty"`
}

type tsplayBrowserRunAudit struct {
	Run       TSPlayBrowserRun `json:"run"`
	Arguments any              `json:"arguments,omitempty"`
}

type tsplayBrowserRunHandle struct {
	run       TSPlayBrowserRun
	arguments any
	limiter   *tsplayBrowserRunLimiter
	session   string
	cancel    context.CancelFunc
	release   sync.Once
}

type tsplayBrowserRunLimiter struct {
	global       chan struct{}
	sessionLimit int
	sessions     sync.Map
}

var tsplayBrowserRunLimiterCache sync.Map

func getTSPlayBrowserRunLimiter(globalLimit int, sessionLimit int) *tsplayBrowserRunLimiter {
	if globalLimit <= 0 {
		globalLimit = defaultTSPlayBrowserRunGlobalLimit
	}
	if sessionLimit <= 0 {
		sessionLimit = defaultTSPlayBrowserRunSessionLimit
	}
	key := fmt.Sprintf("%d:%d", globalLimit, sessionLimit)
	if limiter, ok := tsplayBrowserRunLimiterCache.Load(key); ok {
		return limiter.(*tsplayBrowserRunLimiter)
	}
	limiter := &tsplayBrowserRunLimiter{
		global:       make(chan struct{}, globalLimit),
		sessionLimit: sessionLimit,
	}
	actual, _ := tsplayBrowserRunLimiterCache.LoadOrStore(key, limiter)
	return actual.(*tsplayBrowserRunLimiter)
}

func (limiter *tsplayBrowserRunLimiter) acquire(ctx context.Context, session string) error {
	sessionChan := limiter.sessionSemaphore(session)
	if err := acquireTSPlayBrowserRunToken(ctx, sessionChan); err != nil {
		return err
	}
	if err := acquireTSPlayBrowserRunToken(ctx, limiter.global); err != nil {
		releaseTSPlayBrowserRunToken(sessionChan)
		return err
	}
	return nil
}

func (limiter *tsplayBrowserRunLimiter) releaseForSession(session string) {
	releaseTSPlayBrowserRunToken(limiter.global)
	releaseTSPlayBrowserRunToken(limiter.sessionSemaphore(session))
}

func (limiter *tsplayBrowserRunLimiter) sessionSemaphore(session string) chan struct{} {
	session = strings.TrimSpace(session)
	if session == "" {
		session = "anonymous"
	}
	if semaphore, ok := limiter.sessions.Load(session); ok {
		return semaphore.(chan struct{})
	}
	created := make(chan struct{}, limiter.sessionLimit)
	actual, _ := limiter.sessions.LoadOrStore(session, created)
	return actual.(chan struct{})
}

func acquireTSPlayBrowserRunToken(ctx context.Context, semaphore chan struct{}) error {
	select {
	case semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func releaseTSPlayBrowserRunToken(semaphore chan struct{}) {
	select {
	case <-semaphore:
	default:
	}
}

func beginTSPlayBrowserRun(
	ctx context.Context,
	request mcp.CallToolRequest,
	toolName string,
	options TSPlayMCPServerOptions,
	grants *FlowSecurityPolicy,
) (*tsplayBrowserRunHandle, context.Context, error) {
	caller := tsplayMCPCallerFromContext(ctx)
	baseRoot, rootErr := prepareTSPlayBrowserRunArtifactRoot(options.ArtifactRoot)
	timeoutMS := tsplayBrowserRunTimeoutMS(toolName, request.GetInt("run_timeout", 0), options)
	queueTimeoutMS := options.QueueTimeoutMS
	if queueTimeoutMS <= 0 {
		queueTimeoutMS = defaultTSPlayBrowserRunQueueTimeout
	}
	if timeoutMS > 0 && queueTimeoutMS > timeoutMS {
		queueTimeoutMS = timeoutMS
	}

	run := TSPlayBrowserRun{
		ID:             newTSPlayBrowserRunID(toolName),
		Tool:           toolName,
		Status:         "queued",
		QueuedAt:       time.Now().Format(time.RFC3339Nano),
		TimeoutMS:      timeoutMS,
		QueueTimeoutMS: queueTimeoutMS,
		Caller:         caller,
		ArtifactRoot:   strings.TrimSpace(options.ArtifactRoot),
	}
	if grants != nil {
		run.Grants = *grants
	}

	if rootErr == nil {
		run.ArtifactRoot = baseRoot
		run.RunRoot = filepath.Join(baseRoot, defaultTSPlayBrowserRunFolderName, sanitizeArtifactSegment(tsplayMCPCallerSessionKey(caller)), run.ID)
		run.AuditPath = filepath.Join(run.RunRoot, "run.json")
	}

	handle := &tsplayBrowserRunHandle{
		run:       run,
		arguments: compactTraceValue(request.GetArguments(), 0),
	}
	if err := handle.writeAudit(); err != nil && rootErr == nil {
		rootErr = err
	}
	if rootErr != nil {
		handle.finish(rootErr, nil)
		return handle, nil, rootErr
	}

	limiter := getTSPlayBrowserRunLimiter(options.MaxConcurrentBrowserRuns, options.MaxConcurrentBrowserRunsPerSession)
	handle.limiter = limiter
	handle.session = tsplayMCPCallerSessionKey(caller)

	acquireCtx, cancelAcquire := context.WithTimeout(ctx, time.Duration(queueTimeoutMS)*time.Millisecond)
	defer cancelAcquire()
	if err := limiter.acquire(acquireCtx, handle.session); err != nil {
		if err == context.DeadlineExceeded {
			err = fmt.Errorf("%s concurrency queue timed out after %dms", toolName, queueTimeoutMS)
		}
		handle.finish(err, nil)
		return handle, nil, err
	}

	handle.run.Status = "running"
	handle.run.StartedAt = time.Now().Format(time.RFC3339Nano)
	if queuedAt, err := time.Parse(time.RFC3339Nano, handle.run.QueuedAt); err == nil {
		handle.run.QueueWaitMS = time.Since(queuedAt).Milliseconds()
	}
	remainingTimeoutMS := timeoutMS - int(handle.run.QueueWaitMS)
	if remainingTimeoutMS < 1 {
		remainingTimeoutMS = 1
	}
	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(remainingTimeoutMS)*time.Millisecond)
	handle.cancel = cancelRun
	if err := handle.writeAudit(); err != nil {
		handle.finish(err, nil)
		return handle, nil, err
	}
	return handle, runCtx, nil
}

func (handle *tsplayBrowserRunHandle) snapshot() TSPlayBrowserRun {
	run := handle.run
	if len(handle.run.Details) > 0 {
		run.Details = map[string]any{}
		for key, value := range handle.run.Details {
			run.Details[key] = value
		}
	}
	return run
}

func (handle *tsplayBrowserRunHandle) finish(err error, details map[string]any) TSPlayBrowserRun {
	if handle == nil {
		return TSPlayBrowserRun{}
	}
	handle.release.Do(func() {
		if handle.cancel != nil {
			handle.cancel()
		}
		if handle.limiter != nil {
			handle.limiter.releaseForSession(handle.session)
		}
		handle.run.FinishedAt = time.Now().Format(time.RFC3339Nano)
		if startedAt, parseErr := time.Parse(time.RFC3339Nano, handle.run.StartedAt); parseErr == nil {
			handle.run.DurationMS = time.Since(startedAt).Milliseconds()
		}
		if err != nil {
			handle.run.Status = tsplayBrowserRunStatusForError(err)
			handle.run.Error = err.Error()
		} else if handle.run.Status == "queued" || handle.run.Status == "running" {
			handle.run.Status = "ok"
		}
		if len(details) > 0 {
			handle.run.Details = compactTraceValue(details, 0).(map[string]any)
		}
		_ = handle.writeAudit()
	})
	return handle.snapshot()
}

func (handle *tsplayBrowserRunHandle) writeAudit() error {
	if handle == nil || strings.TrimSpace(handle.run.AuditPath) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(handle.run.AuditPath), 0755); err != nil {
		return fmt.Errorf("create browser run audit directory: %w", err)
	}
	content, err := json.MarshalIndent(tsplayBrowserRunAudit{
		Run:       handle.snapshot(),
		Arguments: handle.arguments,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal browser run audit: %w", err)
	}
	if err := os.WriteFile(handle.run.AuditPath, content, 0644); err != nil {
		return fmt.Errorf("write browser run audit %q: %w", handle.run.AuditPath, err)
	}
	return nil
}

func prepareTSPlayBrowserRunArtifactRoot(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		root = DefaultMCPArtifactRoot
	}
	return prepareRuntimeFileRoot(root)
}

func newTSPlayBrowserRunID(toolName string) string {
	toolName = strings.TrimSpace(strings.TrimPrefix(toolName, "tsplay."))
	if toolName == "" {
		toolName = "browser_run"
	}
	return sanitizeArtifactSegment(toolName) + "-" + time.Now().Format("20060102-150405.000000000")
}

func tsplayBrowserRunTimeoutMS(toolName string, requested int, options TSPlayMCPServerOptions) int {
	defaultTimeout := options.DefaultRunTimeoutMS
	switch strings.TrimSpace(toolName) {
	case "tsplay.run_flow":
		if defaultTimeout <= 0 || defaultTimeout == defaultTSPlayBrowserRunTimeoutMS {
			defaultTimeout = defaultTSPlayFlowRunTimeoutMS
		}
	default:
		if defaultTimeout <= 0 {
			defaultTimeout = defaultTSPlayBrowserRunTimeoutMS
		}
	}
	if requested > 0 {
		defaultTimeout = requested
	}
	maxTimeout := options.MaxRunTimeoutMS
	if maxTimeout <= 0 {
		maxTimeout = defaultTSPlayBrowserRunTimeoutMaxMS
	}
	if defaultTimeout > maxTimeout {
		defaultTimeout = maxTimeout
	}
	if defaultTimeout <= 0 {
		return defaultTSPlayBrowserRunTimeoutMS
	}
	return defaultTimeout
}

func tsplayMCPCallerFromContext(ctx context.Context) TSPlayMCPCaller {
	session := server.ClientSessionFromContext(ctx)
	if session == nil {
		return TSPlayMCPCaller{}
	}
	caller := TSPlayMCPCaller{
		SessionID: strings.TrimSpace(session.SessionID()),
	}
	if withInfo, ok := session.(server.SessionWithClientInfo); ok {
		info := withInfo.GetClientInfo()
		caller.ClientName = strings.TrimSpace(info.Name)
		caller.ClientVersion = strings.TrimSpace(info.Version)
	}
	return caller
}

func flowSavedSessionAccessFromContext(ctx context.Context) FlowSavedSessionAccessInfo {
	caller := tsplayMCPCallerFromContext(ctx)
	return FlowSavedSessionAccessInfo{
		SessionID:     caller.SessionID,
		ClientName:    caller.ClientName,
		ClientVersion: caller.ClientVersion,
	}
}

func tsplayMCPCallerSessionKey(caller TSPlayMCPCaller) string {
	if strings.TrimSpace(caller.SessionID) != "" {
		return caller.SessionID
	}
	return "anonymous"
}

func tsplayBrowserRunStatusForError(err error) string {
	if err == nil {
		return "ok"
	}
	switch {
	case strings.Contains(strings.ToLower(err.Error()), "timed out"):
		return "timed_out"
	case err == context.Canceled:
		return "canceled"
	case err == context.DeadlineExceeded:
		return "timed_out"
	default:
		return "error"
	}
}

func watchContextCancel(ctx context.Context, cancel func()) func() {
	if ctx == nil || cancel == nil {
		return func() {}
	}
	stopped := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-stopped:
		}
	}()
	return func() {
		close(stopped)
	}
}
