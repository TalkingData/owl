package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/redis.v5"
)

func FetchData(c *cli.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0, // use default DB
	})
	info := client.Info().String()
	s := strings.NewReader(info)
	buf := bufio.NewReader(s)
	metrics := []Metric{}
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		kv := strings.Split(strings.TrimSpace(line), ":")
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		if t, ok := infoKeys[key]; ok {
			metric := Metric{
				Metric:   fmt.Sprintf("redis.%s", strings.ToLower(key)),
				DataType: t,
				Value:    0,
				Tags: map[string]string{
					"port": port,
				},
			}
			if value == "up" {
				metric.Value = 1
			} else {
				metric.Value, _ = strconv.ParseFloat(value, 64)
			}
			metrics = append(metrics, metric)
		}
	}
	res, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}
