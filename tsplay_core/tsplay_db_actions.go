package tsplay_core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	lua "github.com/yuin/gopher-lua"
)

const dbDefaultConnection = "default"

type dbDialect string

const (
	dbDialectMySQL     dbDialect = "mysql"
	dbDialectPostgres  dbDialect = "postgres"
	dbDialectSQLServer dbDialect = "sqlserver"
	dbDialectOracle    dbDialect = "oracle"
)

type dbInsertConfig struct {
	Table      string
	Row        map[string]any
	Columns    []string
	Connection string
	Driver     string
}

type dbConnectionConfig struct {
	Name       string
	Dialect    dbDialect
	DriverName string
	DSN        string
}

type flowDatabase interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

type flowDatabaseCacheEntry struct {
	db  flowDatabase
	err error
}

var openFlowDatabase = func(driverName string, dsn string) (flowDatabase, error) {
	return sql.Open(driverName, dsn)
}

var flowDatabaseCache sync.Map
var dbSimpleIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_$]*$`)

func db_insert(L *lua.LState) int {
	return runLuaDBInsert(L, "db_insert", "")
}

func runLuaDBInsert(L *lua.LState, action string, forcedDriver string) int {
	values, err := dbInsertValuesFromLua(L, action, forcedDriver)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	config, err := normalizeDBInsertConfig(values, action, forcedDriver)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	result, err := executeDBInsert(context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func runFlowDBInsertStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBInsertValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}

	action := step.Action
	if action == "" {
		action = "db_insert"
	}
	config, err := normalizeDBInsertConfig(values, action, forcedDriver)
	if err != nil {
		return nil, err
	}

	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBInsert(runCtx, config)
}

func dbInsertValuesFromLua(L *lua.LState, action string, forcedDriver string) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("%s requires either a config table or table/row arguments", action)
	}

	first := luaValueToGo(L.CheckAny(1))
	if values, ok := first.(map[string]any); ok {
		if forcedDriver != "" {
			if _, exists := values["driver"]; !exists {
				values["driver"] = forcedDriver
			}
		}
		return values, nil
	}

	values := map[string]any{"table": first}
	if forcedDriver != "" {
		values["driver"] = forcedDriver
	}
	if L.GetTop() >= 2 {
		values["row"] = luaValueToGo(L.CheckAny(2))
	}
	if L.GetTop() >= 3 {
		third := luaValueToGo(L.CheckAny(3))
		if _, ok := third.(string); ok {
			values["connection"] = third
		} else {
			values["columns"] = third
		}
	}
	if L.GetTop() >= 4 {
		values["connection"] = luaValueToGo(L.CheckAny(4))
	}
	if L.GetTop() >= 5 {
		values["driver"] = luaValueToGo(L.CheckAny(5))
	}
	return values, nil
}

func resolvedDBInsertValues(ctx *FlowContext, step FlowStep, forcedDriver string) (map[string]any, error) {
	values := map[string]any{}
	for _, name := range []string{"table", "row", "columns", "connection", "driver"} {
		value, ok, err := flowStepResolvedParam(ctx, step, name)
		if err != nil {
			return nil, err
		}
		if ok {
			values[name] = value
		}
	}
	if forcedDriver != "" {
		if _, exists := values["driver"]; !exists {
			values["driver"] = forcedDriver
		}
	}
	return values, nil
}

func normalizeDBInsertConfig(values map[string]any, action string, forcedDriver string) (dbInsertConfig, error) {
	tableValue, ok := values["table"]
	if !ok || strings.TrimSpace(fmt.Sprint(tableValue)) == "" {
		return dbInsertConfig{}, fmt.Errorf("%s requires table", action)
	}

	rowValue, ok := values["row"]
	if !ok {
		return dbInsertConfig{}, fmt.Errorf("%s requires row", action)
	}
	row, err := objectMapValue(rowValue, "row")
	if err != nil {
		return dbInsertConfig{}, fmt.Errorf("%s %w", action, err)
	}
	if len(row) == 0 {
		return dbInsertConfig{}, fmt.Errorf("%s row must contain at least one field", action)
	}

	var columns []string
	if rawColumns, ok := values["columns"]; ok && rawColumns != nil {
		columns, err = stringListValue(rawColumns)
		if err != nil {
			return dbInsertConfig{}, fmt.Errorf("%s columns %w", action, err)
		}
	}
	columns, err = normalizeDBInsertColumns(action, columns, row)
	if err != nil {
		return dbInsertConfig{}, err
	}

	connection := ""
	if rawConnection, ok := values["connection"]; ok && rawConnection != nil {
		text, ok := rawConnection.(string)
		if !ok {
			return dbInsertConfig{}, fmt.Errorf("%s connection must be a string", action)
		}
		connection = strings.TrimSpace(text)
	}

	driver := strings.TrimSpace(forcedDriver)
	if rawDriver, ok := values["driver"]; ok && rawDriver != nil {
		text, ok := rawDriver.(string)
		if !ok {
			return dbInsertConfig{}, fmt.Errorf("%s driver must be a string", action)
		}
		text = strings.TrimSpace(text)
		if forcedDriver != "" && text != "" && !strings.EqualFold(text, forcedDriver) {
			return dbInsertConfig{}, fmt.Errorf("%s does not allow driver override; use db_insert for non-MySQL targets", action)
		}
		if driver == "" {
			driver = text
		}
	}

	return dbInsertConfig{
		Table:      strings.TrimSpace(fmt.Sprint(tableValue)),
		Row:        row,
		Columns:    columns,
		Connection: connection,
		Driver:     driver,
	}, nil
}

func normalizeDBInsertColumns(action string, columns []string, row map[string]any) ([]string, error) {
	if len(columns) == 0 {
		columns = make([]string, 0, len(row))
		for key := range row {
			columns = append(columns, key)
		}
		sort.Strings(columns)
	}

	normalized := make([]string, 0, len(columns))
	seen := map[string]struct{}{}
	for _, column := range columns {
		column = strings.TrimSpace(column)
		if column == "" {
			return nil, fmt.Errorf("%s columns cannot contain blank names", action)
		}
		if _, ok := seen[column]; ok {
			return nil, fmt.Errorf("%s columns cannot contain duplicate name %q", action, column)
		}
		seen[column] = struct{}{}
		normalized = append(normalized, column)
	}
	if len(normalized) == 0 {
		return nil, fmt.Errorf("%s requires at least one column", action)
	}
	return normalized, nil
}

func executeDBInsert(ctx context.Context, config dbInsertConfig) (map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}

	query, args, err := buildDBInsertQuery(config, connection.Dialect)
	if err != nil {
		return nil, err
	}

	db, err := getFlowDatabase(connection)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db_insert on connection %q failed: %w", connection.Name, err)
	}

	payload := map[string]any{
		"ok":         true,
		"connection": connection.Name,
		"driver":     string(connection.Dialect),
		"table":      config.Table,
		"columns":    append([]string(nil), config.Columns...),
	}
	if rowsAffected, err := result.RowsAffected(); err == nil {
		payload["rows_affected"] = rowsAffected
	}
	if lastInsertID, err := result.LastInsertId(); err == nil {
		payload["last_insert_id"] = lastInsertID
	}
	return payload, nil
}

func getFlowDatabase(config dbConnectionConfig) (flowDatabase, error) {
	key := config.DriverName + "\n" + config.DSN
	if cached, ok := flowDatabaseCache.Load(key); ok {
		entry := cached.(flowDatabaseCacheEntry)
		return entry.db, entry.err
	}

	if err := ensureDatabaseDriverRegistered(config.DriverName, config.Dialect); err != nil {
		return nil, err
	}
	db, err := openFlowDatabase(config.DriverName, config.DSN)
	if err == nil && db != nil {
		db.SetMaxOpenConns(4)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(3 * time.Minute)
	}

	entry := flowDatabaseCacheEntry{db: db, err: err}
	actual, _ := flowDatabaseCache.LoadOrStore(key, entry)
	resolved := actual.(flowDatabaseCacheEntry)
	return resolved.db, resolved.err
}

func ensureDatabaseDriverRegistered(driverName string, dialect dbDialect) error {
	for _, name := range sql.Drivers() {
		if name == driverName {
			return nil
		}
	}
	if dialect == dbDialectOracle {
		return fmt.Errorf("oracle driver %q is not linked into this binary; build with -tags tsplay_oracle to enable the go-ora/v2 driver", driverName)
	}
	if dialect == dbDialectSQLServer {
		return fmt.Errorf("sqlserver driver %q is not linked into this binary; build with -tags tsplay_sqlserver to enable SQL Server support", driverName)
	}
	return fmt.Errorf("database driver %q is not linked into this binary", driverName)
}

func buildDBInsertQuery(config dbInsertConfig, dialect dbDialect) (string, []any, error) {
	tableName, err := quoteDBIdentifier(config.Table, dialect)
	if err != nil {
		return "", nil, fmt.Errorf("db_insert table %w", err)
	}

	columnNames := make([]string, 0, len(config.Columns))
	placeholders := make([]string, 0, len(config.Columns))
	args := make([]any, 0, len(config.Columns))
	for index, column := range config.Columns {
		quotedColumn, err := quoteDBIdentifier(column, dialect)
		if err != nil {
			return "", nil, fmt.Errorf("db_insert column %q %w", column, err)
		}
		value, err := normalizeDBArgument(config.Row[column])
		if err != nil {
			return "", nil, fmt.Errorf("db_insert column %q %w", column, err)
		}
		columnNames = append(columnNames, quotedColumn)
		placeholders = append(placeholders, dbPlaceholder(dialect, index+1))
		args = append(args, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columnNames, ", "),
		strings.Join(placeholders, ", "),
	)
	return query, args, nil
}

func normalizeDBArgument(value any) (any, error) {
	switch typed := value.(type) {
	case []any, []string, map[string]any, map[string]string:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return nil, fmt.Errorf("marshal json value: %w", err)
		}
		return string(encoded), nil
	default:
		return value, nil
	}
}

func dbPlaceholder(dialect dbDialect, index int) string {
	switch dialect {
	case dbDialectPostgres:
		return fmt.Sprintf("$%d", index)
	case dbDialectSQLServer:
		return fmt.Sprintf("@p%d", index)
	case dbDialectOracle:
		return fmt.Sprintf(":%d", index)
	default:
		return "?"
	}
}

func quoteDBIdentifier(name string, dialect dbDialect) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("must be a non-empty identifier")
	}

	parts := strings.Split(name, ".")
	quoted := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return "", fmt.Errorf("must not contain empty identifier segments")
		}
		if dbSimpleIdentifierPattern.MatchString(part) {
			quoted = append(quoted, part)
			continue
		}
		switch dialect {
		case dbDialectMySQL:
			quoted = append(quoted, "`"+strings.ReplaceAll(part, "`", "``")+"`")
		case dbDialectSQLServer:
			quoted = append(quoted, "["+strings.ReplaceAll(part, "]", "]]")+"]")
		default:
			quoted = append(quoted, `"`+strings.ReplaceAll(part, `"`, `""`)+`"`)
		}
	}
	return strings.Join(quoted, "."), nil
}

func resolveDBConnectionConfig(connection string, driverHint string) (dbConnectionConfig, error) {
	name := strings.TrimSpace(connection)
	if name == "" {
		name = dbDefaultConnection
	}

	driverText := strings.TrimSpace(driverHint)
	if configuredDriver := lookupDBConfigValue(name, "DRIVER"); configuredDriver != "" {
		driverText = configuredDriver
	}

	if dsn := lookupDBConfigValue(name, "DSN"); dsn != "" {
		dialect, driverName, err := normalizeDBDriver(driverText)
		if err != nil {
			return dbConnectionConfig{}, fmt.Errorf("database connection %q requires DRIVER when DSN is used: %w", name, err)
		}
		return dbConnectionConfig{Name: name, Dialect: dialect, DriverName: driverName, DSN: dsn}, nil
	}

	if rawURL := lookupDBConfigValue(name, "URL"); rawURL != "" {
		config, err := resolveDBURLConnectionConfig(name, driverText, rawURL)
		if err != nil {
			return dbConnectionConfig{}, err
		}
		return config, nil
	}

	if hasStructuredDBConfig(name) || driverText != "" {
		config, err := resolveStructuredDBConnectionConfig(name, driverText)
		if err == nil {
			return config, nil
		}
		if hasStructuredDBConfig(name) {
			return dbConnectionConfig{}, err
		}
	}

	if config, ok, err := resolveLegacyMySQLConnectionConfig(name, driverText); ok {
		return config, err
	}

	return dbConnectionConfig{}, fmt.Errorf("database connection %q is not configured; set %s", name, strings.Join(dbConfigEnvHints(name), " or "))
}

func resolveDBURLConnectionConfig(name string, driverHint string, rawURL string) (dbConnectionConfig, error) {
	driverText := strings.TrimSpace(driverHint)
	if driverText == "" {
		inferred, err := inferDBDriverFromURL(rawURL)
		if err != nil {
			return dbConnectionConfig{}, fmt.Errorf("database connection %q url %w", name, err)
		}
		driverText = inferred
	}

	dialect, driverName, err := normalizeDBDriver(driverText)
	if err != nil {
		return dbConnectionConfig{}, err
	}

	dsn, err := dbDSNFromURL(dialect, rawURL)
	if err != nil {
		return dbConnectionConfig{}, fmt.Errorf("database connection %q url %w", name, err)
	}
	return dbConnectionConfig{Name: name, Dialect: dialect, DriverName: driverName, DSN: dsn}, nil
}

func resolveStructuredDBConnectionConfig(name string, driverHint string) (dbConnectionConfig, error) {
	dialect, driverName, err := normalizeDBDriver(driverHint)
	if err != nil {
		return dbConnectionConfig{}, fmt.Errorf("database connection %q requires DRIVER when URL/DSN is not set: %w", name, err)
	}

	dsn, err := structuredDBDSN(name, dialect)
	if err != nil {
		return dbConnectionConfig{}, err
	}
	return dbConnectionConfig{Name: name, Dialect: dialect, DriverName: driverName, DSN: dsn}, nil
}

func normalizeDBDriver(raw string) (dbDialect, string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "mysql", "mariadb":
		return dbDialectMySQL, "mysql", nil
	case "postgres", "postgresql", "pgsql", "pgx", "pq":
		return dbDialectPostgres, "postgres", nil
	case "sqlserver", "mssql":
		return dbDialectSQLServer, "sqlserver", nil
	case "oracle", "gora", "goora", "go-ora", "go_ora", "godror":
		return dbDialectOracle, "oracle", nil
	default:
		return "", "", fmt.Errorf("unsupported driver %q; expected one of mysql, postgres, sqlserver, or oracle", raw)
	}
}

func inferDBDriverFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url %q: %w", rawURL, err)
	}
	switch strings.ToLower(strings.TrimSpace(parsed.Scheme)) {
	case "mysql", "mariadb":
		return "mysql", nil
	case "postgres", "postgresql", "pgsql":
		return "postgres", nil
	case "sqlserver", "mssql":
		return "sqlserver", nil
	case "oracle", "gora", "goora", "go-ora", "go_ora", "godror":
		return "oracle", nil
	default:
		return "", fmt.Errorf("cannot infer driver from scheme %q", parsed.Scheme)
	}
}

func dbDSNFromURL(dialect dbDialect, rawURL string) (string, error) {
	switch dialect {
	case dbDialectMySQL:
		return mysqlDSNFromURL(rawURL)
	case dbDialectPostgres, dbDialectSQLServer:
		return rawURL, nil
	case dbDialectOracle:
		return oracleDSNFromURL(rawURL)
	default:
		return "", fmt.Errorf("unsupported driver %q", dialect)
	}
}

func structuredDBDSN(connection string, dialect dbDialect) (string, error) {
	switch dialect {
	case dbDialectMySQL:
		return structuredMySQLDSN(connection)
	case dbDialectPostgres:
		return structuredPostgresDSN(connection)
	case dbDialectSQLServer:
		return structuredSQLServerDSN(connection)
	case dbDialectOracle:
		return structuredOracleDSN(connection)
	default:
		return "", fmt.Errorf("unsupported driver %q", dialect)
	}
}

func mysqlDSNFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse mysql url %q: %w", rawURL, err)
	}
	config := mysqlDriver.NewConfig()
	if parsed.User != nil {
		config.User = parsed.User.Username()
		if password, ok := parsed.User.Password(); ok {
			config.Passwd = password
		}
	}

	socket := strings.TrimSpace(parsed.Query().Get("socket"))
	if socket != "" {
		config.Net = "unix"
		config.Addr = socket
	} else {
		host := parsed.Hostname()
		port := parsed.Port()
		if host == "" {
			return "", fmt.Errorf("mysql url %q requires host or query param socket", rawURL)
		}
		if port == "" {
			port = "3306"
		}
		config.Net = "tcp"
		config.Addr = net.JoinHostPort(host, port)
	}
	config.DBName = strings.TrimPrefix(parsed.Path, "/")
	params := map[string]string{}
	for key, values := range parsed.Query() {
		if key == "socket" || len(values) == 0 {
			continue
		}
		params[key] = values[0]
	}
	if len(params) > 0 {
		config.Params = params
	}
	return config.FormatDSN(), nil
}

func structuredMySQLDSN(connection string) (string, error) {
	config := mysqlDriver.NewConfig()
	config.User = lookupDBConfigValue(connection, "USER")
	config.Passwd = lookupDBConfigValue(connection, "PASSWORD")

	socket := lookupDBConfigValue(connection, "SOCKET")
	if socket != "" {
		config.Net = "unix"
		config.Addr = socket
	} else {
		host := lookupDBConfigValue(connection, "HOST")
		if host == "" {
			return "", fmt.Errorf("database connection %q mysql requires HOST or SOCKET", connection)
		}
		port := lookupDBConfigValue(connection, "PORT")
		if port == "" {
			port = "3306"
		}
		config.Net = "tcp"
		config.Addr = net.JoinHostPort(host, port)
	}
	config.DBName = lookupDBConfigValue(connection, "DATABASE")
	if params := lookupDBConfigValue(connection, "PARAMS"); params != "" {
		queryValues, err := url.ParseQuery(params)
		if err != nil {
			return "", fmt.Errorf("database connection %q mysql PARAMS %w", connection, err)
		}
		config.Params = map[string]string{}
		for key, values := range queryValues {
			if len(values) > 0 {
				config.Params[key] = values[0]
			}
		}
	}
	return config.FormatDSN(), nil
}

func structuredPostgresDSN(connection string) (string, error) {
	host := lookupDBConfigValue(connection, "HOST")
	if host == "" {
		return "", fmt.Errorf("database connection %q postgres requires HOST", connection)
	}
	port := lookupDBConfigValue(connection, "PORT")
	if port == "" {
		port = "5432"
	}

	dsn := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(host, port),
		Path:   lookupDBConfigValue(connection, "DATABASE"),
	}
	user := lookupDBConfigValue(connection, "USER")
	password := lookupDBConfigValue(connection, "PASSWORD")
	if user != "" {
		if password != "" {
			dsn.User = url.UserPassword(user, password)
		} else {
			dsn.User = url.User(user)
		}
	}
	query := url.Values{}
	if sslmode := lookupDBConfigValue(connection, "SSLMODE"); sslmode != "" {
		query.Set("sslmode", sslmode)
	}
	if params := lookupDBConfigValue(connection, "PARAMS"); params != "" {
		values, err := url.ParseQuery(params)
		if err != nil {
			return "", fmt.Errorf("database connection %q postgres PARAMS %w", connection, err)
		}
		for key, items := range values {
			for _, item := range items {
				query.Add(key, item)
			}
		}
	}
	if len(query) > 0 {
		dsn.RawQuery = query.Encode()
	}
	return dsn.String(), nil
}

func structuredSQLServerDSN(connection string) (string, error) {
	host := lookupDBConfigValue(connection, "HOST")
	if host == "" {
		return "", fmt.Errorf("database connection %q sqlserver requires HOST", connection)
	}
	port := lookupDBConfigValue(connection, "PORT")
	if port == "" {
		port = "1433"
	}

	dsn := &url.URL{
		Scheme: "sqlserver",
		Host:   net.JoinHostPort(host, port),
	}
	user := lookupDBConfigValue(connection, "USER")
	password := lookupDBConfigValue(connection, "PASSWORD")
	if user != "" {
		if password != "" {
			dsn.User = url.UserPassword(user, password)
		} else {
			dsn.User = url.User(user)
		}
	}
	query := url.Values{}
	if database := lookupDBConfigValue(connection, "DATABASE"); database != "" {
		query.Set("database", database)
	}
	if instance := lookupDBConfigValue(connection, "INSTANCE"); instance != "" {
		query.Set("instance", instance)
	}
	if params := lookupDBConfigValue(connection, "PARAMS"); params != "" {
		values, err := url.ParseQuery(params)
		if err != nil {
			return "", fmt.Errorf("database connection %q sqlserver PARAMS %w", connection, err)
		}
		for key, items := range values {
			for _, item := range items {
				query.Add(key, item)
			}
		}
	}
	if len(query) > 0 {
		dsn.RawQuery = query.Encode()
	}
	return dsn.String(), nil
}

func structuredOracleDSN(connection string) (string, error) {
	connectString := lookupDBConfigValue(connection, "CONNECT_STRING")
	host := lookupDBConfigValue(connection, "HOST")
	port := lookupDBConfigValue(connection, "PORT")
	service := lookupDBConfigValue(connection, "SERVICE")
	if service == "" {
		service = lookupDBConfigValue(connection, "DATABASE")
	}
	user := lookupDBConfigValue(connection, "USER")
	password := lookupDBConfigValue(connection, "PASSWORD")
	options := map[string]string{}

	if connectString != "" {
		if parsedHost, parsedPort, parsedService, ok := parseOracleConnectString(connectString); ok {
			host = parsedHost
			port = parsedPort
			service = parsedService
		} else {
			options["connStr"] = connectString
		}
	}
	if host == "" && options["connStr"] == "" {
		return "", fmt.Errorf("database connection %q oracle requires CONNECT_STRING or HOST", connection)
	}
	if host != "" {
		if port == "" {
			port = "1521"
		}
		if service == "" {
			return "", fmt.Errorf("database connection %q oracle requires SERVICE, DATABASE, or CONNECT_STRING", connection)
		}
	}
	if sid := lookupDBConfigValue(connection, "SID"); sid != "" {
		options["SID"] = sid
	}
	if params := lookupDBConfigValue(connection, "PARAMS"); params != "" {
		values, err := url.ParseQuery(params)
		if err != nil {
			return "", fmt.Errorf("database connection %q oracle PARAMS %w", connection, err)
		}
		for key, items := range values {
			if len(items) > 0 {
				options[key] = items[0]
			}
		}
	}
	return buildOracleURL(user, password, host, port, service, options), nil
}

func oracleDSNFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse oracle url %q: %w", rawURL, err)
	}
	host := parsed.Hostname()
	port := parsed.Port()
	service := strings.TrimPrefix(parsed.Path, "/")
	if port == "" && host != "" {
		port = "1521"
	}
	options := map[string]string{}
	connectString := strings.TrimSpace(parsed.Query().Get("connectString"))
	if connectString == "" {
		connectString = strings.TrimSpace(parsed.Query().Get("connStr"))
	}
	if connectString != "" {
		if parsedHost, parsedPort, parsedService, ok := parseOracleConnectString(connectString); ok && host == "" && service == "" {
			host = parsedHost
			port = parsedPort
			service = parsedService
		} else {
			options["connStr"] = connectString
		}
	}
	if parsed.User != nil {
		// handled below when constructing the final URL
	}
	for key, values := range parsed.Query() {
		if (strings.EqualFold(key, "connectString") || key == "connStr") || len(values) == 0 {
			continue
		}
		options[key] = values[0]
	}
	if host == "" && options["connStr"] == "" {
		return "", fmt.Errorf("oracle url %q requires host and service path or query param connStr", rawURL)
	}
	if host != "" && service == "" {
		return "", fmt.Errorf("oracle url %q requires service path when host is present", rawURL)
	}
	user := ""
	password := ""
	if parsed.User != nil {
		user = parsed.User.Username()
		password, _ = parsed.User.Password()
	}
	return buildOracleURL(user, password, host, port, service, options), nil
}

func buildOracleURL(user string, password string, host string, port string, service string, options map[string]string) string {
	if host == "" && port == "" {
		port = "0"
	}
	if host != "" && port == "" {
		port = "1521"
	}

	dsn := &url.URL{Scheme: "oracle"}
	if user != "" || password != "" {
		if password != "" {
			dsn.User = url.UserPassword(user, password)
		} else {
			dsn.User = url.User(user)
		}
	}
	if host != "" || port != "" {
		dsn.Host = net.JoinHostPort(host, port)
	}
	if service != "" {
		dsn.Path = "/" + service
	} else {
		dsn.Path = "/"
	}

	query := url.Values{}
	for key, value := range options {
		if strings.TrimSpace(key) == "" || value == "" {
			continue
		}
		query.Set(key, value)
	}
	if len(query) > 0 {
		dsn.RawQuery = query.Encode()
	}
	return dsn.String()
}

func parseOracleConnectString(connectString string) (string, string, string, bool) {
	connectString = strings.TrimSpace(connectString)
	if connectString == "" || strings.Contains(connectString, "(") {
		return "", "", "", false
	}

	separator := strings.LastIndex(connectString, "/")
	if separator <= 0 || separator == len(connectString)-1 {
		return "", "", "", false
	}

	hostPort := connectString[:separator]
	service := strings.TrimSpace(connectString[separator+1:])
	if service == "" {
		return "", "", "", false
	}

	host, port, err := net.SplitHostPort(hostPort)
	if err == nil {
		if port == "" {
			port = "1521"
		}
		return host, port, service, true
	}
	if strings.Contains(hostPort, ":") {
		return "", "", "", false
	}
	return hostPort, "1521", service, true
}

func lookupDBConfigValue(connection string, suffix string) string {
	for _, key := range dbConfigEnvKeys(connection, suffix) {
		if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func dbConfigEnvHints(connection string) []string {
	keys := []string{}
	for _, suffix := range []string{"DRIVER", "DSN", "URL"} {
		keys = append(keys, dbConfigEnvKeys(connection, suffix)...)
	}
	return keys
}

func dbConfigEnvKeys(connection string, suffix string) []string {
	normalized := normalizeDBConnectionName(connection)
	if normalized == normalizeDBConnectionName(dbDefaultConnection) {
		return []string{
			"TSPLAY_DB_DEFAULT_" + suffix,
			"TSPLAY_DB_" + suffix,
		}
	}
	return []string{"TSPLAY_DB_" + normalized + "_" + suffix}
}

func normalizeDBConnectionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = dbDefaultConnection
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

func hasStructuredDBConfig(connection string) bool {
	for _, suffix := range []string{
		"HOST",
		"PORT",
		"USER",
		"PASSWORD",
		"DATABASE",
		"SOCKET",
		"INSTANCE",
		"SERVICE",
		"CONNECT_STRING",
		"PARAMS",
		"SSLMODE",
		"LIB_DIR",
	} {
		if lookupDBConfigValue(connection, suffix) != "" {
			return true
		}
	}
	return false
}

func resolveLegacyMySQLConnectionConfig(connection string, driverHint string) (dbConnectionConfig, bool, error) {
	if driverHint != "" && !strings.EqualFold(strings.TrimSpace(driverHint), string(dbDialectMySQL)) {
		return dbConnectionConfig{}, false, nil
	}
	if !legacyMySQLConnectionHasConfig(connection) {
		return dbConnectionConfig{}, false, nil
	}

	if rawURL := lookupLegacyMySQLConfigValue(connection, "URL"); rawURL != "" {
		dsn, err := mysqlDSNFromURL(rawURL)
		if err != nil {
			return dbConnectionConfig{}, true, fmt.Errorf("legacy mysql connection %q url %w", connection, err)
		}
		return dbConnectionConfig{
			Name:       connection,
			Dialect:    dbDialectMySQL,
			DriverName: "mysql",
			DSN:        dsn,
		}, true, nil
	}

	config := mysqlDriver.NewConfig()
	config.User = lookupLegacyMySQLConfigValue(connection, "USER")
	if config.User == "" {
		config.User = lookupLegacyMySQLConfigValue(connection, "USERNAME")
	}
	config.Passwd = lookupLegacyMySQLConfigValue(connection, "PASSWORD")
	config.DBName = lookupLegacyMySQLConfigValue(connection, "DATABASE")
	if config.DBName == "" {
		config.DBName = lookupLegacyMySQLConfigValue(connection, "DB")
	}

	socket := lookupLegacyMySQLConfigValue(connection, "SOCKET")
	if socket != "" {
		config.Net = "unix"
		config.Addr = socket
	} else if addr := lookupLegacyMySQLConfigValue(connection, "ADDR"); addr != "" {
		config.Net = "tcp"
		config.Addr = addr
	} else if host := lookupLegacyMySQLConfigValue(connection, "HOST"); host != "" {
		port := lookupLegacyMySQLConfigValue(connection, "PORT")
		if port == "" {
			port = "3306"
		}
		config.Net = "tcp"
		config.Addr = net.JoinHostPort(host, port)
	} else {
		return dbConnectionConfig{}, true, fmt.Errorf("legacy mysql connection %q is not configured; set %s", connection, strings.Join(legacyMySQLConfigEnvHints(connection), " or "))
	}

	return dbConnectionConfig{
		Name:       connection,
		Dialect:    dbDialectMySQL,
		DriverName: "mysql",
		DSN:        config.FormatDSN(),
	}, true, nil
}

func legacyMySQLConnectionHasConfig(connection string) bool {
	for _, suffix := range []string{"URL", "ADDR", "HOST", "SOCKET"} {
		if lookupLegacyMySQLConfigValue(connection, suffix) != "" {
			return true
		}
	}
	return false
}

func lookupLegacyMySQLConfigValue(connection string, suffix string) string {
	for _, key := range legacyMySQLConfigEnvKeys(connection, suffix) {
		if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func legacyMySQLConfigEnvHints(connection string) []string {
	keys := []string{}
	for _, suffix := range []string{"URL", "ADDR", "HOST", "SOCKET"} {
		keys = append(keys, legacyMySQLConfigEnvKeys(connection, suffix)...)
	}
	return keys
}

func legacyMySQLConfigEnvKeys(connection string, suffix string) []string {
	normalized := normalizeLegacyMySQLConnectionName(connection)
	if normalized == normalizeLegacyMySQLConnectionName(dbDefaultConnection) {
		return []string{
			"TSPLAY_MYSQL_DEFAULT_" + suffix,
			"TSPLAY_MYSQL_" + suffix,
		}
	}
	return []string{"TSPLAY_MYSQL_" + normalized + "_" + suffix}
}

func normalizeLegacyMySQLConnectionName(name string) string {
	return normalizeDBConnectionName(name)
}
