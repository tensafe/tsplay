package tsplay_core

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type fakeFlowDatabase struct {
	query             string
	args              []any
	rowsAffected      int64
	lastInsertID      int64
	lastInsertIDError error
	execErr           error
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

func (db *fakeFlowDatabase) SetMaxOpenConns(_ int)              {}
func (db *fakeFlowDatabase) SetMaxIdleConns(_ int)              {}
func (db *fakeFlowDatabase) SetConnMaxLifetime(_ time.Duration) {}

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
					"driver": "postgres",
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

	t.Setenv("TSPLAY_DB_REPORTING_DRIVER", "postgres")
	t.Setenv("TSPLAY_DB_REPORTING_DSN", "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "db_insert_postgres",
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

	if fakeDB.query != `INSERT INTO public.crawl_results (keyword, rank, meta) VALUES ($1, $2, $3)` {
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
	if insertResult["driver"] != "postgres" {
		t.Fatalf("driver = %#v", insertResult["driver"])
	}
	if fmt.Sprint(insertResult["rows_affected"]) != "1" {
		t.Fatalf("rows_affected = %#v", insertResult["rows_affected"])
	}
	if _, exists := insertResult["last_insert_id"]; exists {
		t.Fatalf("last_insert_id should be omitted for postgres: %#v", insertResult["last_insert_id"])
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

	if fakeDB.query != "INSERT INTO crawl_results (keyword, rank) VALUES (?, ?)" {
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
	if query != "INSERT INTO dbo.crawl_results (keyword, rank) VALUES (@p1, @p2)" {
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
