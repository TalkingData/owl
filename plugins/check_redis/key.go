package main

const (
	COUNTER = "counter"
	GAUGE   = "gauge"
)

var infoKeys = map[string]string{
	"mem_fragmentation_ratio":   GAUGE,
	"rdb_bgsave_in_progress":    GAUGE,
	"instantaneous_ops_per_sec": GAUGE,
	"total_commands_processed":  GAUGE,
	"connected_clients":         GAUGE,
	"master_link_status":        GAUGE,
	"connected_slaves":          GAUGE,
	"used_memory":               GAUGE,
}

type Metric struct {
	Metric   string            `json:"metric"`    //sys.cpu.idle
	DataType string            `json:"data_type"` //COUNTER,GAUGE,DERIVE
	Value    interface{}       `json:"value"`     //99.00
	Tags     map[string]string `json:"tags"`      //{"product":"app01", "group":"dev02"}
}
