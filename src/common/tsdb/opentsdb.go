package tsdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OpenTsdbClient struct {
	url        *url.URL
	httpClient *http.Client
	tr         *http.Transport
}

func NewOpenTsdbClient(addr string, timeout time.Duration) (*OpenTsdbClient, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s", addr))
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}

	return &OpenTsdbClient{
		url: u,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
		tr: tr,
	}, nil
}

func (c *OpenTsdbClient) Close() error {
	c.tr.CloseIdleConnections()
	return nil
}

type Query struct {
	Aggregator string            `json:"aggregator"`
	Metric     string            `json:"metric"`
	Rate       bool              `json:"rate,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
}

type QueryParams struct {
	Start             interface{} `json:"start"`
	End               interface{} `json:"end,omitempty"`
	Queries           []Query     `json:"queries,omitempty"`
	NoAnnotations     bool        `json:"no_annotations,omitempty"`
	GlobalAnnotations bool        `json:"global_annotations,omitempty"`
	MsResolution      bool        `json:"ms,omitempty"`
	ShowTSUIDs        bool        `json:"show_tsuids,omitempty"`
	ShowSummary       bool        `json:"show_summary,omitempty"`
	ShowQuery         bool        `json:"show_query,omitempty"`
	Delete            bool        `json:"delete,omitempty"`
}

type InnerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorResp struct {
	Error InnerError `json:"error"`
}

func (this ErrorResp) String() string {
	return fmt.Sprintf("{code: %d, message: %s, detail: %s}", this.Error.Code, this.Error.Message, this.Error.Details)
}

func (c *OpenTsdbClient) newQueryParams(start, end, rawTags, aggregator, metric string, isRelative bool) *QueryParams {
	tags := make(map[string]string)
	if rawTags != "" {
		tagsPairs := strings.Split(rawTags, ",")
		for _, tagPair := range tagsPairs {
			tagKV := strings.Split(tagPair, "=")
			tags[tagKV[0]] = tagKV[1]
		}
	}
	if isRelative {
		start = fmt.Sprintf("%sm-ago", start)
	}
	queries := []Query{Query{Aggregator: aggregator, Metric: metric, Tags: tags}}
	return &QueryParams{Start: start, End: end, Queries: queries}
}

func (c *OpenTsdbClient) Query(start, end, rawTags, aggregator, metric string, isRelative bool) ([]Result, error) {
	q := c.newQueryParams(start, end, rawTags, aggregator, metric, isRelative)
	data, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	u := c.url
	u.Path = "api/query"

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var results []Result
		if err := json.Unmarshal(body, &results); err != nil {
			return nil, err
		}
		return results, nil
	}
	var errResp ErrorResp
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, errors.New(string(data))
	}
	return nil, errors.New((&errResp).String())
}
