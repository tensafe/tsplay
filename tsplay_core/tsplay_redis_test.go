package tsplay_core

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

type redisTestServer struct {
	listener net.Listener
	mu       sync.Mutex
	store    map[int]map[string]string
}

func newRedisTestServer(t *testing.T) *redisTestServer {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen redis test server: %v", err)
	}
	server := &redisTestServer{
		listener: listener,
		store:    map[int]map[string]string{},
	}
	go server.serve()
	return server
}

func (s *redisTestServer) Addr() string {
	if s == nil || s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func (s *redisTestServer) Close() {
	if s == nil || s.listener == nil {
		return
	}
	_ = s.listener.Close()
}

func (s *redisTestServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}

func (s *redisTestServer) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	db := 0

	for {
		value, err := readRedisValue(reader)
		if err != nil {
			return
		}
		args, err := redisArrayToStrings(value)
		if err != nil || len(args) == 0 {
			_ = writeRedisError(writer, "ERR invalid command")
			return
		}

		command := strings.ToUpper(args[0])
		switch command {
		case "AUTH", "SELECT":
			if command == "SELECT" && len(args) > 1 {
				nextDB, err := strconv.Atoi(args[1])
				if err == nil {
					db = nextDB
				}
			}
			_ = writeRedisSimple(writer, "OK")
		case "GET":
			if len(args) != 2 {
				_ = writeRedisError(writer, "ERR wrong number of arguments for GET")
				continue
			}
			value, ok := s.get(db, args[1])
			if !ok {
				_ = writeRedisNilBulk(writer)
				continue
			}
			_ = writeRedisBulk(writer, value)
		case "SET":
			if len(args) < 3 {
				_ = writeRedisError(writer, "ERR wrong number of arguments for SET")
				continue
			}
			s.set(db, args[1], args[2])
			_ = writeRedisSimple(writer, "OK")
		case "DEL":
			if len(args) != 2 {
				_ = writeRedisError(writer, "ERR wrong number of arguments for DEL")
				continue
			}
			deleted := s.del(db, args[1])
			_ = writeRedisInteger(writer, int64(deleted))
		case "INCRBY":
			if len(args) != 3 {
				_ = writeRedisError(writer, "ERR wrong number of arguments for INCRBY")
				continue
			}
			delta, err := strconv.Atoi(args[2])
			if err != nil {
				_ = writeRedisError(writer, "ERR delta must be integer")
				continue
			}
			value, err := s.incr(db, args[1], delta)
			if err != nil {
				_ = writeRedisError(writer, "ERR "+err.Error())
				continue
			}
			_ = writeRedisInteger(writer, int64(value))
		default:
			_ = writeRedisError(writer, "ERR unsupported command "+command)
			return
		}
	}
}

func (s *redisTestServer) get(db int, key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := s.store[db]
	if values == nil {
		return "", false
	}
	value, ok := values[key]
	return value, ok
}

func (s *redisTestServer) set(db int, key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := s.store[db]
	if values == nil {
		values = map[string]string{}
		s.store[db] = values
	}
	values[key] = value
}

func (s *redisTestServer) del(db int, key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := s.store[db]
	if values == nil {
		return 0
	}
	if _, ok := values[key]; !ok {
		return 0
	}
	delete(values, key)
	return 1
}

func (s *redisTestServer) incr(db int, key string, delta int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := s.store[db]
	if values == nil {
		values = map[string]string{}
		s.store[db] = values
	}
	current := 0
	if raw, ok := values[key]; ok && raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("value at %q is not an integer", key)
		}
		current = parsed
	}
	current += delta
	values[key] = strconv.Itoa(current)
	return current, nil
}

func redisArrayToStrings(value redisRESPValue) ([]string, error) {
	if value.kind != '*' {
		return nil, fmt.Errorf("expected redis array, got %q", value.kind)
	}
	items := make([]string, 0, len(value.values))
	for _, item := range value.values {
		if item.nil {
			items = append(items, "")
			continue
		}
		items = append(items, item.text)
	}
	return items, nil
}

func writeRedisSimple(writer *bufio.Writer, text string) error {
	if _, err := fmt.Fprintf(writer, "+%s\r\n", text); err != nil {
		return err
	}
	return writer.Flush()
}

func writeRedisError(writer *bufio.Writer, text string) error {
	if _, err := fmt.Fprintf(writer, "-%s\r\n", text); err != nil {
		return err
	}
	return writer.Flush()
}

func writeRedisInteger(writer *bufio.Writer, value int64) error {
	if _, err := fmt.Fprintf(writer, ":%d\r\n", value); err != nil {
		return err
	}
	return writer.Flush()
}

func writeRedisBulk(writer *bufio.Writer, text string) error {
	if _, err := fmt.Fprintf(writer, "$%d\r\n%s\r\n", len(text), text); err != nil {
		return err
	}
	return writer.Flush()
}

func writeRedisNilBulk(writer *bufio.Writer) error {
	if _, err := writer.WriteString("$-1\r\n"); err != nil {
		return err
	}
	return writer.Flush()
}

func TestValidateFlowSecurityRejectsRedisByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "redis_policy",
		Steps: []FlowStep{
			{Action: "redis_get", Key: "sessions:admin"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected redis security policy error")
	}
	if !strings.Contains(err.Error(), "allow_redis") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowRedisActions(t *testing.T) {
	server := newRedisTestServer(t)
	defer server.Close()

	t.Setenv("TSPLAY_REDIS_ADDR", server.Addr())

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "redis_round_trip",
		Steps: []FlowStep{
			{Action: "redis_set", Key: "sessions:admin_cookie", Value: "SESSION=abc", TTLSeconds: 3600},
			{Action: "redis_get", Key: "sessions:admin_cookie", SaveAs: "cookie_header"},
			{
				Action: "redis_set",
				Key:    "sessions:admin_payload",
				With: map[string]any{
					"value": map[string]any{
						"cookie": "SESSION=abc",
						"user":   "admin",
					},
				},
			},
			{Action: "redis_get", Key: "sessions:admin_payload", SaveAs: "session_payload"},
			{Action: "json_extract", From: "{{session_payload}}", Path: "$.cookie", SaveAs: "cookie_value"},
			{Action: "redis_incr", Key: "orders:counter", Delta: 2, SaveAs: "counter_two"},
			{Action: "redis_incr", Key: "orders:counter", SaveAs: "counter_three"},
			{Action: "redis_del", Key: "sessions:admin_cookie", SaveAs: "deleted_count"},
			{Action: "redis_get", Key: "sessions:admin_cookie", SaveAs: "deleted_value"},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowRedis: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["cookie_header"]; got != "SESSION=abc" {
		t.Fatalf("cookie_header = %#v", got)
	}
	if got := result.Vars["cookie_value"]; got != "SESSION=abc" {
		t.Fatalf("cookie_value = %#v", got)
	}
	if got := result.Vars["counter_two"]; got != float64(2) {
		t.Fatalf("counter_two = %#v", got)
	}
	if got := result.Vars["counter_three"]; got != float64(3) {
		t.Fatalf("counter_three = %#v", got)
	}
	if got := result.Vars["deleted_count"]; got != float64(1) {
		t.Fatalf("deleted_count = %#v", got)
	}
	if got, ok := result.Vars["deleted_value"]; ok && got != nil {
		t.Fatalf("deleted_value = %#v", got)
	}
}

func TestRunFlowRedisNamedConnection(t *testing.T) {
	server := newRedisTestServer(t)
	defer server.Close()

	t.Setenv("TSPLAY_REDIS_SESSIONS_ADDR", server.Addr())

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "redis_named_connection",
		Steps: []FlowStep{
			{Action: "redis_set", Key: "sessions:admin_cookie", Value: "SESSION=xyz", Connection: "sessions"},
			{Action: "redis_get", Key: "sessions:admin_cookie", Connection: "sessions", SaveAs: "cookie_header"},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowRedis: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["cookie_header"]; got != "SESSION=xyz" {
		t.Fatalf("cookie_header = %#v", got)
	}
}

func TestRunFlowLuaRedisHelpersHonorAllowRedis(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_redis_policy",
		Steps: []FlowStep{
			{
				Action: "lua",
				Code:   `return redis_get("sessions:admin")`,
			},
		},
	}

	_, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowLua: true},
	})
	if err == nil {
		t.Fatalf("expected allow_redis runtime error")
	}
	if !strings.Contains(err.Error(), "allow_redis") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowLuaRedisHelpers(t *testing.T) {
	server := newRedisTestServer(t)
	defer server.Close()

	t.Setenv("TSPLAY_REDIS_ADDR", server.Addr())

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_redis_round_trip",
		Steps: []FlowStep{
			{
				Action: "lua",
				SaveAs: "redis_result",
				Code: `local stored = redis_set("sessions:admin_cookie", "SESSION=lua", 3600)
local value = redis_get("sessions:admin_cookie")
local counter = redis_incr("orders:counter", 2)
local deleted = redis_del("sessions:admin_cookie")
local missing = redis_get("sessions:admin_cookie")
return {
  stored = stored,
  value = value,
  counter = counter,
  deleted = deleted,
  missing = missing,
}`,
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowLua:   true,
			AllowRedis: true,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	redisResult, ok := result.Vars["redis_result"].(map[string]any)
	if !ok {
		t.Fatalf("redis_result = %#v", result.Vars["redis_result"])
	}
	if got := redisResult["value"]; got != "SESSION=lua" {
		t.Fatalf("value = %#v", got)
	}
	if got := redisResult["counter"]; got != float64(2) {
		t.Fatalf("counter = %#v", got)
	}
	if got := redisResult["deleted"]; got != float64(1) {
		t.Fatalf("deleted = %#v", got)
	}
	if got, exists := redisResult["missing"]; exists && got != nil {
		t.Fatalf("missing = %#v", got)
	}
	stored, ok := redisResult["stored"].(map[string]any)
	if !ok || stored["key"] != "sessions:admin_cookie" {
		t.Fatalf("stored = %#v", redisResult["stored"])
	}
}
