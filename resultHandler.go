package dbquery

import (
    "database/sql"
)

type ResultHandler struct {
    rows *sql.Rows
    result sql.Result
    sth *sql.Stmt
}

func (this *ResultHandler) Close() {
    if this.rows != nil {
        for this.rows.Next() {
            this.rows.Scan()
        }
        this.rows.Close()
    }
    this.sth.Close()
}

func (this *ResultHandler) LastInsertId() (int64, error) {
    if this.result == nil {
        return 0, nil
    }
    id, err := this.result.LastInsertId()
    return id, err
}

func (this *ResultHandler) RowsAffected() (int64, error) {
    if this.result == nil {
        return 0, nil
    }
    num, err := this.result.RowsAffected()
    return num, err
}

func (this *ResultHandler) Fetch() (DBqueryRow, error) {
    assoc := make(DBqueryRow)
    if this.rows == nil {
        return assoc, nil
    }
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
