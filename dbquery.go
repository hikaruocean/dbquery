package dbquery

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "regexp"
)

type DBquery struct {
    db *sql.DB
    sth *sql.Stmt
}

func (this *DBquery) Connect () (bool, error) {
    db, err := sql.Open("mysql", "root:my19850126sql@(db)/soteria_cloud?charset=utf8&collation=utf8_general_ci")
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
