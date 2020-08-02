# magiccol

Dinamyc columns for database/sql. Magiccol allows to scan rows for a sql query, without known which columns it returns

Example

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
    m.Type(reflect.TypeOf(mysql.NullTime{}), "DATE", "DATETIME", "TIMESTAMP")

    sc, err := magiccol.NewScanner(magiccol.Options{Rows:r, Mapper: m})
    if err != nil {
        log.Fatal(err)
    }
    for sc.Scan() {
        value := sc.Value()
        fmt.Printf("%v\n", value)
    }
    if sc.Err() != nil {
        log.Fatal(sc.Err())
    }
}

```
