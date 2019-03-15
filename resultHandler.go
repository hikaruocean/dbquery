package dbquery

import (
    "database/sql"
    "fmt"
)

type ResultHandler struct {
    rows *sql.Rows
}

func (this *ResultHandler) Close() {
    this.rows.Close()
}

func (this *ResultHandler) Fetch() (DBqueryRow, error) {
    assoc := make(DBqueryRow)
    colNameAry, err := this.rows.Columns()
    if err != nil {
        return assoc, err
    }
    colTypesAry, err := this.rows.ColumnTypes()
    if err != nil {
        return assoc, err
    }
    colTypesStrAry := make([]string, 0)
    for _, v := range colTypesAry {
        colTypesStrAry = append(colTypesStrAry, v.DatabaseTypeName())
    }
    len := len(colNameAry)
    results := make([]interface{}, len)
    ptrAry := make([]interface{}, len)
    for i := range ptrAry {
        ptrAry[i] = &results[i]
    }

    if !this.rows.Next() {
        defer this.Close()
        return assoc, nil
    }

    if err := this.rows.Scan(ptrAry...); err != nil {
        return assoc, err
    }
    for i := range results {
        assoc[colNameAry[i]] = results[i]
        switch colTypesStrAry[i] {
            case "CHAR":
                fallthrough
            case "VARCHAR":
                fallthrough
            case "TEXT":
                val, _ := assoc.String(colNameAry[i])
                assoc[colNameAry[i]] = val
            case "TIME":
                fallthrough
            case "DATETIME":
                val, _ := assoc.Time(colNameAry[i])
                assoc[colNameAry[i]] = val
        }
    }
    return assoc, nil
}
