package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/pingcap/errors"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// MatrixRangeQuery defines a log range query. only for matrix
type MatrixRangeQuery struct {
	Start time.Time
	End   time.Time
	Step  int
	Query string
	Limit uint32
	//Shards    []string
}

func (q *MatrixRangeQuery) ToUrlQuery() string {
	query := fmt.Sprintf("start=%d", q.Start.Unix()) +
		fmt.Sprintf("&end=%d", q.End.Unix()) +
		"&query=" + url.QueryEscape(q.Query)
	if q.Step > 0 {
		query += fmt.Sprintf("&step=%d", q.Step)
	}
	if q.Limit > 0 {
		query += fmt.Sprintf("&limit=%d", q.Limit)
	} else {
		// 防止非matrix query请求造成大量数据返回
		query += "&limit=1"
	}
	return query
}

func (q *MatrixRangeQuery) DoRequest(ctx context.Context, httpClient *http.Client, queryUrl, orgId string) (MatrixQueryResponse, error) {
	var res MatrixQueryResponse
	err := httputil.Get(
		ctx,
		httpClient,
		queryUrl+"?"+q.ToUrlQuery(),
		map[string]string{
			"X-Scope-OrgID": orgId,
		},
		&res,
	)
	if err != nil {
		return MatrixQueryResponse{}, err
	}
	if res.Status != "success" || res.Data.ResultType != "matrix" {
		return MatrixQueryResponse{}, errors.New("fail")
	}
	return res, nil
}

const (
	minimumTick = time.Millisecond
	second      = int64(time.Second / minimumTick)
)

var (
	dotPrecision = int(math.Log10(float64(second)))
)

type Time int64

// String returns a string representation of the Time.
func (t Time) String() string {
	return strconv.FormatFloat(float64(t)/float64(second), 'f', -1, 64)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(b []byte) error {
	p := strings.Split(string(b), ".")
	switch len(p) {
	case 1:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		*t = Time(v * second)

	case 2:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		v *= second

		prec := dotPrecision - len(p[1])
		if prec < 0 {
			p[1] = p[1][:dotPrecision]
		} else if prec > 0 {
			p[1] = p[1] + strings.Repeat("0", prec)
		}

		va, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil {
			return err
		}

		// If the value was something like -0.1 the negative is lost in the
		// parsing because of the leading zero, this ensures that we capture it.
		if len(p[0]) > 0 && p[0][0] == '-' && v+va > 0 {
			*t = Time(v+va) * -1
		} else {
			*t = Time(v + va)
		}

	default:
		return fmt.Errorf("invalid time %q", string(b))
	}
	return nil
}

// A SampleValue is a representation of a value for a given sample at a given
// time.
type SampleValue float64

// MarshalJSON implements json.Marshaler.
func (v SampleValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *SampleValue) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("sample value must be a quoted string")
	}
	f, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
	if err != nil {
		return err
	}
	*v = SampleValue(f)
	return nil
}

func (v SampleValue) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct {
	Timestamp Time
	Value     SampleValue
}

func (s SamplePair) MarshalJSON() ([]byte, error) {
	t, err := json.Marshal(s.Timestamp)
	if err != nil {
		return nil, err
	}
	v, err := json.Marshal(s.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

type SampleStream struct {
	Metric map[string]string `json:"metric"`
	Values []SamplePair      `json:"values"`
}

// MatrixQueryResponseData only for resultType == "matrix"
type MatrixQueryResponseData struct {
	ResultType string         `json:"resultType"`
	Result     []SampleStream `json:"result"`
}

// MatrixQueryResponse represents the http json response to a Loki range query only for matrix
type MatrixQueryResponse struct {
	Status string                  `json:"status"`
	Data   MatrixQueryResponseData `json:"data"`
}

func (r *MatrixQueryResponse) SumAllValue() float64 {
	var ret float64 = 0
	for _, stream := range r.Data.Result {
		for _, val := range stream.Values {
			ret += float64(val.Value)
		}
	}
	return ret
}
