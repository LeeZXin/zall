package command

import (
	"database/sql"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/mysqltool"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser/ast"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
)

func ValidateMysqlSelectSql(sql string) (string, string, error) {
	p := parser.New()
	s, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlWrongSyntaxMsg))
	}
	stmt, ok := s.(*ast.SelectStmt)
	if !ok {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	if stmt.Limit != nil {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlNotAllowHasLimitMsg))
	}
	checker := new(mysqltool.Checker)
	stmt.Accept(checker)
	tableNames := checker.GetTableNames()
	if len(tableNames) != 1 {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	if tableNames[0].Schema.String() != "" {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	return tableNames[0].Name.String(), stmt.Text(), nil
}

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
