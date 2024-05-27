// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package testdb provides go db testing utilities
package testdb

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-corelibs/slices"
)

var (
	gSelectSqliteSchemaSQL   = `SELECT sql  FROM sqlite_schema WHERE type = ? AND name = ?`
	gSelectSqliteSchemaName  = `SELECT name FROM sqlite_schema WHERE type = ? AND name = ?                 ORDER BY name`
	gSelectSqliteSchemaNames = `SELECT name FROM sqlite_schema WHERE type = ? AND name NOT LIKE "sqlite_%" ORDER BY name`
)

var _ TestDB = (*testdb)(nil)

// TestDB is the interface for Sqlite3 based test databases. These databases
// are intended to be ephemeral and easily reset for the purposes of unit
// testing within other projects
type TestDB interface {
	// Open starts a new sqlite database, destroying any previous one with
	// TestDB.Close first
	Open() (err error)
	// Close stops the existing sqlite connection and if the database is not
	// ":memory:", removes the db file
	Close()
	// DBH returns the sql.DB instance opened during NewTestDBWith
	DBH() *sql.DB
	// SqliteDB returns the main database path from the "pragma_database_list"
	// table
	SqliteDB() (file string)
	// Tables returns the list of non-sqlite tables from the "sqlite_schema"
	// table
	Tables() (names []string)
	// HasTable returns true if the given table name is present in the
	// "sqlite_schema" table
	HasTable(name string) (present bool)
	// TableSchema returns the SQL definition from the "sqlite_schema" table
	// for the named table
	TableSchema(name string) (schema string)
	// Indexes returns the list of non-sqlite indexes from the "sqlite_schema"
	// table
	Indexes() (names []string)
	// HasIndex returns true if the given table name is present in the
	// "sqlite_schema" table
	HasIndex(name string) (present bool)
	// IndexSchema returns the SQL definition from the "sqlite_schema" table
	// for the named index
	IndexSchema(name string) (schema string)
	// Select is a very simple wrapper around a sql.DB Query call, gathering
	// all the results in a simple mapping of column names to interface{}
	// values
	Select(query string, argv ...interface{}) (results []map[string]interface{}, err error)
	// SelectOne is a wrapper around Select and returning just the first result's
	// specific column value
	SelectOne(column, query string, argv ...interface{}) (value interface{}, err error)
	// SelectList is a wrapper around Select and returning a list of just the
	// specific column values
	SelectList(column, query string, argv ...interface{}) (values []interface{}, err error)
}

type testdb struct {
	file string
	dbh  *sql.DB
}

// NewTestDB is a wrapper around ":memory:" call to NewTestDBWith
func NewTestDB() (tdb TestDB, err error) {
	return NewTestDBWith(":memory:")
}

// NewTestDBWith opens the given database file and returns a new TestDB instance,
// if the file argument is empty, a ":memory:" database is used
//
// Note that if the database file given is actually a file, the file will be
// deleted when Close is called
func NewTestDBWith(file string) (tdb TestDB, err error) {
	if file == "" {
		file = ":memory:"
	}
	file += "?cache=shared"
	t := &testdb{file: file}
	if err = t.Open(); err == nil {
		tdb = t
	}
	return
}

func (t *testdb) Open() (err error) {
	t.dbh, err = sql.Open("sqlite3", t.file)
	return
}

func (t *testdb) Close() {
	if t.dbh != nil {
		file := t.SqliteDB()
		_ = t.dbh.Close()
		file, _, _ = strings.Cut(file, "?")
		if file != "" && file != ":memory:" {
			_ = os.Remove(file)
		}
	}
}

func (t *testdb) DBH() *sql.DB {
	return t.dbh
}

func (t *testdb) SqliteDB() (file string) {
	query := `SELECT file FROM pragma_database_list WHERE name = "main"`
	if results, err := t.Select(query); err == nil {
		if len(results) == 1 {
			file, _ = results[0]["file"].(string)
		}
	}
	return
}

func (t *testdb) Tables() (names []string) {
	if values, err := t.SelectList("name", gSelectSqliteSchemaNames, "table"); err == nil {
		names = slices.ToStrings(values)
	}
	return
}

func (t *testdb) HasTable(name string) (present bool) {
	if value, err := t.SelectOne("name", gSelectSqliteSchemaName, "table", name); err == nil {
		present = value != nil
	}
	return
}

func (t *testdb) TableSchema(name string) (schema string) {
	if value, err := t.SelectOne("sql", gSelectSqliteSchemaSQL, "table", name); err == nil {
		schema, _ = value.(string)
	}
	return
}

func (t *testdb) Indexes() (names []string) {
	if values, err := t.SelectList("name", gSelectSqliteSchemaNames, "index"); err == nil {
		names = slices.ToStrings(values)
	}
	return
}

func (t *testdb) HasIndex(name string) (present bool) {
	if value, err := t.SelectOne("name", gSelectSqliteSchemaName, "index", name); err == nil {
		present = value != nil
	}
	return
}

func (t *testdb) IndexSchema(name string) (schema string) {
	if value, err := t.SelectOne("sql", gSelectSqliteSchemaSQL, "index", name); err == nil {
		schema, _ = value.(string)
	}
	return
}

func (t *testdb) SelectOne(column, query string, argv ...interface{}) (value interface{}, err error) {
	var results []map[string]interface{}
	if results, err = t.Select(query, argv...); err == nil {
		if len(results) > 0 {
			result := results[0]
			var ok bool
			if value, ok = result[column]; ok {
				return
			}
			err = fmt.Errorf("%w: %q", ErrResultColumnNotFound, column)
		}
	}
	return
}

func (t *testdb) SelectList(column, query string, argv ...interface{}) (values []interface{}, err error) {
	var results []map[string]interface{}
	if results, err = t.Select(query, argv...); err == nil {
		for _, result := range results {
			if value, ok := result[column]; !ok {
				err = fmt.Errorf("%w: %q", ErrResultColumnNotFound, column)
				return
			} else {
				values = append(values, value)
			}
		}
	}
	return
}

func (t *testdb) Select(query string, argv ...interface{}) (results []map[string]interface{}, err error) {
	var rows *sql.Rows
	if rows, err = t.dbh.Query(query, argv...); err == nil {
		for rows.Next() {
			var values []interface{}
			columns, _ := rows.Columns()
			for range columns {
				var v interface{}
				values = append(values, &v)
			}
			_ = rows.Scan(values...) // no custom value scanners used, can ignore err safely
			result := make(map[string]interface{})
			for idx, name := range columns {
				if v, ok := values[idx].(*interface{}); ok {
					result[name] = *v
				}
			}
			results = append(results, result)
		}
	}
	return
}
