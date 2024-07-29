package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/mysqltool"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser/ast"
	"github.com/LeeZXin/zall/util"
	"github.com/pingcap/errors"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SourceType int

const (
	MysqlType SourceType = iota + 1
	PromType
)

func (t SourceType) Readable() string {
	switch t {
	case MysqlType:
		return "mysql"
	case PromType:
		return "prom"
	default:
		return "unknown"
	}
}

type MysqlConfig struct {
	Host      string `json:"host"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	SelectSql string `json:"selectSql"`
	Condition string `json:"condition"`
}

func (c *MysqlConfig) IsValid() bool {
	return util.GenIpPortPattern().MatchString(c.Host) &&
		c.Username != "" && c.Password != "" &&
		validateMysqlSelectSql(c.SelectSql) && c.Database != "" &&
		c.Condition != ""
}

func (c *MysqlConfig) Execute() (map[string]string, error) {
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Database,
	)
	result, err := util.MysqlQuery(datasourceName, strings.TrimSuffix(c.SelectSql, ";")+" limit 1")
	if err != nil {
		return nil, err
	}
	ret := result.ToMap()
	if len(ret) > 0 {
		return ret[0], nil
	}
	return nil, nil
}

type PromConfig struct {
	HostUrl   string `json:"HostUrl"`
	PromQl    string `json:"promQl"`
	Condition string `json:"condition"`
}

func (c *PromConfig) IsValid() bool {
	parsedUrl, err := url.Parse(c.HostUrl)
	if err != nil {
		return false
	}
	return strings.HasPrefix(parsedUrl.Scheme, "http") && c.PromQl != "" && c.Condition != ""
}

func (c *PromConfig) Execute(httpClient *http.Client) (*model.Sample, error) {
	client, err := util.NewPromHttpClient(c.HostUrl, httpClient)
	if err != nil {
		return nil, err
	}
	api := promapi.NewAPI(client)
	result, _, err := api.Query(context.Background(), c.PromQl, time.Now())
	if err != nil {
		return nil, err
	}
	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		if len(vector) > 0 {
			return vector[0], nil
		}
		return nil, nil
	default:
		return nil, errors.New("unexpected result format")
	}
}

type Alert struct {
	Source SourceType   `json:"source"`
	Mysql  *MysqlConfig `json:"mysql"`
	Prom   *PromConfig  `json:"prom"`
	Api    util.Api     `json:"api"`
}

func (a *Alert) FromDB(content []byte) error {
	if a == nil {
		*a = Alert{}
	}
	return json.Unmarshal(content, a)
}

func (a *Alert) ToDB() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Alert) IsValid() bool {
	if !a.Api.IsValid() {
		return false
	}
	switch a.Source {
	case MysqlType:
		return a.Mysql != nil && a.Mysql.IsValid()
	case PromType:
		return a.Prom != nil && a.Prom.IsValid()
	default:
		return false
	}
}

func validateMysqlSelectSql(sql string) bool {
	sqlParser := parser.New()
	parsedStmt, err := sqlParser.ParseOneStmt(sql, "", "")
	// 语法错误
	if err != nil {
		return false
	}
	switch parsedStmt.(type) {
	case *ast.SelectStmt:
		stmt := parsedStmt.(*ast.SelectStmt)
		// 存在limit
		if stmt.Limit != nil {
			return false
		}
	default:
		// 只支持select
		return false
	}
	checker := new(mysqltool.Checker)
	parsedStmt.Accept(checker)
	tableNames := checker.GetTableNames()
	// 必须有表名
	if len(tableNames) != 1 || tableNames[0].Schema.String() != "" {
		return false
	}
	return true
}
