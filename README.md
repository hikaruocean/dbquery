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
    var dbquery = new(dbquery.DBquery)
    dbquery.Config(map[string]string{"username": "root", "password": "mysqlPassWord", "host": "db", "dbname": "databaseName"})
    dbquery.Connect()
    bindData := make(map[string]interface{})
    bindData["enabled"] = "Y"
    bindData["deleted"] = "N"
    rh, err := dbquery.Execute("SELECT * FROM account WHERE enabled = :enabled: AND deleted = :deleted:", bindData)
    if err != nil {
        panic(err.Error())
    }
    for row, err := rh.Fetch() ; err == nil && len(row) != 0 ; row, err = rh.Fetch() {
        fmt.Println(row)
    }
    if err != nil {
        panic(err.Error())
    }
}
```

> row, err := rh.Fetch()
> row is map[string]interface{}
