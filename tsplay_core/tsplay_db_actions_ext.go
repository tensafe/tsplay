package tsplay_core

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type dbInsertManyConfig struct {
	Table      string
	Rows       []map[string]any
	Columns    []string
	Connection string
	Driver     string
	Returning  []string
	TimeoutMS  int
}

type dbUpsertConfig struct {
	Table         string
	Row           map[string]any
	Columns       []string
	KeyColumns    []string
	UpdateColumns []string
	Connection    string
	Driver        string
	Returning     []string
	TimeoutMS     int
	DoNothing     bool
}

type dbStatementConfig struct {
	SQL        string
	Args       []any
	Connection string
	Driver     string
	TimeoutMS  int
}

type dbRuntimeSettings struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	QueryTimeout    time.Duration
}

type flowDBTransactionScope struct {
	mu       sync.Mutex
	sessions map[string]*flowDBTransactionSession
}

type flowDBTransactionSession struct {
	connection dbConnectionConfig
	tx         flowDBTransaction
}

func validateDBTransactionFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use steps and optional timeout", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{"steps": true, "timeout": true, "timeout_ms": true, "timeout_seconds": true}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if name == "steps" {
			continue
		}
		if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if len(step.Steps) == 0 {
		return fmt.Errorf("step %s action %q requires nested steps", stepPath, step.Action)
	}
	if value, ok := step.param("timeout"); ok && len(flowReferences(value)) == 0 {
		timeoutMS, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "timeout", err)
		}
		if timeoutMS < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "timeout")
		}
	}
	if value, ok := step.param("timeout_ms"); ok && len(flowReferences(value)) == 0 {
		timeoutMS, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "timeout_ms", err)
		}
		if timeoutMS < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "timeout_ms")
		}
	}
	if value, ok := step.param("timeout_seconds"); ok && len(flowReferences(value)) == 0 {
		timeoutSeconds, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "timeout_seconds", err)
		}
		if timeoutSeconds < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "timeout_seconds")
		}
	}
	return validateFlowStepSequence(step.Steps, copyKnownVars(knownVars), stepPath)
}

func db_insert_many(L *lua.LState) int {
	values, err := dbInsertManyValuesFromLua(L, "db_insert_many", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err := normalizeDBInsertManyConfig(values, "db_insert_many", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeDBInsertManyWithFlow(nil, context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func db_upsert(L *lua.LState) int {
	values, err := dbUpsertValuesFromLua(L, "db_upsert", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err := normalizeDBUpsertConfig(values, "db_upsert", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeDBUpsertWithFlow(nil, context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func db_query(L *lua.LState) int {
	values, err := dbStatementValuesFromLua(L, "db_query", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err := normalizeDBStatementConfig(values, "db_query", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeDBQueryWithFlow(nil, context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func db_query_one(L *lua.LState) int {
	values, err := dbStatementValuesFromLua(L, "db_query_one", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err := normalizeDBStatementConfig(values, "db_query_one", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeDBQueryOneWithFlow(nil, context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func db_execute(L *lua.LState) int {
	values, err := dbStatementValuesFromLua(L, "db_execute", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err := normalizeDBStatementConfig(values, "db_execute", "")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeDBExecuteWithFlow(nil, context.Background(), config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func runFlowDBInsertManyStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBInsertManyValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}
	config, err := normalizeDBInsertManyConfig(values, step.Action, forcedDriver)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBInsertManyWithFlow(ctx, runCtx, config)
}

func runFlowDBUpsertStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBUpsertValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}
	config, err := normalizeDBUpsertConfig(values, step.Action, forcedDriver)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBUpsertWithFlow(ctx, runCtx, config)
}

func runFlowDBQueryStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBStatementValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}
	config, err := normalizeDBStatementConfig(values, step.Action, forcedDriver)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBQueryWithFlow(ctx, runCtx, config)
}

func runFlowDBQueryOneStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBStatementValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}
	config, err := normalizeDBStatementConfig(values, step.Action, forcedDriver)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBQueryOneWithFlow(ctx, runCtx, config)
}

func runFlowDBExecuteStep(ctx *FlowContext, step FlowStep, forcedDriver string) (any, error) {
	values, err := resolvedDBStatementValues(ctx, step, forcedDriver)
	if err != nil {
		return nil, err
	}
	config, err := normalizeDBStatementConfig(values, step.Action, forcedDriver)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeDBExecuteWithFlow(ctx, runCtx, config)
}

func runFlowDBTransactionStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, []FlowStepTrace, error) {
	if ctx == nil {
		return nil, nil, fmt.Errorf("db_transaction requires flow context")
	}
	if ctx.DBTransaction != nil {
		return nil, nil, fmt.Errorf("nested db_transaction is not supported")
	}
	timeoutMS, err := resolvedDBTimeoutForStep(ctx, step, step.Action)
	if err != nil {
		return nil, nil, err
	}

	runCtx := ctx.Context
	if runCtx == nil {
		runCtx = context.Background()
	}
	cancel := func() {}
	if timeoutMS > 0 {
		runCtx, cancel = context.WithTimeout(runCtx, time.Duration(timeoutMS)*time.Millisecond)
	}
	defer cancel()

	scope := &flowDBTransactionScope{sessions: map[string]*flowDBTransactionSession{}}
	child := *ctx
	child.Context = runCtx
	child.DBTransaction = scope

	children, err := runFlowStepSequence(L, &child, step.Steps, stepPath, 0, 0)
	if err != nil {
		if rollbackErr := scope.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("%w (rollback failed: %v)", err, rollbackErr)
		}
		return nil, children, err
	}
	if err := scope.Commit(); err != nil {
		return nil, children, err
	}
	return map[string]any{
		"ok":           true,
		"connections":  scope.connectionNames(),
		"transactions": len(scope.connectionNames()),
	}, children, nil
}

func dbInsertManyValuesFromLua(L *lua.LState, action string, forcedDriver string) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("%s requires either a config table or table/rows arguments", action)
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
	if L.GetTop() >= 2 {
		values["rows"] = luaValueToGo(L.CheckAny(2))
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

func resolvedDBInsertManyValues(ctx *FlowContext, step FlowStep, forcedDriver string) (map[string]any, error) {
	return resolvedDBValues(ctx, step, forcedDriver, "table", "rows", "columns", "connection", "driver", "returning", "timeout", "timeout_ms", "timeout_seconds")
}

func normalizeDBInsertManyConfig(values map[string]any, action string, forcedDriver string) (dbInsertManyConfig, error) {
	tableValue, ok := values["table"]
	if !ok || strings.TrimSpace(fmt.Sprint(tableValue)) == "" {
		return dbInsertManyConfig{}, fmt.Errorf("%s requires table", action)
	}
	rowsValue, ok := values["rows"]
	if !ok {
		return dbInsertManyConfig{}, fmt.Errorf("%s requires rows", action)
	}
	rows, err := objectListValue(rowsValue, "rows")
	if err != nil {
		return dbInsertManyConfig{}, fmt.Errorf("%s %w", action, err)
	}
	if len(rows) == 0 {
		return dbInsertManyConfig{}, fmt.Errorf("%s rows must contain at least one object", action)
	}

	var columns []string
	if rawColumns, ok := values["columns"]; ok && rawColumns != nil {
		columns, err = stringListValue(rawColumns)
		if err != nil {
			return dbInsertManyConfig{}, fmt.Errorf("%s columns %w", action, err)
		}
	}
	columns, err = normalizeDBInsertManyColumns(action, columns, rows)
	if err != nil {
		return dbInsertManyConfig{}, err
	}

	connection, driver, returning, timeoutMS, err := normalizeDBCommonWriteOptions(values, action, forcedDriver)
	if err != nil {
		return dbInsertManyConfig{}, err
	}
	return dbInsertManyConfig{
		Table:      strings.TrimSpace(fmt.Sprint(tableValue)),
		Rows:       rows,
		Columns:    columns,
		Connection: connection,
		Driver:     driver,
		Returning:  returning,
		TimeoutMS:  timeoutMS,
	}, nil
}

func dbUpsertValuesFromLua(L *lua.LState, action string, forcedDriver string) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("%s requires either a config table or table/row/key_columns arguments", action)
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
	if L.GetTop() >= 2 {
		values["row"] = luaValueToGo(L.CheckAny(2))
	}
	if L.GetTop() >= 3 {
		values["key_columns"] = luaValueToGo(L.CheckAny(3))
	}
	if L.GetTop() >= 4 {
		fourth := luaValueToGo(L.CheckAny(4))
		if _, ok := fourth.(string); ok {
			values["connection"] = fourth
		} else {
			values["update_columns"] = fourth
		}
	}
	if L.GetTop() >= 5 {
		values["connection"] = luaValueToGo(L.CheckAny(5))
	}
	if L.GetTop() >= 6 {
		values["driver"] = luaValueToGo(L.CheckAny(6))
	}
	return values, nil
}

func resolvedDBUpsertValues(ctx *FlowContext, step FlowStep, forcedDriver string) (map[string]any, error) {
	return resolvedDBValues(ctx, step, forcedDriver, "table", "row", "columns", "key_columns", "update_columns", "do_nothing", "connection", "driver", "returning", "timeout", "timeout_ms", "timeout_seconds")
}

func normalizeDBUpsertConfig(values map[string]any, action string, forcedDriver string) (dbUpsertConfig, error) {
	tableValue, ok := values["table"]
	if !ok || strings.TrimSpace(fmt.Sprint(tableValue)) == "" {
		return dbUpsertConfig{}, fmt.Errorf("%s requires table", action)
	}
	rowValue, ok := values["row"]
	if !ok {
		return dbUpsertConfig{}, fmt.Errorf("%s requires row", action)
	}
	row, err := objectMapValue(rowValue, "row")
	if err != nil {
		return dbUpsertConfig{}, fmt.Errorf("%s %w", action, err)
	}
	if len(row) == 0 {
		return dbUpsertConfig{}, fmt.Errorf("%s row must contain at least one field", action)
	}
	keyColumnsValue, ok := values["key_columns"]
	if !ok {
		return dbUpsertConfig{}, fmt.Errorf("%s requires key_columns", action)
	}
	keyColumns, err := stringListValue(keyColumnsValue)
	if err != nil {
		return dbUpsertConfig{}, fmt.Errorf("%s key_columns %w", action, err)
	}
	if len(keyColumns) == 0 {
		return dbUpsertConfig{}, fmt.Errorf("%s key_columns must contain at least one column", action)
	}

	var columns []string
	if rawColumns, ok := values["columns"]; ok && rawColumns != nil {
		columns, err = stringListValue(rawColumns)
		if err != nil {
			return dbUpsertConfig{}, fmt.Errorf("%s columns %w", action, err)
		}
	}
	columns, err = normalizeDBInsertColumns(action, columns, row)
	if err != nil {
		return dbUpsertConfig{}, err
	}
	columnSet := map[string]struct{}{}
	for _, column := range columns {
		columnSet[column] = struct{}{}
	}
	for _, keyColumn := range keyColumns {
		if _, ok := row[keyColumn]; !ok {
			return dbUpsertConfig{}, fmt.Errorf("%s row is missing key column %q", action, keyColumn)
		}
		if _, ok := columnSet[keyColumn]; !ok {
			return dbUpsertConfig{}, fmt.Errorf("%s key column %q must appear in columns", action, keyColumn)
		}
	}

	var updateColumns []string
	if rawUpdateColumns, ok := values["update_columns"]; ok && rawUpdateColumns != nil {
		updateColumns, err = stringListValue(rawUpdateColumns)
		if err != nil {
			return dbUpsertConfig{}, fmt.Errorf("%s update_columns %w", action, err)
		}
	}
	if len(updateColumns) == 0 {
		keySet := map[string]struct{}{}
		for _, keyColumn := range keyColumns {
			keySet[keyColumn] = struct{}{}
		}
		for _, column := range columns {
			if _, isKey := keySet[column]; !isKey {
				updateColumns = append(updateColumns, column)
			}
		}
	}
	doNothing := false
	if rawDoNothing, ok := values["do_nothing"]; ok && rawDoNothing != nil {
		doNothing, err = boolParam(rawDoNothing)
		if err != nil {
			return dbUpsertConfig{}, fmt.Errorf("%s do_nothing %w", action, err)
		}
	}
	if len(updateColumns) == 0 {
		doNothing = true
	}
	for _, updateColumn := range updateColumns {
		if _, ok := columnSet[updateColumn]; !ok {
			return dbUpsertConfig{}, fmt.Errorf("%s update column %q must appear in columns", action, updateColumn)
		}
	}

	connection, driver, returning, timeoutMS, err := normalizeDBCommonWriteOptions(values, action, forcedDriver)
	if err != nil {
		return dbUpsertConfig{}, err
	}
	return dbUpsertConfig{
		Table:         strings.TrimSpace(fmt.Sprint(tableValue)),
		Row:           row,
		Columns:       columns,
		KeyColumns:    append([]string(nil), keyColumns...),
		UpdateColumns: append([]string(nil), updateColumns...),
		Connection:    connection,
		Driver:        driver,
		Returning:     returning,
		TimeoutMS:     timeoutMS,
		DoNothing:     doNothing,
	}, nil
}

func dbStatementValuesFromLua(L *lua.LState, action string, forcedDriver string) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("%s requires either a config table or sql/args arguments", action)
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
	values := map[string]any{"sql": first}
	if L.GetTop() >= 2 {
		values["args"] = luaValueToGo(L.CheckAny(2))
	}
	if L.GetTop() >= 3 {
		values["connection"] = luaValueToGo(L.CheckAny(3))
	}
	if L.GetTop() >= 4 {
		values["driver"] = luaValueToGo(L.CheckAny(4))
	}
	return values, nil
}

func resolvedDBStatementValues(ctx *FlowContext, step FlowStep, forcedDriver string) (map[string]any, error) {
	return resolvedDBValues(ctx, step, forcedDriver, "sql", "args", "connection", "driver", "timeout", "timeout_ms", "timeout_seconds")
}

func normalizeDBStatementConfig(values map[string]any, action string, forcedDriver string) (dbStatementConfig, error) {
	sqlValue, ok := values["sql"]
	if !ok || strings.TrimSpace(fmt.Sprint(sqlValue)) == "" {
		return dbStatementConfig{}, fmt.Errorf("%s requires sql", action)
	}
	args, err := normalizeDBStatementArgs(values["args"])
	if err != nil {
		return dbStatementConfig{}, fmt.Errorf("%s args %w", action, err)
	}
	connection, driver, _, timeoutMS, err := normalizeDBCommonWriteOptions(values, action, forcedDriver)
	if err != nil {
		return dbStatementConfig{}, err
	}
	return dbStatementConfig{
		SQL:        strings.TrimSpace(fmt.Sprint(sqlValue)),
		Args:       args,
		Connection: connection,
		Driver:     driver,
		TimeoutMS:  timeoutMS,
	}, nil
}

func resolvedDBValues(ctx *FlowContext, step FlowStep, forcedDriver string, names ...string) (map[string]any, error) {
	values := map[string]any{}
	for _, name := range names {
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

func normalizeDBCommonWriteOptions(values map[string]any, action string, forcedDriver string) (string, string, []string, int, error) {
	connection := ""
	if rawConnection, ok := values["connection"]; ok && rawConnection != nil {
		text, ok := rawConnection.(string)
		if !ok {
			return "", "", nil, 0, fmt.Errorf("%s connection must be a string", action)
		}
		connection = strings.TrimSpace(text)
	}
	driver := strings.TrimSpace(forcedDriver)
	if rawDriver, ok := values["driver"]; ok && rawDriver != nil {
		text, ok := rawDriver.(string)
		if !ok {
			return "", "", nil, 0, fmt.Errorf("%s driver must be a string", action)
		}
		text = strings.TrimSpace(text)
		if forcedDriver != "" && text != "" && !strings.EqualFold(text, forcedDriver) {
			return "", "", nil, 0, fmt.Errorf("%s does not allow driver override; use db_insert for non-MySQL targets", action)
		}
		if driver == "" {
			driver = text
		}
	}
	returning, err := normalizeDBReturningColumns(values, action)
	if err != nil {
		return "", "", nil, 0, err
	}
	timeoutMS, err := normalizeDBTimeoutMS(values, action)
	if err != nil {
		return "", "", nil, 0, err
	}
	return connection, driver, returning, timeoutMS, nil
}

func normalizeDBReturningColumns(values map[string]any, action string) ([]string, error) {
	rawReturning, ok := values["returning"]
	if !ok || rawReturning == nil {
		return nil, nil
	}
	returning, err := stringListValue(rawReturning)
	if err != nil {
		return nil, fmt.Errorf("%s returning %w", action, err)
	}
	return normalizeDBIdentifierList(action, "returning", returning)
}

func normalizeDBTimeoutMS(values map[string]any, action string) (int, error) {
	if raw, ok := values["timeout_ms"]; ok && raw != nil {
		timeoutMS, err := intParam(raw)
		if err != nil {
			return 0, fmt.Errorf("%s timeout_ms %w", action, err)
		}
		if timeoutMS < 1 {
			return 0, fmt.Errorf("%s timeout_ms must be at least 1", action)
		}
		return timeoutMS, nil
	}
	if raw, ok := values["timeout"]; ok && raw != nil {
		timeoutMS, err := intParam(raw)
		if err != nil {
			return 0, fmt.Errorf("%s timeout %w", action, err)
		}
		if timeoutMS < 1 {
			return 0, fmt.Errorf("%s timeout must be at least 1", action)
		}
		return timeoutMS, nil
	}
	if raw, ok := values["timeout_seconds"]; ok && raw != nil {
		timeoutSeconds, err := intParam(raw)
		if err != nil {
			return 0, fmt.Errorf("%s timeout_seconds %w", action, err)
		}
		if timeoutSeconds < 1 {
			return 0, fmt.Errorf("%s timeout_seconds must be at least 1", action)
		}
		return timeoutSeconds * 1000, nil
	}
	return 0, nil
}

func resolvedDBTimeoutForStep(ctx *FlowContext, step FlowStep, action string) (int, error) {
	values, err := resolvedDBValues(ctx, step, "", "timeout", "timeout_ms", "timeout_seconds")
	if err != nil {
		return 0, err
	}
	return normalizeDBTimeoutMS(values, action)
}

func normalizeDBIdentifierList(action string, name string, items []string) ([]string, error) {
	normalized := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, fmt.Errorf("%s %s cannot contain blank names", action, name)
		}
		if _, exists := seen[item]; exists {
			return nil, fmt.Errorf("%s %s cannot contain duplicate name %q", action, name, item)
		}
		seen[item] = struct{}{}
		normalized = append(normalized, item)
	}
	return normalized, nil
}

func objectListValue(value any, name string) ([]map[string]any, error) {
	items, err := toList(value)
	if err != nil {
		return nil, fmt.Errorf("%s must be a list of objects", name)
	}
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		row, err := objectMapValue(item, name)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func normalizeDBInsertManyColumns(action string, columns []string, rows []map[string]any) ([]string, error) {
	if len(columns) == 0 {
		columnSet := map[string]struct{}{}
		for _, row := range rows {
			for key := range row {
				columnSet[key] = struct{}{}
			}
		}
		for column := range columnSet {
			columns = append(columns, column)
		}
		sort.Strings(columns)
	}
	normalized, err := normalizeDBIdentifierList(action, "columns", columns)
	if err != nil {
		return nil, err
	}
	for rowIndex, row := range rows {
		for _, column := range normalized {
			if _, ok := row[column]; !ok {
				return nil, fmt.Errorf("%s row %d is missing column %q", action, rowIndex+1, column)
			}
		}
	}
	return normalized, nil
}

func normalizeDBStatementArgs(value any) ([]any, error) {
	if value == nil {
		return nil, nil
	}
	switch typed := value.(type) {
	case []any:
		args := make([]any, 0, len(typed))
		for _, item := range typed {
			normalized, err := normalizeDBStatementArg(item)
			if err != nil {
				return nil, err
			}
			args = append(args, normalized)
		}
		return args, nil
	case []string:
		args := make([]any, 0, len(typed))
		for _, item := range typed {
			args = append(args, item)
		}
		return args, nil
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		args := make([]any, 0, len(keys))
		for _, key := range keys {
			normalized, err := normalizeDBStatementArg(typed[key])
			if err != nil {
				return nil, err
			}
			args = append(args, sql.Named(key, normalized))
		}
		return args, nil
	case map[string]string:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		args := make([]any, 0, len(keys))
		for _, key := range keys {
			args = append(args, sql.Named(key, typed[key]))
		}
		return args, nil
	default:
		normalized, err := normalizeDBStatementArg(typed)
		if err != nil {
			return nil, err
		}
		return []any{normalized}, nil
	}
}

func normalizeDBStatementArg(value any) (any, error) {
	switch typed := value.(type) {
	case sql.NamedArg:
		normalized, err := normalizeDBArgument(typed.Value)
		if err != nil {
			return nil, err
		}
		typed.Value = normalized
		return typed, nil
	default:
		return normalizeDBArgument(value)
	}
}

func executeDBInsertWithFlow(flowCtx *FlowContext, ctx context.Context, config dbInsertConfig) (map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}
	query, args, hasReturning, err := buildDBInsertStatement(config, connection.Dialect)
	if err != nil {
		return nil, err
	}
	execCtx, cancel, executor, err := resolveDBExecutor(flowCtx, ctx, connection, config.TimeoutMS)
	if err != nil {
		return nil, err
	}
	defer cancel()
	payload := map[string]any{
		"ok":         true,
		"connection": connection.Name,
		"driver":     string(connection.Dialect),
		"table":      config.Table,
		"columns":    append([]string(nil), config.Columns...),
	}
	if hasReturning {
		rows, err := executor.QueryContext(execCtx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("db_insert on connection %q failed: %w", connection.Name, err)
		}
		defer rows.Close()
		resultRows, err := scanDBRows(rows)
		if err != nil {
			return nil, fmt.Errorf("db_insert on connection %q failed: %w", connection.Name, err)
		}
		if len(resultRows) > 0 {
			payload["returned"] = resultRows[0]
		}
		payload["returning"] = append([]string(nil), config.Returning...)
		payload["rows_affected"] = int64(len(resultRows))
		return payload, nil
	}
	result, err := executor.ExecContext(execCtx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db_insert on connection %q failed: %w", connection.Name, err)
	}
	if rowsAffected, err := result.RowsAffected(); err == nil {
		payload["rows_affected"] = rowsAffected
	}
	if lastInsertID, err := result.LastInsertId(); err == nil {
		payload["last_insert_id"] = lastInsertID
	}
	return payload, nil
}

func executeDBInsertManyWithFlow(flowCtx *FlowContext, ctx context.Context, config dbInsertManyConfig) (map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}
	query, args, hasReturning, err := buildDBInsertManyStatement(config, connection.Dialect)
	if err != nil {
		return nil, err
	}
	execCtx, cancel, executor, err := resolveDBExecutor(flowCtx, ctx, connection, config.TimeoutMS)
	if err != nil {
		return nil, err
	}
	defer cancel()
	payload := map[string]any{
		"ok":         true,
		"connection": connection.Name,
		"driver":     string(connection.Dialect),
		"table":      config.Table,
		"columns":    append([]string(nil), config.Columns...),
		"rows":       len(config.Rows),
	}
	if hasReturning {
		rows, err := executor.QueryContext(execCtx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("db_insert_many on connection %q failed: %w", connection.Name, err)
		}
		defer rows.Close()
		resultRows, err := scanDBRows(rows)
		if err != nil {
			return nil, fmt.Errorf("db_insert_many on connection %q failed: %w", connection.Name, err)
		}
		payload["returned_rows"] = resultRows
		payload["returning"] = append([]string(nil), config.Returning...)
		payload["rows_affected"] = int64(len(resultRows))
		return payload, nil
	}
	result, err := executor.ExecContext(execCtx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db_insert_many on connection %q failed: %w", connection.Name, err)
	}
	if rowsAffected, err := result.RowsAffected(); err == nil {
		payload["rows_affected"] = rowsAffected
	}
	return payload, nil
}

func executeDBUpsertWithFlow(flowCtx *FlowContext, ctx context.Context, config dbUpsertConfig) (map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}
	query, args, hasReturning, err := buildDBUpsertStatement(config, connection.Dialect)
	if err != nil {
		return nil, err
	}
	execCtx, cancel, executor, err := resolveDBExecutor(flowCtx, ctx, connection, config.TimeoutMS)
	if err != nil {
		return nil, err
	}
	defer cancel()
	payload := map[string]any{
		"ok":             true,
		"connection":     connection.Name,
		"driver":         string(connection.Dialect),
		"table":          config.Table,
		"columns":        append([]string(nil), config.Columns...),
		"key_columns":    append([]string(nil), config.KeyColumns...),
		"update_columns": append([]string(nil), config.UpdateColumns...),
	}
	if hasReturning {
		rows, err := executor.QueryContext(execCtx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("db_upsert on connection %q failed: %w", connection.Name, err)
		}
		defer rows.Close()
		resultRows, err := scanDBRows(rows)
		if err != nil {
			return nil, fmt.Errorf("db_upsert on connection %q failed: %w", connection.Name, err)
		}
		if len(resultRows) > 0 {
			payload["returned"] = resultRows[0]
		}
		payload["returning"] = append([]string(nil), config.Returning...)
		payload["rows_affected"] = int64(len(resultRows))
		return payload, nil
	}
	result, err := executor.ExecContext(execCtx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db_upsert on connection %q failed: %w", connection.Name, err)
	}
	if rowsAffected, err := result.RowsAffected(); err == nil {
		payload["rows_affected"] = rowsAffected
	}
	if lastInsertID, err := result.LastInsertId(); err == nil {
		payload["last_insert_id"] = lastInsertID
	}
	return payload, nil
}

func executeDBQueryWithFlow(flowCtx *FlowContext, ctx context.Context, config dbStatementConfig) ([]map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}
	execCtx, cancel, executor, err := resolveDBExecutor(flowCtx, ctx, connection, config.TimeoutMS)
	if err != nil {
		return nil, err
	}
	defer cancel()
	rows, err := executor.QueryContext(execCtx, config.SQL, config.Args...)
	if err != nil {
		return nil, fmt.Errorf("db_query on connection %q failed: %w", connection.Name, err)
	}
	defer rows.Close()
	return scanDBRows(rows)
}

func executeDBQueryOneWithFlow(flowCtx *FlowContext, ctx context.Context, config dbStatementConfig) (map[string]any, error) {
	rows, err := executeDBQueryWithFlow(flowCtx, ctx, config)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func executeDBExecuteWithFlow(flowCtx *FlowContext, ctx context.Context, config dbStatementConfig) (map[string]any, error) {
	connection, err := resolveDBConnectionConfig(config.Connection, config.Driver)
	if err != nil {
		return nil, err
	}
	execCtx, cancel, executor, err := resolveDBExecutor(flowCtx, ctx, connection, config.TimeoutMS)
	if err != nil {
		return nil, err
	}
	defer cancel()
	result, err := executor.ExecContext(execCtx, config.SQL, config.Args...)
	if err != nil {
		return nil, fmt.Errorf("db_execute on connection %q failed: %w", connection.Name, err)
	}
	payload := map[string]any{
		"ok":         true,
		"connection": connection.Name,
		"driver":     string(connection.Dialect),
	}
	if rowsAffected, err := result.RowsAffected(); err == nil {
		payload["rows_affected"] = rowsAffected
	}
	if lastInsertID, err := result.LastInsertId(); err == nil {
		payload["last_insert_id"] = lastInsertID
	}
	return payload, nil
}

func resolveDBExecutor(flowCtx *FlowContext, fallbackCtx context.Context, connection dbConnectionConfig, timeoutMS int) (context.Context, context.CancelFunc, flowDBExecutor, error) {
	runCtx := fallbackCtx
	if flowCtx != nil && flowCtx.Context != nil {
		runCtx = flowCtx.Context
	}
	if runCtx == nil {
		runCtx = context.Background()
	}

	if timeoutMS == 0 {
		settings, err := resolveDBRuntimeSettings(connection.Name)
		if err != nil {
			return nil, nil, nil, err
		}
		if settings.QueryTimeout > 0 {
			var cancel context.CancelFunc
			runCtx, cancel = context.WithTimeout(runCtx, settings.QueryTimeout)
			if flowCtx != nil && flowCtx.DBTransaction != nil {
				executor, err := flowCtx.DBTransaction.executor(runCtx, connection)
				return runCtx, cancel, executor, err
			}
			db, err := getFlowDatabase(connection)
			return runCtx, cancel, db, err
		}
	}

	cancel := func() {}
	if timeoutMS > 0 {
		runCtx, cancel = context.WithTimeout(runCtx, time.Duration(timeoutMS)*time.Millisecond)
	}
	if flowCtx != nil && flowCtx.DBTransaction != nil {
		executor, err := flowCtx.DBTransaction.executor(runCtx, connection)
		return runCtx, cancel, executor, err
	}
	db, err := getFlowDatabase(connection)
	return runCtx, cancel, db, err
}

func (scope *flowDBTransactionScope) executor(ctx context.Context, connection dbConnectionConfig) (flowDBExecutor, error) {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	if scope.sessions == nil {
		scope.sessions = map[string]*flowDBTransactionSession{}
	}
	key := connection.DriverName + "\n" + connection.DSN
	if session, ok := scope.sessions[key]; ok {
		return session.tx, nil
	}
	db, err := getFlowDatabase(connection)
	if err != nil {
		return nil, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction on connection %q: %w", connection.Name, err)
	}
	scope.sessions[key] = &flowDBTransactionSession{connection: connection, tx: tx}
	return tx, nil
}

func (scope *flowDBTransactionScope) Commit() error {
	return scope.finish(true)
}

func (scope *flowDBTransactionScope) Rollback() error {
	return scope.finish(false)
}

func (scope *flowDBTransactionScope) finish(commit bool) error {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	if len(scope.sessions) == 0 {
		return nil
	}
	keys := make([]string, 0, len(scope.sessions))
	for key := range scope.sessions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var firstErr error
	for _, key := range keys {
		session := scope.sessions[key]
		var err error
		if commit {
			err = session.tx.Commit()
		} else {
			err = session.tx.Rollback()
		}
		if err != nil && firstErr == nil {
			verb := "commit"
			if !commit {
				verb = "rollback"
			}
			firstErr = fmt.Errorf("%s transaction on connection %q: %w", verb, session.connection.Name, err)
		}
	}
	scope.sessions = map[string]*flowDBTransactionSession{}
	return firstErr
}

func (scope *flowDBTransactionScope) connectionNames() []string {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	names := make([]string, 0, len(scope.sessions))
	seen := map[string]struct{}{}
	for _, session := range scope.sessions {
		if _, exists := seen[session.connection.Name]; exists {
			continue
		}
		seen[session.connection.Name] = struct{}{}
		names = append(names, session.connection.Name)
	}
	sort.Strings(names)
	return names
}

func resolveDBRuntimeSettings(connection string) (dbRuntimeSettings, error) {
	settings := dbRuntimeSettings{
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxLifetime: 3 * time.Minute,
	}
	if value := lookupDBConfigValue(connection, "MAX_OPEN_CONNS"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return dbRuntimeSettings{}, fmt.Errorf("database connection %q MAX_OPEN_CONNS must be a positive integer", connection)
		}
		settings.MaxOpenConns = n
	}
	if value := lookupDBConfigValue(connection, "MAX_IDLE_CONNS"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return dbRuntimeSettings{}, fmt.Errorf("database connection %q MAX_IDLE_CONNS must be at least 0", connection)
		}
		settings.MaxIdleConns = n
	}
	if value := lookupDBConfigValue(connection, "CONN_MAX_LIFETIME_SECONDS"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return dbRuntimeSettings{}, fmt.Errorf("database connection %q CONN_MAX_LIFETIME_SECONDS must be a positive integer", connection)
		}
		settings.ConnMaxLifetime = time.Duration(n) * time.Second
	}
	if value := lookupDBConfigValue(connection, "QUERY_TIMEOUT_MS"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return dbRuntimeSettings{}, fmt.Errorf("database connection %q QUERY_TIMEOUT_MS must be a positive integer", connection)
		}
		settings.QueryTimeout = time.Duration(n) * time.Millisecond
	} else if value := lookupDBConfigValue(connection, "QUERY_TIMEOUT_SECONDS"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return dbRuntimeSettings{}, fmt.Errorf("database connection %q QUERY_TIMEOUT_SECONDS must be a positive integer", connection)
		}
		settings.QueryTimeout = time.Duration(n) * time.Second
	}
	return settings, nil
}

func buildDBInsertStatement(config dbInsertConfig, dialect dbDialect) (string, []any, bool, error) {
	tableName, quotedColumns, args, err := buildDBWriteParts(config.Table, config.Columns, func(index int, column string) (any, error) {
		return normalizeDBArgument(config.Row[column])
	}, dialect)
	if err != nil {
		return "", nil, false, err
	}
	placeholders := make([]string, 0, len(config.Columns))
	for index := range config.Columns {
		placeholders = append(placeholders, dbPlaceholder(dialect, index+1))
	}
	returningClause, outputClause, hasReturning, err := buildDBReturningFragments(dialect, config.Returning)
	if err != nil {
		return "", nil, false, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s)", tableName, strings.Join(quotedColumns, ", "))
	if outputClause != "" {
		query += " " + outputClause
	}
	query += fmt.Sprintf(" VALUES (%s)", strings.Join(placeholders, ", "))
	if returningClause != "" {
		query += " " + returningClause
	}
	return query, args, hasReturning, nil
}

func buildDBInsertManyStatement(config dbInsertManyConfig, dialect dbDialect) (string, []any, bool, error) {
	tableName, quotedColumns, _, err := buildDBWriteParts(config.Table, config.Columns, nil, dialect)
	if err != nil {
		return "", nil, false, err
	}
	returningClause, outputClause, hasReturning, err := buildDBReturningFragments(dialect, config.Returning)
	if err != nil {
		return "", nil, false, err
	}
	args := make([]any, 0, len(config.Rows)*len(config.Columns))
	switch dialect {
	case dbDialectOracle:
		if hasReturning {
			return "", nil, false, fmt.Errorf("db_insert_many returning is not supported for oracle")
		}
		intoClauses := make([]string, 0, len(config.Rows))
		argIndex := 1
		for _, row := range config.Rows {
			placeholders := make([]string, 0, len(config.Columns))
			for _, column := range config.Columns {
				value, err := normalizeDBArgument(row[column])
				if err != nil {
					return "", nil, false, fmt.Errorf("db_insert_many column %q %w", column, err)
				}
				args = append(args, value)
				placeholders = append(placeholders, dbPlaceholder(dialect, argIndex))
				argIndex++
			}
			intoClauses = append(intoClauses, fmt.Sprintf("INTO %s (%s) VALUES (%s)", tableName, strings.Join(quotedColumns, ", "), strings.Join(placeholders, ", ")))
		}
		query := "INSERT ALL " + strings.Join(intoClauses, " ") + " SELECT 1 FROM DUAL"
		return query, args, false, nil
	default:
		valueSets := make([]string, 0, len(config.Rows))
		argIndex := 1
		for _, row := range config.Rows {
			placeholders := make([]string, 0, len(config.Columns))
			for _, column := range config.Columns {
				value, err := normalizeDBArgument(row[column])
				if err != nil {
					return "", nil, false, fmt.Errorf("db_insert_many column %q %w", column, err)
				}
				args = append(args, value)
				placeholders = append(placeholders, dbPlaceholder(dialect, argIndex))
				argIndex++
			}
			valueSets = append(valueSets, "("+strings.Join(placeholders, ", ")+")")
		}
		query := fmt.Sprintf("INSERT INTO %s (%s)", tableName, strings.Join(quotedColumns, ", "))
		if outputClause != "" {
			query += " " + outputClause
		}
		query += " VALUES " + strings.Join(valueSets, ", ")
		if returningClause != "" {
			query += " " + returningClause
		}
		return query, args, hasReturning, nil
	}
}

func buildDBUpsertStatement(config dbUpsertConfig, dialect dbDialect) (string, []any, bool, error) {
	tableName, quotedColumns, args, err := buildDBWriteParts(config.Table, config.Columns, func(index int, column string) (any, error) {
		return normalizeDBArgument(config.Row[column])
	}, dialect)
	if err != nil {
		return "", nil, false, err
	}
	insertPlaceholders := make([]string, 0, len(config.Columns))
	for index := range config.Columns {
		insertPlaceholders = append(insertPlaceholders, dbPlaceholder(dialect, index+1))
	}
	returningClause, outputClause, hasReturning, err := buildDBReturningFragments(dialect, config.Returning)
	if err != nil {
		return "", nil, false, err
	}
	keyQuoted, err := quoteDBIdentifiers(config.KeyColumns, dialect)
	if err != nil {
		return "", nil, false, err
	}
	switch dialect {
	case dbDialectMySQL:
		if hasReturning {
			return "", nil, false, fmt.Errorf("db_upsert returning is not supported for mysql")
		}
		assignments := make([]string, 0, len(config.UpdateColumns))
		for _, column := range config.UpdateColumns {
			quotedColumn, err := quoteDBIdentifier(column, dialect)
			if err != nil {
				return "", nil, false, err
			}
			assignments = append(assignments, fmt.Sprintf("%s = VALUES(%s)", quotedColumn, quotedColumn))
		}
		if len(assignments) == 0 {
			assignments = append(assignments, fmt.Sprintf("%s = %s", keyQuoted[0], keyQuoted[0]))
		}
		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
			tableName,
			strings.Join(quotedColumns, ", "),
			strings.Join(insertPlaceholders, ", "),
			strings.Join(assignments, ", "),
		)
		return query, args, false, nil
	case dbDialectPostgres:
		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s)",
			tableName,
			strings.Join(quotedColumns, ", "),
			strings.Join(insertPlaceholders, ", "),
			strings.Join(keyQuoted, ", "),
		)
		if config.DoNothing {
			query += " DO NOTHING"
		} else {
			assignments := make([]string, 0, len(config.UpdateColumns))
			for _, column := range config.UpdateColumns {
				quotedColumn, err := quoteDBIdentifier(column, dialect)
				if err != nil {
					return "", nil, false, err
				}
				assignments = append(assignments, fmt.Sprintf("%s = EXCLUDED.%s", quotedColumn, quotedColumn))
			}
			query += " DO UPDATE SET " + strings.Join(assignments, ", ")
		}
		if returningClause != "" {
			query += " " + returningClause
		}
		return query, args, hasReturning, nil
	case dbDialectSQLServer, dbDialectOracle:
		sourceSelect, sourceRefs, err := buildDBMergeSourceSelect(config.Columns, dialect)
		if err != nil {
			return "", nil, false, err
		}
		onParts := make([]string, 0, len(config.KeyColumns))
		for _, keyColumn := range config.KeyColumns {
			quotedColumn, err := quoteDBIdentifier(keyColumn, dialect)
			if err != nil {
				return "", nil, false, err
			}
			onParts = append(onParts, fmt.Sprintf("target.%s = source.%s", quotedColumn, quotedColumn))
		}
		query := fmt.Sprintf("MERGE INTO %s target USING (%s) source ON (%s)", tableName, sourceSelect, strings.Join(onParts, " AND "))
		if !config.DoNothing {
			assignments := make([]string, 0, len(config.UpdateColumns))
			for _, column := range config.UpdateColumns {
				quotedColumn, err := quoteDBIdentifier(column, dialect)
				if err != nil {
					return "", nil, false, err
				}
				assignments = append(assignments, fmt.Sprintf("target.%s = source.%s", quotedColumn, quotedColumn))
			}
			if len(assignments) > 0 {
				query += " WHEN MATCHED THEN UPDATE SET " + strings.Join(assignments, ", ")
			}
		}
		query += fmt.Sprintf(
			" WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s)",
			strings.Join(quotedColumns, ", "),
			strings.Join(sourceRefs, ", "),
		)
		if outputClause != "" {
			query += " " + outputClause
		}
		if dialect == dbDialectSQLServer {
			query += ";"
		}
		if returningClause != "" {
			query += " " + returningClause
		}
		return query, args, hasReturning, nil
	default:
		return "", nil, false, fmt.Errorf("db_upsert does not support dialect %q", dialect)
	}
}

func buildDBWriteParts(table string, columns []string, valueAt func(index int, column string) (any, error), dialect dbDialect) (string, []string, []any, error) {
	tableName, err := quoteDBIdentifier(table, dialect)
	if err != nil {
		return "", nil, nil, fmt.Errorf("db table %w", err)
	}
	quotedColumns, err := quoteDBIdentifiers(columns, dialect)
	if err != nil {
		return "", nil, nil, err
	}
	args := make([]any, 0, len(columns))
	if valueAt != nil {
		for index, column := range columns {
			value, err := valueAt(index, column)
			if err != nil {
				return "", nil, nil, fmt.Errorf("db column %q %w", column, err)
			}
			args = append(args, value)
		}
	}
	return tableName, quotedColumns, args, nil
}

func buildDBReturningFragments(dialect dbDialect, columns []string) (string, string, bool, error) {
	if len(columns) == 0 {
		return "", "", false, nil
	}
	quotedColumns, err := quoteDBIdentifiers(columns, dialect)
	if err != nil {
		return "", "", false, err
	}
	switch dialect {
	case dbDialectPostgres:
		return "RETURNING " + strings.Join(quotedColumns, ", "), "", true, nil
	case dbDialectSQLServer:
		prefixed := make([]string, 0, len(quotedColumns))
		for _, column := range quotedColumns {
			prefixed = append(prefixed, "inserted."+column)
		}
		return "", "OUTPUT " + strings.Join(prefixed, ", "), true, nil
	case dbDialectMySQL:
		return "", "", false, fmt.Errorf("returning is not supported for mysql")
	case dbDialectOracle:
		return "", "", false, fmt.Errorf("returning is not supported for oracle")
	default:
		return "", "", false, fmt.Errorf("returning is not supported for dialect %q", dialect)
	}
}

func buildDBMergeSourceSelect(columns []string, dialect dbDialect) (string, []string, error) {
	items := make([]string, 0, len(columns))
	sourceRefs := make([]string, 0, len(columns))
	for index, column := range columns {
		quotedColumn, err := quoteDBIdentifier(column, dialect)
		if err != nil {
			return "", nil, err
		}
		placeholder := dbPlaceholder(dialect, index+1)
		if dialect == dbDialectOracle {
			items = append(items, fmt.Sprintf("%s %s", placeholder, quotedColumn))
		} else {
			items = append(items, fmt.Sprintf("%s AS %s", placeholder, quotedColumn))
		}
		sourceRefs = append(sourceRefs, "source."+quotedColumn)
	}
	if dialect == dbDialectOracle {
		return "SELECT " + strings.Join(items, ", ") + " FROM dual", sourceRefs, nil
	}
	return "SELECT " + strings.Join(items, ", "), sourceRefs, nil
}

func quoteDBIdentifiers(columns []string, dialect dbDialect) ([]string, error) {
	quoted := make([]string, 0, len(columns))
	for _, column := range columns {
		quotedColumn, err := quoteDBIdentifier(column, dialect)
		if err != nil {
			return nil, err
		}
		quoted = append(quoted, quotedColumn)
	}
	return quoted, nil
}

func scanDBRows(rows flowRows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		scanTargets := make([]any, len(columns))
		for index := range values {
			scanTargets[index] = &values[index]
		}
		if err := rows.Scan(scanTargets...); err != nil {
			return nil, err
		}
		row := map[string]any{}
		for index, column := range columns {
			row[column] = normalizeScannedDBValue(values[index])
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func normalizeScannedDBValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	default:
		return typed
	}
}

func dbShouldQuoteSimpleIdentifier(name string) bool {
	_, reserved := dbReservedIdentifierSet[strings.ToLower(strings.TrimSpace(name))]
	return reserved
}
