package dbquery

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "regexp"
)

type DBquery struct {
    db *sql.DB
    sth *sql.Stmt
    dsn string
    config map[string]string
}

func (this *DBquery) Config (config map[string]string) {
    this.dsn = ""
    this.config = config
}

func (this *DBquery) SetDSN () {
    if val, isset := this.config["proto"] ; !isset || val == "" {
        this.config["proto"] = "tcp"
    }
    if val, isset := this.config["port"] ; !isset || val == "" {
        this.config["port"] = "3306"
    }
    if val, isset := this.config["charset"] ; !isset || val == "" {
        this.config["charset"] = "utf8"
    }
    if val, isset := this.config["collation"] ; !isset || val == "" {
        this.config["collation"] = "utf8_general_ci"
    }

    this.dsn = this.config["username"] + ":" + this.config["password"] + "@" + this.config["proto"] + "(" + this.config["host"] + ":" + this.config["port"] + ")/" + this.config["dbname"] + "?charset=" + this.config["charset"] + "&collation=" + this.config["collation"]
}

func (this *DBquery) Connect () (bool, error) {
    if this.dsn == "" {
        this.SetDSN()
    }
    db, err := sql.Open("mysql", this.dsn)
    if (err != nil) {
        return false, err
    }
    this.db = db
    return true, nil
}

func (this *DBquery) Query (sqlStr string,params map[string]interface{}) (ResultHandler, error) {
    var rh ResultHandler
    realSql, markSortAry := this.getRealSql(sqlStr)
    bind := this.getBindData(markSortAry, params)
    sth, err := this.SthProcess(realSql)
    if err != nil {
        return rh, err
    }
    defer sth.Close()
    this.sth = sth

    rows, err := sth.Query(bind...)
    if err != nil {
        panic(err.Error())
        return rh, err
    }
    rh.rows = rows
    return rh, nil
}

func (this *DBquery) Execute (sqlStr string,params map[string]interface{}) (ResultHandler, error) {
    var rh ResultHandler
    realSql, markSortAry := this.getRealSql(sqlStr)
    bind := this.getBindData(markSortAry, params)
    sth, err := this.SthProcess(realSql)
    if err != nil {
        return rh, err
    }
    defer sth.Close()
    this.sth = sth

    result, err := sth.Exec(bind...)
    if err != nil {
        panic(err.Error())
        return rh, err
    }
    rh.result = result
    return rh, nil
}

func (this *DBquery) Insert (table string, data map[string]interface{}) (ResultHandler, error) {
    colStr := ""
    placeholderStr := ""
    for colName := range data {
        if colStr != "" {
            colStr += ", "
            placeholderStr += ", "
        }
        colStr += colName
        placeholderStr += ":" + colName + ":"
    }
    sqlStr := "INSERT INTO " + table + " (" + colStr + ") VALUES (" + placeholderStr + ")"
    rh, err := this.Execute(sqlStr, data)
    return rh, err
}

func (this *DBquery) Update (table string, data map[string]interface{}, conditionStr string, cdata map[string]interface{}) (ResultHandler, error) {
    setStr := ""
    bindData := make(map[string]interface{})
    for colName := range data {
        if setStr != "" {
            setStr += ", "
        }
        setStr += colName + " = :d_" + colName + ": "
        bindData["d_" + colName] = data[colName]
    }

    sqlStr := "UPDATE " + table + " SET " + setStr + " WHERE " + conditionStr
    for placeholder, val := range cdata {
        bindData[placeholder] = val
    }
    rh, err := this.Execute(sqlStr, bindData)
    return rh, err
}

func (this *DBquery) Delete (table string, conditionStr string, cdata map[string]interface{}) (ResultHandler, error) {

    sqlStr := "DELETE FROM " + table + " WHERE " + conditionStr
    rh, err := this.Execute(sqlStr, cdata)
    return rh, err
}

func (this *DBquery) SthProcess (sqlStr string) (*sql.Stmt, error) {

    stmt, err := this.db.Prepare(sqlStr)
    return stmt, err
}

func (this *DBquery) getRealSql (sqlStr string) (string, []string){
    markSortAry := make([]string, 0)
    re := regexp.MustCompile(`:([a-zA-Z_]+[a-zA-Z0-9_]+):`)
    matchAryAry := re.FindAllStringSubmatch(sqlStr, -1)
    for _, matchAry := range matchAryAry {
        markSortAry = append(markSortAry, matchAry[1])
    }
    realSql := re.ReplaceAllString(sqlStr, "?")
    return realSql, markSortAry
}

func (this *DBquery) getBindData (markSortAry []string, params map[string]interface{}) ([]interface{}){
    bind := make([]interface{}, 0)
    for _, val := range markSortAry {
        data, isset := params[val]
        if !isset {
            panic("Placeholder not found in bind data")
        }
        bind = append(bind, data)
    }
    return bind
}
