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
	sqlParser := parser.New()
	parsedStmt, err := sqlParser.ParseOneStmt(sql, "", "")
	if err != nil {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlWrongSyntaxMsg))
	}
	switch parsedStmt.(type) {
	case *ast.SelectStmt:
		stmt := parsedStmt.(*ast.SelectStmt)
		if stmt.Limit != nil {
			return "", "", errors.New(i18n.GetByKey(i18n.SqlNotAllowHasLimitMsg))
		}
	case *ast.ExplainStmt:
	default:
		return "", "", errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	checker := new(mysqltool.Checker)
	parsedStmt.Accept(checker)
	tableNames := checker.GetTableNames()
	if len(tableNames) != 1 || tableNames[0].Schema.String() != "" {
		return "", "", errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	return tableNames[0].Name.String(), parsedStmt.Text(), nil
}

type ValidateUpdateResult struct {
	TableName string `json:"tableName"`
	Sql       string `json:"sql"`
	Pass      bool   `json:"pass"`
	ErrMsg    string `json:"errMsg"`
}

func ValidateMysqlUpdateSql(sql string) ([]ValidateUpdateResult, bool, error) {
	sqlParser := parser.New()
	parsedStmts, warns, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		return nil, false, errors.New(i18n.GetByKey(i18n.SqlWrongSyntaxMsg))
	}
	if len(warns) > 0 {
		return nil, false, warns[0]
	}
	allPass := true
	ret := make([]ValidateUpdateResult, 0)
	for _, stmt := range parsedStmts {
		result := ValidateUpdateResult{
			Sql: stmt.Text(),
		}
		checker := new(mysqltool.Checker)
		stmt.Accept(checker)
		tableNames := checker.GetTableNames()
		if len(tableNames) != 1 || tableNames[0].Schema.String() != "" {
			result.ErrMsg = i18n.GetByKey(i18n.SqlUnsupportedMsg)
		} else {
			result.TableName = tableNames[0].Name.String()
			switch stmt.(type) {
			case *ast.CreateTableStmt, *ast.InsertStmt, *ast.AlterTableStmt:
				result.Pass = true
			case *ast.DeleteStmt:
				d := stmt.(*ast.DeleteStmt)
				if d.Where == nil {
					result.ErrMsg = i18n.GetByKey(i18n.SqlNotAllowNoWhereMsg)
				} else {
					result.Pass = true
				}
			case *ast.UpdateStmt:
				u := stmt.(*ast.UpdateStmt)
				if u.Where == nil {
					result.ErrMsg = i18n.GetByKey(i18n.SqlNotAllowNoWhereMsg)
				} else {
					result.Pass = true
				}
			default:
				result.ErrMsg = i18n.GetByKey(i18n.SqlUnsupportedMsg)
			}
		}
		if result.ErrMsg != "" {
			allPass = false
		}
		ret = append(ret, result)
	}
	return ret, allPass, nil
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
