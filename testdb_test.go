// Copyright (c) 2024  The Go-Enjin Authors
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

package testdb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/tdata"
)

func TestTestDB(t *testing.T) {
	Convey("Constructor", t, func() {

		tempdir, err := tdata.NewTempData("", "tdata-testdb.*.d")
		So(err, ShouldBeNil)
		So(tempdir, ShouldNotEqual, "")
		defer tempdir.Destroy()

		Convey("in-memory", func() {
			tdb, err := NewTestDB()
			So(err, ShouldBeNil)
			So(tdb, ShouldNotBeNil)
			defer tdb.Close()
			So(tdb.DBH(), ShouldNotBeNil)
			So(tdb.SqliteDB(), ShouldEqual, "")
		})

		Convey("tempdir.db", func() {
			dbpath := tempdir.Join("test.db")
			tdb, err := NewTestDBWith(dbpath)
			So(err, ShouldBeNil)
			So(tdb, ShouldNotBeNil)
			defer tdb.Close()
			So(tdb.DBH(), ShouldNotBeNil)
			So(tdb.SqliteDB(), ShouldEqual, dbpath)
		})

		Convey("testing-testdb?mode=memory", func() {
			tdb, err := NewTestDBWith("file:testing-testdb?mode=memory")
			So(err, ShouldBeNil)
			So(tdb, ShouldNotBeNil)
			defer tdb.Close()
			So(tdb.DBH(), ShouldNotBeNil)
			So(tdb.SqliteDB(), ShouldEqual, "")
		})

	})

	Convey("Tables", t, func() {

		tempdir, err := tdata.NewTempData("", "tdata-testdb.*.d")
		So(err, ShouldBeNil)
		So(tempdir, ShouldNotEqual, "")
		defer tempdir.Destroy()

		tdb, err := NewTestDBWith("")
		So(err, ShouldBeNil)
		So(tdb, ShouldNotBeNil)
		defer tdb.Close()
		So(tdb.DBH(), ShouldNotBeNil)
		So(tdb.SqliteDB(), ShouldEqual, "")

		So(tdb.Tables(), ShouldEqual, []string(nil))
		So(tdb.HasTable("nope"), ShouldBeFalse)
		So(tdb.TableSchema("nope"), ShouldEqual, "")

		create := `CREATE TABLE "check" ( id INTEGER PRIMARY KEY NOT NULL, value TEXT )`
		r, err := tdb.DBH().Exec(create)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		r, err = tdb.DBH().Exec(`INSERT INTO "check" ("value") VALUES ("yes")`)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)

		So(tdb.Tables(), ShouldEqual, []string{"check"})
		So(tdb.HasTable("check"), ShouldBeTrue)
		So(tdb.TableSchema("check"), ShouldEqual, create)

		Convey("Indexes", func() {

			createIndex := `CREATE INDEX "check_value" ON "check" ( "value" )`
			r, err := tdb.DBH().Exec(createIndex)
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)

			So(tdb.Indexes(), ShouldEqual, []string{"check_value"})
			So(tdb.HasIndex("check_value"), ShouldBeTrue)
			So(tdb.IndexSchema("check_value"), ShouldEqual, createIndex)

		})

		Convey("SelectOne", func() {
			value, err := tdb.SelectOne("nope", `SELECT 1 FROM "check"`)
			So(err, ShouldNotBeNil)
			So(value, ShouldBeNil)
			value, err = tdb.SelectOne("value", `SELECT * FROM "check"`)
			So(err, ShouldBeNil)
			So(value, ShouldNotBeNil)
		})

		Convey("SelectList", func() {
			value, err := tdb.SelectList("nope", `SELECT 1 FROM "check"`)
			So(err, ShouldNotBeNil)
			So(value, ShouldBeNil)
			value, err = tdb.SelectList("value", `SELECT * FROM "check"`)
			So(err, ShouldBeNil)
			So(value, ShouldNotBeNil)
		})

	})
}
