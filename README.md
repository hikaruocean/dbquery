# dbquery
golang mysql struct implement with go-sql-driver

## Installation
```bash
$ go get -u github.com/hikaruocean/dbquery
```

## Usage
```go
package main

import (
    "fmt"
    "github.com/hikaruocean/dbquery"
)

func main () {
    var dbquery = dbquery.New(map[string]string{"username": "root", "password": "mysqlPassWord", "host": "db", "dbname": "databaseName"})
    dbquery.SetConnect()
    bindData := make(map[string]interface{})
    bindData["enabled"] = "Y"
    bindData["deleted"] = "N"
    rh, err := dbquery.Query("SELECT * FROM account WHERE enabled = :enabled: AND deleted = :deleted:", bindData)
    if err != nil {
        panic(err.Error())
    }
    for row, err := rh.Fetch() ; err == nil && len(row) != 0 ; row, err = rh.Fetch() {
        fmt.Println(row)
    }
    rh.Close()
    if err != nil {
        panic(err.Error())
    }

    bindData = make(map[string]interface{})
    bindData["data1"] = "hikaru"
    bindData["created_at"] = "2019-01-01 00:00:00"
    bindData["updated_at"] = "2019-01-01 00:00:00"
    dbquery.Begin()
    rh, err = dbquery.Execute("INSERT INTO test (data1, created_at, updated_at) VALUES (:data1:, :created_at:, :updated_at:)", bindData)
    if err != nil {
        panic(err.Error())
    }
    id, err := rh.LastInsertId()
    num, err := rh.RowsAffected()
    rh.Close()
    dbquery.Commit()
    fmt.Println(id, num)
}
```

> row, err := rh.Fetch()
> row is map[string]interface{}
