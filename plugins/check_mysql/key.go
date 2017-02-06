package main

const (
	COUNTER = "counter"
	GAUGE   = "gauge"
)

var globalStatusKeys = map[string]string{
	"Questions":         COUNTER,
	"Com_commit":        COUNTER,
	"Com_rollback":      COUNTER,
	"Com_begin":         COUNTER,
	"Com_select":        COUNTER,
	"Com_insert":        COUNTER,
	"Com_update":        COUNTER,
	"Com_delete":        COUNTER,
	"Com_repleace":      COUNTER,
	"Bytes_received":    COUNTER,
	"Bytes_send":        COUNTER,
	"Threads_running":   COUNTER,
	"Threads_connected": COUNTER,
	"Aborted_connects":  COUNTER,
	"Open_tables":       COUNTER,
	"Opend_table":       COUNTER,
}

var slaveStatusKey = map[string]string{
	"Slave_IO_Running":      GAUGE,
	"Slave_SQL_Running":     GAUGE,
	"Seconds_Behind_Master": GAUGE,
}

type Metric struct {
	Metric   string            `json:"metric"`    //sys.cpu.idle
	DataType string            `json:"data_type"` //COUNTER,GAUGE,DERIVE
	Value    interface{}       `json:"value"`     //99.00
	Tags     map[string]string `json:"tags"`      //{"product":"app01", "group":"dev02"}
}
