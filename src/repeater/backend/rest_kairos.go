package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"owl/common/types"
	"time"
)

type restKairosTimeSeriesData struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Cycle     int               `json:"cycle,omitempty"`
	Tags      map[string]string `json:"tags"`
}

type RestKairosdbBackend struct {
	tcpAddr *net.TCPAddr
}

func NewRestKairosdbBackend(addr string) (*RestKairosdbBackend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	backend := new(RestKairosdbBackend)
	backend.tcpAddr = tcpAddr
	return backend, nil
}

func (rkb *RestKairosdbBackend) Write(in *types.TimeSeriesData) error {

	dto := &restKairosTimeSeriesData{
		Name:      in.Metric,
		Timestamp: in.Timestamp * 1000,
		Value:     in.Value,
		Tags:      in.Tags,
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s/api/v1/datapoints", rkb.tcpAddr.String())

	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Failed to wirte rest kairosdb, status code: %d, response: %s", resp.StatusCode, resp.Body)
	}

	return nil
}
