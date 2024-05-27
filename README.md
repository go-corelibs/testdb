[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/testdb)
[![codecov](https://codecov.io/gh/go-corelibs/testdb/graph/badge.svg?token=mHjCyOnLMM)](https://codecov.io/gh/go-corelibs/testdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/testdb)](https://goreportcard.com/report/github.com/go-corelibs/testdb)

# testdb - go db testing utilities

A collection of utilities for simple unit testing with a sqlite db

# Installation

``` shell
> go get github.com/go-corelibs/testdb@latest
```

# Examples

## NewTestDB

``` go
func main() {
    tdb, err := testdb.NewTestDB()
    defer tdb.Close()               // shutdown and delete temporary things
    sqlDB := tdb.DBH()              // *sql.DB
    file := tdb.SqliteDB()          // absolute path to temporary db file
    present := tdb.HasTable("blah") // check if a table is present
    
    // lots of other methods too, see the godoc for more detail
}
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

```
Copyright 2024 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
