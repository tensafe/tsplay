package tsplay_core

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const redisDefaultConnection = "default"

type redisConnectionConfig struct {
	Name      string
	Addr      string
	Username  string
	Password  string
	DB        int
	UseTLS    bool
	TimeoutMS int
}

type redisRESPValue struct {
	kind    byte
	text    string
	integer int64
	values  []redisRESPValue
	nil     bool
}

func redis_get(L *lua.LState) int {
	key := L.CheckString(1)
	connection := ""
	if L.GetTop() >= 2 {
		connection = L.OptString(2, "")
	}
	if err := luaRedisExecutionAllowed(L, "redis_get"); err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	value, err := redisGet(key, connection)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, value))
	return 1
}

func redis_set(L *lua.LState) int {
	key := L.CheckString(1)
	value := luaValueToGo(L.CheckAny(2))
	ttlSeconds := 0
	connection := ""
	if L.GetTop() >= 3 {
		switch third := L.Get(3).(type) {
		case lua.LNumber:
			ttlSeconds = int(third)
			if L.GetTop() >= 4 {
				connection = L.OptString(4, "")
			}
		case lua.LString:
			connection = string(third)
		}
	}
	if err := luaRedisExecutionAllowed(L, "redis_set"); err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	result, err := redisSet(key, value, ttlSeconds, connection)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func redis_del(L *lua.LState) int {
	key := L.CheckString(1)
	connection := ""
	if L.GetTop() >= 2 {
		connection = L.OptString(2, "")
	}
	if err := luaRedisExecutionAllowed(L, "redis_del"); err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	deleted, err := redisDel(key, connection)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, deleted))
	return 1
}

func redis_incr(L *lua.LState) int {
	key := L.CheckString(1)
	delta := 1
	connection := ""
	if L.GetTop() >= 2 {
		switch second := L.Get(2).(type) {
		case lua.LNumber:
			delta = int(second)
			if L.GetTop() >= 3 {
				connection = L.OptString(3, "")
			}
		case lua.LString:
			connection = string(second)
		}
	}
	if err := luaRedisExecutionAllowed(L, "redis_incr"); err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	value, err := redisIncr(key, delta, connection)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, value))
	return 1
}

func luaRedisExecutionAllowed(L *lua.LState, action string) error {
	flowCtx := flowContextFromState(L)
	if flowCtx != nil && flowCtx.Security != nil && !flowCtx.Security.AllowRedis {
		return fmt.Errorf("%s is disabled by security policy; set allow_redis=true only for trusted flows", action)
	}
	return nil
}

func redisGet(key string, connection string) (any, error) {
	config, err := resolveRedisConnectionConfig(connection)
	if err != nil {
		return nil, err
	}
	reply, err := executeRedisCommand(config, "GET", key)
	if err != nil {
		return nil, err
	}
	if reply.nil {
		return nil, nil
	}
	return reply.text, nil
}

func redisSet(key string, value any, ttlSeconds int, connection string) (map[string]any, error) {
	config, err := resolveRedisConnectionConfig(connection)
	if err != nil {
		return nil, err
	}
	encoded, err := encodeRedisValue(value)
	if err != nil {
		return nil, err
	}

	args := []string{"SET", key, encoded}
	if ttlSeconds > 0 {
		args = append(args, "EX", strconv.Itoa(ttlSeconds))
	}
	if _, err := executeRedisCommand(config, args...); err != nil {
		return nil, err
	}

	result := map[string]any{
		"ok":         true,
		"key":        key,
		"connection": config.Name,
	}
	if ttlSeconds > 0 {
		result["ttl_seconds"] = ttlSeconds
	}
	return result, nil
}

func redisDel(key string, connection string) (int, error) {
	config, err := resolveRedisConnectionConfig(connection)
	if err != nil {
		return 0, err
	}
	reply, err := executeRedisCommand(config, "DEL", key)
	if err != nil {
		return 0, err
	}
	return int(reply.integer), nil
}

func redisIncr(key string, delta int, connection string) (int, error) {
	config, err := resolveRedisConnectionConfig(connection)
	if err != nil {
		return 0, err
	}
	reply, err := executeRedisCommand(config, "INCRBY", key, strconv.Itoa(delta))
	if err != nil {
		return 0, err
	}
	return int(reply.integer), nil
}

func redisConnectionHasConfig(connection string) bool {
	name := strings.TrimSpace(connection)
	if name == "" {
		name = redisDefaultConnection
	}
	return strings.TrimSpace(lookupRedisConfigValue(name, "URL")) != "" ||
		strings.TrimSpace(lookupRedisConfigValue(name, "ADDR")) != ""
}

func resolveRedisConnectionConfig(connection string) (redisConnectionConfig, error) {
	name := strings.TrimSpace(connection)
	if name == "" {
		name = redisDefaultConnection
	}

	config := redisConnectionConfig{
		Name:      name,
		TimeoutMS: 5000,
	}
	urlValue := lookupRedisConfigValue(name, "URL")
	if strings.TrimSpace(urlValue) != "" {
		if err := applyRedisURLConfig(&config, urlValue); err != nil {
			return redisConnectionConfig{}, err
		}
	}
	if config.Addr == "" {
		config.Addr = lookupRedisConfigValue(name, "ADDR")
	}
	if username := lookupRedisConfigValue(name, "USERNAME"); username != "" {
		config.Username = username
	}
	if password := lookupRedisConfigValue(name, "PASSWORD"); password != "" {
		config.Password = password
	}
	if dbText := lookupRedisConfigValue(name, "DB"); dbText != "" {
		db, err := strconv.Atoi(strings.TrimSpace(dbText))
		if err != nil {
			return redisConnectionConfig{}, fmt.Errorf("redis connection %q env DB must be an integer: %w", name, err)
		}
		config.DB = db
	}
	if timeoutText := lookupRedisConfigValue(name, "TIMEOUT_MS"); timeoutText != "" {
		timeoutMS, err := strconv.Atoi(strings.TrimSpace(timeoutText))
		if err != nil {
			return redisConnectionConfig{}, fmt.Errorf("redis connection %q env TIMEOUT_MS must be an integer: %w", name, err)
		}
		if timeoutMS > 0 {
			config.TimeoutMS = timeoutMS
		}
	}
	if strings.TrimSpace(config.Addr) == "" {
		return redisConnectionConfig{}, fmt.Errorf("redis connection %q is not configured; set %s", name, strings.Join(redisConfigEnvHints(name), " or "))
	}
	return config, nil
}

func lookupRedisConfigValue(connection string, suffix string) string {
	for _, key := range redisConfigEnvKeys(connection, suffix) {
		if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func redisConfigEnvHints(connection string) []string {
	keys := []string{}
	for _, suffix := range []string{"URL", "ADDR"} {
		keys = append(keys, redisConfigEnvKeys(connection, suffix)...)
	}
	return keys
}

func redisConfigEnvKeys(connection string, suffix string) []string {
	normalized := normalizeRedisConnectionName(connection)
	keys := []string{}
	if normalized == normalizeRedisConnectionName(redisDefaultConnection) {
		keys = append(keys,
			"TSPLAY_REDIS_DEFAULT_"+suffix,
			"TSPLAY_REDIS_"+suffix,
		)
	} else {
		keys = append(keys, "TSPLAY_REDIS_"+normalized+"_"+suffix)
	}
	return keys
}

func normalizeRedisConnectionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = redisDefaultConnection
	}
	builder := strings.Builder{}
	for _, r := range strings.ToUpper(name) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		} else {
			builder.WriteByte('_')
		}
	}
	return strings.Trim(builder.String(), "_")
}

func applyRedisURLConfig(config *redisConnectionConfig, rawURL string) error {
	if config == nil {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse redis url %q: %w", rawURL, err)
	}
	switch parsed.Scheme {
	case "redis":
		config.UseTLS = false
	case "rediss":
		config.UseTLS = true
	default:
		return fmt.Errorf("redis url %q must use redis:// or rediss://", rawURL)
	}
	config.Addr = parsed.Host
	if parsed.User != nil {
		config.Username = parsed.User.Username()
		if password, ok := parsed.User.Password(); ok {
			config.Password = password
		}
	}
	if path := strings.TrimPrefix(parsed.Path, "/"); path != "" {
		db, err := strconv.Atoi(path)
		if err != nil {
			return fmt.Errorf("redis url %q has invalid database %q: %w", rawURL, path, err)
		}
		config.DB = db
	}
	return nil
}

func encodeRedisValue(value any) (string, error) {
	switch typed := value.(type) {
	case nil:
		return "null", nil
	case string:
		return typed, nil
	case []byte:
		return string(typed), nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(typed), nil
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return "", fmt.Errorf("encode redis value: %w", err)
		}
		return string(encoded), nil
	}
}

func executeRedisCommand(config redisConnectionConfig, args ...string) (redisRESPValue, error) {
	conn, reader, writer, err := openRedisConnection(config)
	if err != nil {
		return redisRESPValue{}, err
	}
	defer conn.Close()

	if err := writeRedisCommand(writer, args...); err != nil {
		return redisRESPValue{}, err
	}
	return readRedisValue(reader)
}

func openRedisConnection(config redisConnectionConfig) (net.Conn, *bufio.Reader, *bufio.Writer, error) {
	timeout := time.Duration(config.TimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	dialer := &net.Dialer{Timeout: timeout}

	var (
		conn net.Conn
		err  error
	)
	if config.UseTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", config.Addr, &tls.Config{MinVersion: tls.VersionTLS12})
	} else {
		conn, err = dialer.Dial("tcp", config.Addr)
	}
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connect redis %q at %s: %w", config.Name, config.Addr, err)
	}
	_ = conn.SetDeadline(time.Now().Add(timeout))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	if config.Password != "" || config.Username != "" {
		authArgs := []string{"AUTH"}
		if config.Username != "" {
			authArgs = append(authArgs, config.Username, config.Password)
		} else {
			authArgs = append(authArgs, config.Password)
		}
		if err := writeRedisCommand(writer, authArgs...); err != nil {
			conn.Close()
			return nil, nil, nil, err
		}
		if _, err := readRedisValue(reader); err != nil {
			conn.Close()
			return nil, nil, nil, err
		}
	}
	if config.DB != 0 {
		if err := writeRedisCommand(writer, "SELECT", strconv.Itoa(config.DB)); err != nil {
			conn.Close()
			return nil, nil, nil, err
		}
		if _, err := readRedisValue(reader); err != nil {
			conn.Close()
			return nil, nil, nil, err
		}
	}
	return conn, reader, writer, nil
}

func writeRedisCommand(writer *bufio.Writer, args ...string) error {
	if writer == nil {
		return fmt.Errorf("redis writer is nil")
	}
	if _, err := fmt.Fprintf(writer, "*%d\r\n", len(args)); err != nil {
		return fmt.Errorf("write redis command header: %w", err)
	}
	for _, arg := range args {
		if _, err := fmt.Fprintf(writer, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return fmt.Errorf("write redis command argument: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush redis command: %w", err)
	}
	return nil
}

func readRedisValue(reader *bufio.Reader) (redisRESPValue, error) {
	if reader == nil {
		return redisRESPValue{}, fmt.Errorf("redis reader is nil")
	}
	prefix, err := reader.ReadByte()
	if err != nil {
		return redisRESPValue{}, fmt.Errorf("read redis reply prefix: %w", err)
	}
	switch prefix {
	case '+':
		line, err := readRedisLine(reader)
		if err != nil {
			return redisRESPValue{}, err
		}
		return redisRESPValue{kind: prefix, text: line}, nil
	case '-':
		line, err := readRedisLine(reader)
		if err != nil {
			return redisRESPValue{}, err
		}
		return redisRESPValue{}, fmt.Errorf("redis error: %s", line)
	case ':':
		line, err := readRedisLine(reader)
		if err != nil {
			return redisRESPValue{}, err
		}
		integer, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return redisRESPValue{}, fmt.Errorf("parse redis integer reply %q: %w", line, err)
		}
		return redisRESPValue{kind: prefix, integer: integer}, nil
	case '$':
		line, err := readRedisLine(reader)
		if err != nil {
			return redisRESPValue{}, err
		}
		size, err := strconv.Atoi(line)
		if err != nil {
			return redisRESPValue{}, fmt.Errorf("parse redis bulk size %q: %w", line, err)
		}
		if size < 0 {
			return redisRESPValue{kind: prefix, nil: true}, nil
		}
		payload := make([]byte, size+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return redisRESPValue{}, fmt.Errorf("read redis bulk payload: %w", err)
		}
		return redisRESPValue{kind: prefix, text: string(payload[:size])}, nil
	case '*':
		line, err := readRedisLine(reader)
		if err != nil {
			return redisRESPValue{}, err
		}
		size, err := strconv.Atoi(line)
		if err != nil {
			return redisRESPValue{}, fmt.Errorf("parse redis array size %q: %w", line, err)
		}
		if size < 0 {
			return redisRESPValue{kind: prefix, nil: true}, nil
		}
		values := make([]redisRESPValue, 0, size)
		for i := 0; i < size; i++ {
			value, err := readRedisValue(reader)
			if err != nil {
				return redisRESPValue{}, err
			}
			values = append(values, value)
		}
		return redisRESPValue{kind: prefix, values: values}, nil
	default:
		return redisRESPValue{}, fmt.Errorf("unsupported redis reply prefix %q", prefix)
	}
}

func readRedisLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read redis reply line: %w", err)
	}
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}
