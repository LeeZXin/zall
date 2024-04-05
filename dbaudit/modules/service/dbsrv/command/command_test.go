package command

import (
	"fmt"
	"testing"
)

func TestMysqlExecutor_Execute(t *testing.T) {
	e := MysqlExecutor{
		DatasourceName: "root:root@tcp(127.0.0.1:3306)/hhhh?charset=utf8",
	}
	fmt.Println(e.Execute("select * from zall_user limit 2"))

}

func TestValidateMysqlSelectSql(t *testing.T) {
	fmt.Println(ValidateMysqlSelectSql("select * from zall_user       ;;;; ; ;"))
}
