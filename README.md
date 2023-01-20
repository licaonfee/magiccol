# magiccol

![Test](https://github.com/licaonfee/magiccol/workflows/Run%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/licaonfee/magiccol)](https://goreportcard.com/report/github.com/licaonfee/magiccol)
[![Coverage Status](https://coveralls.io/repos/github/licaonfee/magiccol/badge.svg?branch=master)](https://coveralls.io/github/licaonfee/magiccol?branch=master)
[![godoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/licaonfee/magiccol?tab=doc)

Dinamyc columns for database/sql. Magiccol allows to scan rows for a sql query, without known which columns it returns and create a `map[string]any`

## But Why?

In some cases is not possible to create a struct or variable to scan a value from a database query, then you can scan it into `any` but this lead to incorrect scanned values, for example MYSQL by default scan `DATETIME` into a `[]uint8` magiccol ensures that a DATETIME value will be scanned into a mysql.NullTime. This behaviour can be changed according to any `sql/db` driver.

## Example

```golang
package main

import (
    "database/sql"
    "fmt"
    "log"
    "reflect"

    "github.com/go-sql-driver/mysql"
    "github.com/licaonfee/magiccol"
)

func main() {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/employees?charset=utf8mb4")
    if err != nil {
        log.Fatal(err)
    }
    r, err := db.Query("select * from employees limit 10;")
    if err != nil {
        log.Fatal(err)
    }
    m := magiccol.DefaultMapper()
    //Use mysql native Time type see
    //https://github.com/go-sql-driver/mysql#timetime-support
    custom := reflect.TypeOf(mysql.NullTime{})
    match := []magiccol.Matcher{
        magiccol.DatabaseTypeAs("DATE", custom),
        magiccol.DatabaseTypeAs("DATETIME", custom),
        magiccol.DatabaseTypeAs("TIMESTAMP", custom),
    }
    m.Match(match)

    sc, err := magiccol.NewScanner(magiccol.Options{Rows:r, Mapper: m})
    if err != nil {
        log.Fatal(err)
    }
    var value map[string]any
    for sc.Scan() {
        value = sc.Value()
        fmt.Printf("%v\n", value)
    }
    if sc.Err() != nil {
        log.Fatal(sc.Err())
    }
}

```
