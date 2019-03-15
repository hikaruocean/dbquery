package dbquery

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "regexp"
    "fmt"
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
    fmt.Println(this.config)
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

func (this *DBquery) Execute (sqlStr string,params map[string]interface{}) (ResultHandler, error) {
    var rh ResultHandler
    realSql, markSortAry := this.getRealSql(sqlStr)
    sth, err := this.SthProcess(realSql)
    if err != nil {
        return rh, err
    }
    defer sth.Close()
    this.sth = sth

    bind := make([]interface{}, 0)
    for _, val := range markSortAry {
        bind = append(bind, params[val])
    }
    rows, err := sth.Query(bind...)
    // defer rows.Close()
    if err != nil {
        panic(err.Error())
        return rh, err
    }
    rh.rows = rows
    return rh, nil
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
