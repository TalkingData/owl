package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"owl/dto"
	"time"
)

// kairosRestTsData struct
type kairosRestTsData struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Cycle     int               `json:"cycle,omitempty"`
	Tags      map[string]string `json:"tags"`
}

// kairosdbRestBackend struct
type kairosdbRestBackend struct {
	tcpAddr *net.TCPAddr
}

// newKairosdbRestBackend
func newKairosdbRestBackend(addr string) (Backend, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	bEnd := &kairosdbRestBackend{
		tcpAddr: tcpAddr,
	}
	return bEnd, nil
}

// Write
func (rkb *kairosdbRestBackend) Write(in *dto.TsData) error {
	kTsData := &kairosRestTsData{
		Name:      in.Metric,
		Timestamp: in.Timestamp * 1000,
		Value:     in.Value,
		Tags:      in.Tags,
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s/api/v1/datapoints", rkb.tcpAddr.String())

	data, err := json.Marshal(kTsData)
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
