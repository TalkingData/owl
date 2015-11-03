package main

import (
	"fmt"
	"net"
	"time"
)

type Host struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	AssestId   int    `json:"assest_id"`
	Ip         string `json:"ip"`
	Group      string `json:"group"`
	Status     int    `json:"status"`
	Proxy      string `json:"proxy"`
	last_check string
	Services   []*Service `json:"services"`
}

type Service struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	Plugin       string           `json:"plugin"`
	Args         string           `json:"args"`
	ExecInterval int              `json:"exec_interval"`
	Items        map[string]*Item `json:"items"`
	//lastcheck    time.Time
}

type Item struct {
	Key string `json:"key"`
	DT  string `json:"data_type"` //GAUGE|COUNTER|DERIVE

}

type Port struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	ProcName string `json:"proc_name"`
	Status   int    `json:"status"`
}

func (this *Port) StatusCheck() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port), time.Second*5)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
