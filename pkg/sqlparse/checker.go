package sqlparse

import "github.com/LeeZXin/zall/pkg/sqlparse/parser/ast"

// Checker 获取表名
type Checker struct {
	db         string
	tableNames []*ast.TableName
}

func (s *Checker) GetTableNames() []*ast.TableName {
	return s.tableNames
}

// Enter for node visit
func (s *Checker) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch nn := n.(type) {
	case *ast.TableName:
		s.tableNames = append(s.tableNames, nn)
	}
	return n, false
}

// Leave for node visit
func (s *Checker) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}
