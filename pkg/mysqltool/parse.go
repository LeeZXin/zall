package mysqltool

import (
	"fmt"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser/ast"
	_ "github.com/LeeZXin/zall/pkg/mysqltool/parser/tidb-types/parser_driver"
)

func Parse(sql string) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		panic(err)
	}
	s, ok := stmt.(*ast.SelectStmt)
	if ok {
		checker := new(Checker)
		s.Accept(checker)
		fmt.Println(checker.GetTableNames()[0].Text())

	}
	fmt.Println(stmt.Text())

}
