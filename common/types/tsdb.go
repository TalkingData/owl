package types

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

type Options struct {
	Addr    string
	Timeout time.Duration
}

type Client struct {
	url        *url.URL
	httpClient *http.Client
	tr         *http.Transport
}

func NewTsdbClient(addr string, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s", addr))
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}

	return &Client{
		url: u,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
		tr: tr,
	}, nil
}

func (c *Client) Close() error {
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

var results []Result

type Result struct {
	Metric        string             `json:"metric"`
	Tags          map[string]string  `json:"tags"`
	AggregateTags []string           `json:"aggregateTags"`
	Dps           map[string]float64 `json:"dps"`
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

func NewQueryParams(start, end string, rawTags string, aggregator string, metric string) *QueryParams {
	tags := make(map[string]string)
	if rawTags != "" {
		tags_pairs := strings.Split(rawTags, ",")
		for _, tag_pair := range tags_pairs {
			tag_k_v := strings.Split(tag_pair, "=")
			tags[tag_k_v[0]] = tag_k_v[1]
		}
	}
	queries := []Query{Query{Aggregator: aggregator, Metric: metric, Tags: tags}}
	return &QueryParams{Start: start, End: end, Queries: queries}
}

func (c *Client) Query(q *QueryParams) ([]Result, error) {
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
	} else {
		var err_resp ErrorResp
		if err := json.Unmarshal(body, &err_resp); err != nil {
			return nil, errors.New(string(data))
		}
		return nil, errors.New((&err_resp).String())
	}
}
