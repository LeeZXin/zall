package alert

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/loki"
	"github.com/LeeZXin/zall/pkg/mysqltool"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser"
	"github.com/LeeZXin/zall/pkg/mysqltool/parser/ast"
	"github.com/LeeZXin/zall/pkg/prom"
	"github.com/LeeZXin/zall/util"
	"github.com/pingcap/errors"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

type SourceType int

const (
	MysqlType SourceType = iota + 1
	PromType
	LokiType
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
		validateMysqlSelectSql(c.SelectSql) &&
		c.Database != "" &&
		c.Condition != ""
}

func (c *MysqlConfig) Execute(endTime time.Time) (map[string]any, error) {
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Database,
	)
	selectSql := c.SelectSql
	tpl, err := template.New("").Parse(c.SelectSql)
	if err == nil {
		msg := new(bytes.Buffer)
		err = tpl.Execute(msg, map[string]any{
			"EndTime": endTime.Format(time.DateTime),
		})
		if err == nil {
			selectSql = msg.String()
		}
	}
	result, err := util.MysqlQuery(datasourceName, strings.TrimSuffix(selectSql, ";")+" limit 1")
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
	Host      string `json:"host"`
	PromQl    string `json:"promQl"`
	Condition string `json:"condition"`
}

func (c *PromConfig) IsValid() bool {
	parsedUrl, err := url.Parse(c.Host)
	if err != nil {
		return false
	}
	return strings.HasPrefix(parsedUrl.Scheme, "http") && c.PromQl != "" && c.Condition != ""
}

func (c *PromConfig) Execute(httpClient *http.Client, endTime time.Time) (*model.Sample, error) {
	client, err := prom.NewPromHttpClient(c.Host, httpClient)
	if err != nil {
		return nil, err
	}
	result, _, err := promapi.NewAPI(client).Query(context.Background(), c.PromQl, endTime)
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

type LokiConfig struct {
	Host         string `json:"host"`
	LogQl        string `json:"logQl"`
	OrgId        string `json:"orgId"`
	LastDuration string `json:"lastDuration"`
	Step         int    `json:"step"`
	Condition    string `json:"condition"`
}

func (c *LokiConfig) IsValid() bool {
	parsedUrl, err := url.Parse(c.Host)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return false
	}
	if c.LogQl == "" {
		return false
	}
	duration, err := time.ParseDuration(c.LastDuration)
	if err != nil || duration <= 0 || duration > time.Hour {
		return false
	}
	if c.Step <= 0 || c.Step > 3600 {
		return false
	}
	if c.Condition == "" {
		return false
	}
	return true
}

func (c *LokiConfig) Execute(httpClient *http.Client, endTime time.Time) (time.Time, float64, error) {
	duration, err := time.ParseDuration(c.LastDuration)
	if err != nil {
		return time.Time{}, 0, err
	}
	startTime := endTime.Add(-duration)
	query := loki.MatrixRangeQuery{
		Start: startTime,
		End:   endTime,
		Step:  c.Step,
		Query: c.LogQl,
		Limit: 1,
	}
	response, err := query.DoRequest(
		context.Background(),
		httpClient,
		strings.TrimSuffix(c.Host, "/")+"/loki/api/v1/query_range",
		c.OrgId,
	)
	if err != nil {
		return startTime, 0, err
	}
	return startTime, response.SumAllValue(), nil
}

type Alert struct {
	SourceType SourceType   `json:"sourceType"`
	Mysql      *MysqlConfig `json:"mysql,omitempty"`
	Prom       *PromConfig  `json:"prom,omitempty"`
	Loki       *LokiConfig  `json:"loki,omitempty"`
	commonhook.TypeAndCfg
}

func (a *Alert) IsValid() bool {
	if !a.TypeAndCfg.IsValid() {
		return false
	}
	switch a.SourceType {
	case MysqlType:
		return a.Mysql != nil && a.Mysql.IsValid()
	case PromType:
		return a.Prom != nil && a.Prom.IsValid()
	case LokiType:
		return a.Loki != nil && a.Loki.IsValid()
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
