package util

import (
	"database/sql"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
)

type MysqlQueryResult struct {
	Columns []string
	Data    [][]string
}

type MysqlQueryResultWithErr struct {
	MysqlQueryResult
	Err error
}

func (r *MysqlQueryResult) ToMap() []map[string]string {
	ret := make([]map[string]string, 0, len(r.Data))
	for _, datum := range r.Data {
		if len(datum) != len(r.Columns) {
			continue
		}
		item := make(map[string]string, len(r.Columns))
		for i := range r.Columns {
			item[r.Columns[i]] = datum[i]
		}
		ret = append(ret, item)
	}
	return ret
}

func MysqlQuery(datasourceName, cmd string) (MysqlQueryResult, error) {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return MysqlQueryResult{}, err
	}
	defer db.Close()
	return query(db, cmd)
}

func MysqlQueries(datasourceName string, cmds ...string) ([]MysqlQueryResultWithErr, error) {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ret := make([]MysqlQueryResultWithErr, 0, len(cmds))
	for _, cmd := range cmds {
		result, err := query(db, cmd)
		ret = append(ret, MysqlQueryResultWithErr{
			MysqlQueryResult: result,
			Err:              err,
		})
	}
	return ret, nil
}

func query(db *sql.DB, cmd string) (MysqlQueryResult, error) {
	rows, err := db.Query(cmd)
	if err != nil {
		return MysqlQueryResult{}, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return MysqlQueryResult{}, err
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
			return MysqlQueryResult{}, err
		}
		for i, v := range values {
			row[i] = string(v)
		}
		ret = append(ret, row)
	}
	return MysqlQueryResult{
		Columns: columns,
		Data:    ret,
	}, err
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
