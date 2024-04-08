package util

import (
	"database/sql"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
)

func MysqlQuery(datasourceName, cmd string) ([]string, [][]string, error) {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()
	rows, err := db.Query(cmd)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	ret := make([][]string, 0)
	values := make([][]byte, len(columns))
	scans := make([]any, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	for rows.Next() {
		row := make([]string, len(columns))
		err = rows.Scan(scans...)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range values {
			row[i] = string(v)
		}
		ret = append(ret, row)
	}
	return columns, ret, nil
}

type MysqlExecuteResult struct {
	Sql          string
	AffectedRows int64
	ErrMsg       string
}

func MysqlExecute(datasourceName string, cmd string) ([]MysqlExecuteResult, error) {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ret := make([]MysqlExecuteResult, 0)
	p := parser.New()
	stmts, warns, err := p.Parse(cmd, "", "")
	if err != nil {
		return nil, err
	}
	if len(warns) > 0 {
		return nil, warns[0]
	}
	for _, stmt := range stmts {
		executeSql := stmt.Text()
		er := MysqlExecuteResult{
			Sql: executeSql,
		}
		result, err := db.Exec(executeSql)
		if err != nil {
			er.ErrMsg = err.Error()
		} else {
			er.AffectedRows, err = result.RowsAffected()
			if err != nil {
				er.ErrMsg = err.Error()
			}
		}
		ret = append(ret, er)
	}
	return ret, nil
}
