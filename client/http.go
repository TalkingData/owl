package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"owl/common/types"
	"strings"
	"time"
)

func startHttpMetrcs() {
	var port string
	mblist := strings.Split(GlobalConfig.MetricBind, ":")
	if len(mblist) == 2 {
		port = mblist[1]
	} else {
		lg.Error("metric_bind in the config file is incorrect")
		return
	}

	if port != "0" {
		http.HandleFunc("/", metricsHandler)
		lg.Info("start http listen on: %s", GlobalConfig.MetricBind)
		err := http.ListenAndServe(GlobalConfig.MetricBind, nil)
		lg.Error("start metric interface error: %s", err)
	} else {
		return
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var tsd types.TimeSeriesData
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&tsd)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err.Error()), 400)
		return
	}
	tsd.Timestamp = time.Now().Unix()
	lg.Info("get json data", tsd)
	agent.SendChan <- tsd

	w.Write([]byte(`{"status":"ok"}`))

}
