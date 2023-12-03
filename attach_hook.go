package sqlu

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"modernc.org/sqlite"
)

var (
	ErrAttachDatabaseNameMissing = errors.New("name missing")
)

func init() {
	sqlite.RegisterConnectionHook(AttachConnectionHook)
}

func AttachConnectionHook(conn sqlite.ExecQuerierContext, dsn string) error {
	qPos := strings.IndexRune(dsn, '?')
	if qPos < 1 {
		return nil
	}

	queryStr := dsn[qPos+1:]
	query, err := url.ParseQuery(queryStr)
	if err != nil {
		return fmt.Errorf("parse query: %s: %w", queryStr, err)
	}

	for _, attachStr := range query["_attach"] {
		attachDSN, err := url.QueryUnescape(attachStr)
		if err != nil {
			return fmt.Errorf("unescape attach dsn: %s: %w", attachStr, err)
		}

		err = attachDatabase(conn, attachDSN)
		if err != nil {
			return fmt.Errorf("attach database: %w", err)
		}
	}

	return nil
}

func attachDatabase(
	conn driver.ExecerContext,
	dsn string,
) error {
	var query url.Values
	var err error

	qPos := strings.IndexRune(dsn, '?')
	if qPos < 1 {
		return fmt.Errorf("query not found in DSN: %s", dsn)
	}

	queryStr := dsn[qPos+1:]
	query, err = url.ParseQuery(queryStr)
	if err != nil {
		return fmt.Errorf("parse query: %s: %w", queryStr, err)
	}

	filename := dsn[:qPos]

	databaseName := query.Get("_name")
	if databaseName == "" {
		return ErrAttachDatabaseNameMissing
	}

	_, err = conn.ExecContext(
		context.Background(),
		fmt.Sprintf(`ATTACH DATABASE '%s' AS %s`, filename, databaseName),
		nil,
	)
	if err != nil {
		return fmt.Errorf("attach exec: %w", err)
	}

	for _, value := range query["_pragma"] {
		cmd := "pragma " + databaseName + "." + value

		_, err = conn.ExecContext(context.Background(), cmd, nil)
		if err != nil {
			return fmt.Errorf("pragma exec failed: %w", err)
		}
	}

	return nil
}
