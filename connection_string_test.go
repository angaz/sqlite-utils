package sqlu

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func TestConnectionString(t *testing.T) {
	connStr := ConnParams{
		Filename: "test.db",
		Pragma: []Pragma{
			PragmaBusyTimeout(3000),
			PragmaJournalModeWAL,
		},
		Attach: []AttachParams{
			{
				Filename: "test2.db",
				Database: "test2",
				Pragma: []Pragma{
					PragmaBusyTimeout(2000),
				},
			},
		},
	}

	dsn := connStr.ConnectionString()
	expectedDSN := ("test.db?_attach=test2.db%3F_name%3Dtest2%26" +
		"_pragma%3Dbusy_timeout%253D2000&_pragma=busy_timeout%3D3000&" +
		"_pragma=journal_mode%3Dwal")

	if dsn != expectedDSN {
		t.Errorf(
			"dsn string did not match: want: %s, got: %s",
			expectedDSN,
			dsn,
		)
	}

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}

	busyTimeout, err := queryString(conn, "PRAGMA busy_timeout")
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}

	if busyTimeout != "2000" {
		t.Errorf("busy timeout wrong: want: 2000, got: %s", busyTimeout)
	}

	busyTimeout, err = queryString(conn, "PRAGMA test2.busy_timeout")
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}

	if busyTimeout != "2000" {
		t.Errorf("busy timeout wrong: want: 2000, got: %s", busyTimeout)
	}

	journalMode, err := queryString(conn, "PRAGMA journal_mode")
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}

	if journalMode != "wal" {
		t.Errorf(`journal mode wrong: want: "wal", got: %s`, journalMode)
	}

	journalMode, err = queryString(conn, "PRAGMA test2.journal_mode")
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}

	if journalMode != "delete" {
		t.Errorf(`journal mode wrong: want: "delete", got: %s`, journalMode)
	}

	os.Remove("test.db")
	os.Remove("test2.db")
}

func queryString(conn *sql.DB, query string) (string, error) {
	row := conn.QueryRow(query)
	if row.Err() != nil {
		return "", fmt.Errorf("query exec: %w", row.Err())
	}

	var result string

	err := row.Scan(&result)
	if err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}

	return result, nil
}
