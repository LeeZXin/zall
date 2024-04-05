package sqlparse

import (
	"fmt"
	"github.com/LeeZXin/zall/pkg/sqlparse/parser"
	"github.com/LeeZXin/zall/pkg/sqlparse/parser/ast"
	_ "github.com/LeeZXin/zall/pkg/sqlparse/parser/tidb-types/parser_driver"
)

func Parse(sql string) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		panic(err)
	}
	s, ok := stmt.(*ast.SelectStmt)
	if ok {
		fmt.Println("yes")
		checker := new(Checker)
		s.Accept(checker)
		fmt.Println(checker.GetTableNames()[0].Text())

	}
	fmt.Println(stmt.Text())

}
