package command

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/mysqltool"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser/ast"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
)

func ValidateMysqlSelectSql(sql string) (string, string, bool, error) {
	sqlParser := parser.New()
	parsedStmt, err := sqlParser.ParseOneStmt(sql, "", "")
	if err != nil {
		return "", "", false, errors.New(i18n.GetByKey(i18n.SqlWrongSyntaxMsg))
	}
	isExplain := false
	switch parsedStmt.(type) {
	case *ast.SelectStmt:
		stmt := parsedStmt.(*ast.SelectStmt)
		if stmt.Limit != nil {
			return "", "", false, errors.New(i18n.GetByKey(i18n.SqlNotAllowHasLimitMsg))
		}
	case *ast.ExplainStmt:
		isExplain = true
	default:
		return "", "", false, errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	checker := new(mysqltool.Checker)
	parsedStmt.Accept(checker)
	tableNames := checker.GetTableNames()
	if len(tableNames) != 1 || tableNames[0].Schema.String() != "" {
		return "", "", false, errors.New(i18n.GetByKey(i18n.SqlUnsupportedMsg))
	}
	return tableNames[0].Name.String(), parsedStmt.Text(), isExplain, nil
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
