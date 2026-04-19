package tsplay_core

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type fakeFlowDatabase struct {
	query             string
	args              []any
	queryRows         []map[string]any
	queryColumns      []string
	queryErr          error
	rowsAffected      int64
	lastInsertID      int64
	lastInsertIDError error
	execErr           error
	beginTxCount      int
	tx                *fakeFlowTransaction
}

func (db *fakeFlowDatabase) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	db.query = query
	db.args = append([]any(nil), args...)
	if db.execErr != nil {
		return nil, db.execErr
	}
	return fakeFlowResult{
		rowsAffected:      db.rowsAffected,
		lastInsertID:      db.lastInsertID,
		lastInsertIDError: db.lastInsertIDError,
	}, nil
}

func (db *fakeFlowDatabase) QueryContext(_ context.Context, query string, args ...any) (flowRows, error) {
	db.query = query
	db.args = append([]any(nil), args...)
	if db.queryErr != nil {
		return nil, db.queryErr
	}
	columns := append([]string(nil), db.queryColumns...)
	if len(columns) == 0 && len(db.queryRows) > 0 {
		for key := range db.queryRows[0] {
			columns = append(columns, key)
		}
		sort.Strings(columns)
	}
	return &fakeFlowRows{
		columns: columns,
		rows:    append([]map[string]any(nil), db.queryRows...),
	}, nil
}

func (db *fakeFlowDatabase) BeginTx(_ context.Context, _ *sql.TxOptions) (flowDBTransaction, error) {
	db.beginTxCount++
	if db.tx == nil {
		db.tx = &fakeFlowTransaction{db: db}
	}
	return db.tx, nil
}

func (db *fakeFlowDatabase) SetMaxOpenConns(_ int)              {}
func (db *fakeFlowDatabase) SetMaxIdleConns(_ int)              {}
func (db *fakeFlowDatabase) SetConnMaxLifetime(_ time.Duration) {}

type fakeFlowTransaction struct {
	db         *fakeFlowDatabase
	committed  bool
	rolledBack bool
}

func (tx *fakeFlowTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.db.ExecContext(ctx, query, args...)
}

func (tx *fakeFlowTransaction) QueryContext(ctx context.Context, query string, args ...any) (flowRows, error) {
	return tx.db.QueryContext(ctx, query, args...)
}

func (tx *fakeFlowTransaction) Commit() error {
	tx.committed = true
	return nil
}

func (tx *fakeFlowTransaction) Rollback() error {
	tx.rolledBack = true
	return nil
}

type fakeFlowRows struct {
	columns []string
	rows    []map[string]any
	index   int
}

func (rows *fakeFlowRows) Columns() ([]string, error) {
	return append([]string(nil), rows.columns...), nil
}

func (rows *fakeFlowRows) Next() bool {
	if rows.index >= len(rows.rows) {
		return false
	}
	rows.index++
	return true
}

func (rows *fakeFlowRows) Scan(dest ...any) error {
	if rows.index == 0 || rows.index > len(rows.rows) {
		return fmt.Errorf("scan out of bounds")
	}
	row := rows.rows[rows.index-1]
	for index, column := range rows.columns {
		ptr, ok := dest[index].(*any)
		if !ok {
			return fmt.Errorf("unexpected dest type %T", dest[index])
		}
		*ptr = row[column]
	}
	return nil
}

func (rows *fakeFlowRows) Err() error   { return nil }
func (rows *fakeFlowRows) Close() error { return nil }

type fakeFlowResult struct {
	rowsAffected      int64
	lastInsertID      int64
	lastInsertIDError error
}

func (result fakeFlowResult) LastInsertId() (int64, error) {
	if result.lastInsertIDError != nil {
		return 0, result.lastInsertIDError
	}
	return result.lastInsertID, nil
}

func (result fakeFlowResult) RowsAffected() (int64, error) {
	return result.rowsAffected, nil
}

func withFakeOpenFlowDatabase(t *testing.T, fake flowDatabase) {
	t.Helper()

	previousOpener := openFlowDatabase
	flowDatabaseCache = sync.Map{}
	openFlowDatabase = func(_ string, _ string) (flowDatabase, error) {
		return fake, nil
	}
	t.Cleanup(func() {
		openFlowDatabase = previousOpener
		flowDatabaseCache = sync.Map{}
	})
}

func TestValidateFlowSecurityRejectsDBInsertByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_policy",
		Steps: []FlowStep{
			{
				Action: "db_insert",
				With: map[string]any{
					"table":  "crawl_results",
					"driver": "pgsql",
					"row": map[string]any{
						"keyword": "山东大学",
					},
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected database security policy error")
	}
	if !strings.Contains(err.Error(), "allow_database") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowDBInsertPostgresAction(t *testing.T) {
	fakeDB := &fakeFlowDatabase{
		rowsAffected:      1,
		lastInsertIDError: fmt.Errorf("not supported"),
	}
	withFakeOpenFlowDatabase(t, fakeDB)

	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
	t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_insert_pgsql",
		Steps: []FlowStep{
			{
				Action:     "db_insert",
				Connection: "reporting",
				SaveAs:     "insert_result",
				With: map[string]any{
					"table": "public.crawl_results",
					"columns": []any{
						"keyword",
						"rank",
						"meta",
					},
					"row": map[string]any{
						"keyword": "山东大学",
						"rank":    1,
						"meta": map[string]any{
							"source": "baidu",
						},
					},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowDatabase: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	if fakeDB.query != `INSERT INTO public.crawl_results (keyword, "rank", meta) VALUES ($1, $2, $3)` {
		t.Fatalf("query = %q", fakeDB.query)
	}
	if len(fakeDB.args) != 3 {
		t.Fatalf("args = %#v", fakeDB.args)
	}
	if fakeDB.args[0] != "山东大学" || fakeDB.args[1] != 1 {
		t.Fatalf("args = %#v", fakeDB.args)
	}
	if fakeDB.args[2] != `{"source":"baidu"}` {
		t.Fatalf("meta arg = %#v", fakeDB.args[2])
	}

	insertResult, ok := result.Vars["insert_result"].(map[string]any)
	if !ok {
		t.Fatalf("insert_result = %#v", result.Vars["insert_result"])
	}
	if insertResult["driver"] != "pgsql" {
		t.Fatalf("driver = %#v", insertResult["driver"])
	}
	if fmt.Sprint(insertResult["rows_affected"]) != "1" {
		t.Fatalf("rows_affected = %#v", insertResult["rows_affected"])
	}
	if _, exists := insertResult["last_insert_id"]; exists {
		t.Fatalf("last_insert_id should be omitted for pgsql: %#v", insertResult["last_insert_id"])
	}
}

func TestRunFlowDBInsertMySQLLegacyConnection(t *testing.T) {
	fakeDB := &fakeFlowDatabase{
		rowsAffected: 1,
		lastInsertID: 42,
	}
	withFakeOpenFlowDatabase(t, fakeDB)

	t.Setenv("TSPLAY_MYSQL_REPORTING_HOST", "127.0.0.1")
	t.Setenv("TSPLAY_MYSQL_REPORTING_PORT", "3306")
	t.Setenv("TSPLAY_MYSQL_REPORTING_USER", "collector")
	t.Setenv("TSPLAY_MYSQL_REPORTING_PASSWORD", "secret")
	t.Setenv("TSPLAY_MYSQL_REPORTING_DATABASE", "analytics")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_insert_mysql_legacy",
		Steps: []FlowStep{
			{
				Action:     "db_insert",
				Connection: "reporting",
				SaveAs:     "insert_result",
				With: map[string]any{
					"table":  "crawl_results",
					"driver": "mysql",
					"columns": []any{
						"keyword",
						"rank",
					},
					"row": map[string]any{
						"keyword": "山东大学",
						"rank":    1,
					},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowDatabase: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	if fakeDB.query != "INSERT INTO crawl_results (keyword, `rank`) VALUES (?, ?)" {
		t.Fatalf("query = %q", fakeDB.query)
	}
	insertResult, ok := result.Vars["insert_result"].(map[string]any)
	if !ok {
		t.Fatalf("insert_result = %#v", result.Vars["insert_result"])
	}
	if insertResult["driver"] != "mysql" {
		t.Fatalf("driver = %#v", insertResult["driver"])
	}
	if fmt.Sprint(insertResult["last_insert_id"]) != "42" {
		t.Fatalf("last_insert_id = %#v", insertResult["last_insert_id"])
	}
}

func TestResolveDBConnectionConfigOracleRequiresTaggedBuild(t *testing.T) {
	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "gora")
	t.Setenv("TSPLAY_DB_REPORTING_DSN", "oracle://collector:secret@127.0.0.1:1521/orclpdb1")

	_, err := getFlowDatabase(dbConnectionConfig{
		Name:       "reporting",
		Dialect:    dbDialectOracle,
		DriverName: "oracle",
		DSN:        "oracle://collector:secret@127.0.0.1:1521/orclpdb1",
	})
	if err == nil {
		t.Fatalf("expected oracle driver registration error")
	}
	if !strings.Contains(err.Error(), "tsplay_oracle") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildDBInsertQuerySQLServerPlaceholders(t *testing.T) {
	query, args, err := buildDBInsertQuery(dbInsertConfig{
		Table:   "dbo.crawl_results",
		Columns: []string{"keyword", "rank"},
		Row: map[string]any{
			"keyword": "山东大学",
			"rank":    1,
		},
	}, dbDialectSQLServer)
	if err != nil {
		t.Fatalf("build query: %v", err)
	}
	if query != "INSERT INTO dbo.crawl_results (keyword, [rank]) VALUES (@p1, @p2)" {
		t.Fatalf("query = %q", query)
	}
	if len(args) != 2 || args[0] != "山东大学" || args[1] != 1 {
		t.Fatalf("args = %#v", args)
	}
}

func TestResolveStructuredDBConnectionConfigOracle(t *testing.T) {
	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "gora")
	t.Setenv("TSPLAY_DB_REPORTING_HOST", "127.0.0.1")
	t.Setenv("TSPLAY_DB_REPORTING_PORT", "1521")
	t.Setenv("TSPLAY_DB_REPORTING_SERVICE", "orclpdb1")
	t.Setenv("TSPLAY_DB_REPORTING_USER", "collector")
	t.Setenv("TSPLAY_DB_REPORTING_PASSWORD", "secret")

	config, err := resolveDBConnectionConfig("reporting", "")
	if err != nil {
		t.Fatalf("resolve connection: %v", err)
	}
	if config.Dialect != dbDialectOracle {
		t.Fatalf("dialect = %q", config.Dialect)
	}
	if config.DriverName != "oracle" {
		t.Fatalf("driver name = %q", config.DriverName)
	}
	if config.DSN != "oracle://collector:secret@127.0.0.1:1521/orclpdb1" {
		t.Fatalf("dsn = %q", config.DSN)
	}
}

func TestResolveStructuredDBConnectionConfigOracleConnectString(t *testing.T) {
	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "go-ora")
	t.Setenv("TSPLAY_DB_REPORTING_CONNECT_STRING", `(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=127.0.0.1)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))`)
	t.Setenv("TSPLAY_DB_REPORTING_USER", "collector")
	t.Setenv("TSPLAY_DB_REPORTING_PASSWORD", "secret")

	config, err := resolveDBConnectionConfig("reporting", "")
	if err != nil {
		t.Fatalf("resolve connection: %v", err)
	}
	if config.DriverName != "oracle" {
		t.Fatalf("driver name = %q", config.DriverName)
	}
	if !strings.Contains(config.DSN, "oracle://collector:secret@:0/?") {
		t.Fatalf("dsn = %q", config.DSN)
	}
	if !strings.Contains(config.DSN, "connStr=%28DESCRIPTION%3D") {
		t.Fatalf("dsn = %q", config.DSN)
	}
}

func TestNormalizeDBDriverAliases(t *testing.T) {
	cases := []struct {
		input       string
		wantDialect dbDialect
		wantDriver  string
	}{
		{input: "pgsql", wantDialect: dbDialectPostgres, wantDriver: "postgres"},
		{input: "postgres", wantDialect: dbDialectPostgres, wantDriver: "postgres"},
		{input: "sqlserver", wantDialect: dbDialectSQLServer, wantDriver: "sqlserver"},
		{input: "oracle", wantDialect: dbDialectOracle, wantDriver: "oracle"},
	}

	for _, tc := range cases {
		dialect, driverName, err := normalizeDBDriver(tc.input)
		if err != nil {
			t.Fatalf("normalize %q: %v", tc.input, err)
		}
		if dialect != tc.wantDialect || driverName != tc.wantDriver {
			t.Fatalf("normalize %q = (%q, %q), want (%q, %q)", tc.input, dialect, driverName, tc.wantDialect, tc.wantDriver)
		}
	}
}

func TestRunFlowDBInsertManyAction(t *testing.T) {
	fakeDB := &fakeFlowDatabase{rowsAffected: 2}
	withFakeOpenFlowDatabase(t, fakeDB)

	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
	t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_insert_many_pgsql",
		Steps: []FlowStep{
			{
				Action:     "db_insert_many",
				Connection: "reporting",
				SaveAs:     "insert_many_result",
				With: map[string]any{
					"table": "public.crawl_results",
					"rows": []any{
						map[string]any{"keyword": "山东大学", "rank": 1},
						map[string]any{"keyword": "北京大学", "rank": 2},
					},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowDatabase: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if fakeDB.query != `INSERT INTO public.crawl_results (keyword, "rank") VALUES ($1, $2), ($3, $4)` {
		t.Fatalf("query = %q", fakeDB.query)
	}
	if len(fakeDB.args) != 4 {
		t.Fatalf("args = %#v", fakeDB.args)
	}
	insertResult, ok := result.Vars["insert_many_result"].(map[string]any)
	if !ok {
		t.Fatalf("insert_many_result = %#v", result.Vars["insert_many_result"])
	}
	if fmt.Sprint(insertResult["rows_affected"]) != "2" {
		t.Fatalf("rows_affected = %#v", insertResult["rows_affected"])
	}
}

func TestBuildDBUpsertStatementPostgresReturning(t *testing.T) {
	query, args, hasReturning, err := buildDBUpsertStatement(dbUpsertConfig{
		Table:         "public.crawl_results",
		Columns:       []string{"keyword", "rank"},
		KeyColumns:    []string{"keyword"},
		UpdateColumns: []string{"rank"},
		Returning:     []string{"keyword", "rank"},
		Row: map[string]any{
			"keyword": "山东大学",
			"rank":    1,
		},
	}, dbDialectPostgres)
	if err != nil {
		t.Fatalf("build upsert: %v", err)
	}
	want := `INSERT INTO public.crawl_results (keyword, "rank") VALUES ($1, $2) ON CONFLICT (keyword) DO UPDATE SET "rank" = EXCLUDED."rank" RETURNING keyword, "rank"`
	if query != want {
		t.Fatalf("query = %q", query)
	}
	if !hasReturning {
		t.Fatalf("expected returning query")
	}
	if len(args) != 2 || args[0] != "山东大学" || args[1] != 1 {
		t.Fatalf("args = %#v", args)
	}
}

func TestRunFlowDBQueryAndExecuteActions(t *testing.T) {
	fakeDB := &fakeFlowDatabase{
		queryColumns: []string{"keyword", "rank"},
		queryRows: []map[string]any{
			{"keyword": "山东大学", "rank": 1},
			{"keyword": "北京大学", "rank": 2},
		},
		rowsAffected: 3,
	}
	withFakeOpenFlowDatabase(t, fakeDB)

	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
	t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_query_and_execute",
		Steps: []FlowStep{
			{
				Action:     "db_query",
				Connection: "reporting",
				SaveAs:     "rows",
				With: map[string]any{
					"sql":  "SELECT keyword, rank FROM public.crawl_results WHERE rank >= $1",
					"args": []any{1},
				},
			},
			{
				Action:     "db_query_one",
				Connection: "reporting",
				SaveAs:     "row",
				With: map[string]any{
					"sql":  "SELECT keyword, rank FROM public.crawl_results WHERE rank >= $1",
					"args": []any{1},
				},
			},
			{
				Action:     "db_execute",
				Connection: "reporting",
				SaveAs:     "exec_result",
				With: map[string]any{
					"sql":  "DELETE FROM public.crawl_results WHERE rank < $1",
					"args": []any{10},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowDatabase: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	rows, ok := result.Vars["rows"].([]map[string]any)
	if !ok || len(rows) != 2 {
		t.Fatalf("rows = %#v", result.Vars["rows"])
	}
	row, ok := result.Vars["row"].(map[string]any)
	if !ok || row["keyword"] != "山东大学" {
		t.Fatalf("row = %#v", result.Vars["row"])
	}
	execResult, ok := result.Vars["exec_result"].(map[string]any)
	if !ok || fmt.Sprint(execResult["rows_affected"]) != "3" {
		t.Fatalf("exec_result = %#v", result.Vars["exec_result"])
	}
}

func TestRunFlowDBTransactionCommitAndRollback(t *testing.T) {
	t.Run("commit", func(t *testing.T) {
		fakeDB := &fakeFlowDatabase{rowsAffected: 1}
		withFakeOpenFlowDatabase(t, fakeDB)
		t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
		t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

		L := lua.NewState()
		defer L.Close()

		flow := &Flow{
			SchemaVersion: "1",
			Name:          "db_transaction_commit",
			Steps: []FlowStep{
				{
					Action: "db_transaction",
					Steps: []FlowStep{
						{
							Action:     "db_insert",
							Connection: "reporting",
							With: map[string]any{
								"table": "public.crawl_results",
								"row": map[string]any{
									"keyword": "山东大学",
								},
							},
						},
					},
				},
			},
		}

		if _, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
			Security: &FlowSecurityPolicy{AllowDatabase: true},
		}); err != nil {
			t.Fatalf("run flow: %v", err)
		}
		if fakeDB.beginTxCount != 1 {
			t.Fatalf("beginTxCount = %d", fakeDB.beginTxCount)
		}
		if fakeDB.tx == nil || !fakeDB.tx.committed {
			t.Fatalf("expected committed transaction")
		}
	})

	t.Run("rollback", func(t *testing.T) {
		fakeDB := &fakeFlowDatabase{execErr: fmt.Errorf("boom")}
		withFakeOpenFlowDatabase(t, fakeDB)
		t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
		t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

		L := lua.NewState()
		defer L.Close()

		flow := &Flow{
			SchemaVersion: "1",
			Name:          "db_transaction_rollback",
			Steps: []FlowStep{
				{
					Action: "db_transaction",
					Steps: []FlowStep{
						{
							Action:     "db_insert",
							Connection: "reporting",
							With: map[string]any{
								"table": "public.crawl_results",
								"row": map[string]any{
									"keyword": "山东大学",
								},
							},
						},
					},
				},
			},
		}

		if _, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
			Security: &FlowSecurityPolicy{AllowDatabase: true},
		}); err == nil {
			t.Fatalf("expected transaction failure")
		}
		if fakeDB.tx == nil || !fakeDB.tx.rolledBack {
			t.Fatalf("expected rolled back transaction")
		}
	})
}

func TestRunFlowLuaDatabaseHelpersHonorAllowDatabase(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_db_policy",
		Steps: []FlowStep{
			{
				Action: "lua",
				Code: `return db_query({
  sql = "SELECT 1",
  connection = "reporting",
  driver = "pgsql"
})`,
			},
		},
	}

	_, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowLua: true},
	})
	if err == nil {
		t.Fatalf("expected allow_database runtime error")
	}
	if !strings.Contains(err.Error(), "allow_database") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowLuaDBTransactionCommitAndRollback(t *testing.T) {
	t.Run("commit", func(t *testing.T) {
		fakeDB := &fakeFlowDatabase{rowsAffected: 1}
		withFakeOpenFlowDatabase(t, fakeDB)
		t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
		t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

		L := lua.NewState()
		defer L.Close()

		flow := &Flow{
			SchemaVersion: "1",
			Name:          "lua_db_transaction_commit",
			Steps: []FlowStep{
				{
					Action: "lua",
					SaveAs: "tx_result",
					Code: `return db_transaction(function()
  return db_insert({
    table = "public.crawl_results",
    connection = "reporting",
    row = {
      keyword = "山东大学"
    }
  })
end, 5000)`,
				},
			},
		}

		result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
			Security: &FlowSecurityPolicy{
				AllowLua:      true,
				AllowDatabase: true,
			},
		})
		if err != nil {
			t.Fatalf("run flow: %v", err)
		}
		if fakeDB.beginTxCount != 1 {
			t.Fatalf("beginTxCount = %d", fakeDB.beginTxCount)
		}
		if fakeDB.tx == nil || !fakeDB.tx.committed {
			t.Fatalf("expected committed transaction")
		}
		txResult, ok := result.Vars["tx_result"].(map[string]any)
		if !ok || fmt.Sprint(txResult["rows_affected"]) != "1" {
			t.Fatalf("tx_result = %#v", result.Vars["tx_result"])
		}
	})

	t.Run("rollback", func(t *testing.T) {
		fakeDB := &fakeFlowDatabase{rowsAffected: 1}
		withFakeOpenFlowDatabase(t, fakeDB)
		t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "pgsql")
		t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

		L := lua.NewState()
		defer L.Close()

		flow := &Flow{
			SchemaVersion: "1",
			Name:          "lua_db_transaction_rollback",
			Steps: []FlowStep{
				{
					Action: "lua",
					Code: `return db_transaction(function()
  db_insert({
    table = "public.crawl_results",
    connection = "reporting",
    row = {
      keyword = "山东大学"
    }
  })
  error("boom")
end)`,
				},
			},
		}

		if _, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
			Security: &FlowSecurityPolicy{
				AllowLua:      true,
				AllowDatabase: true,
			},
		}); err == nil {
			t.Fatalf("expected transaction failure")
		}
		if fakeDB.beginTxCount != 1 {
			t.Fatalf("beginTxCount = %d", fakeDB.beginTxCount)
		}
		if fakeDB.tx == nil || !fakeDB.tx.rolledBack {
			t.Fatalf("expected rolled back transaction")
		}
	})
}

func TestResolveDBRuntimeSettings(t *testing.T) {
	t.Setenv("TSPLAY_DB_REPORTING_MAX_OPEN_CONNS", "8")
	t.Setenv("TSPLAY_DB_REPORTING_MAX_IDLE_CONNS", "4")
	t.Setenv("TSPLAY_DB_REPORTING_CONN_MAX_LIFETIME_SECONDS", "120")
	t.Setenv("TSPLAY_DB_REPORTING_QUERY_TIMEOUT_SECONDS", "9")

	settings, err := resolveDBRuntimeSettings("reporting")
	if err != nil {
		t.Fatalf("resolve settings: %v", err)
	}
	if settings.MaxOpenConns != 8 || settings.MaxIdleConns != 4 {
		t.Fatalf("settings = %#v", settings)
	}
	if settings.ConnMaxLifetime != 120*time.Second {
		t.Fatalf("ConnMaxLifetime = %v", settings.ConnMaxLifetime)
	}
	if settings.QueryTimeout != 9*time.Second {
		t.Fatalf("QueryTimeout = %v", settings.QueryTimeout)
	}
}
