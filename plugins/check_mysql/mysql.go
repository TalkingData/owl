package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
)

var mysql *db

type db struct {
	*sql.DB
}

//dsn_format: user:password@[protocol](address)/dbname

func FetchData(c *cli.Context) error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8",
		user, password, host, port))
	if err != nil {
		return err
	}
	row, err := conn.Query("select 1")
	if err != nil {
		return err
	}
	defer row.Close()
	mysql = &db{conn}
	metrics := []Metric{}
	metrics = append(metrics, globalStatus()...)
	metrics = append(metrics, slaveStatus()...)
	res, _ := json.MarshalIndent(metrics, "", "    ")
	fmt.Println(string(res))
	return nil
}

func globalStatus() []Metric {
	rows, err := mysql.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var name string
	var value float64
	datas := []Metric{}
	for rows.Next() {
		if err := rows.Scan(&name, &value); err != nil {
			continue
		}
		if m, ok := globalStatusKeys[name]; ok {
			datas = append(datas, Metric{
				Metric:   fmt.Sprintf("mysql.%s", strings.ToLower(name)),
				DataType: m,
				Value:    value,
				Tags: map[string]string{
					"port": strconv.Itoa(port),
				},
			})
		}
	}
	return datas
}

func slaveStatus() []Metric {
	rows, err := mysql.Query("SHOW SLAVE STATUS")
	if err != nil {
		return nil
	}
	if !rows.Next() {
		return nil
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	for i := range values {
		var v sql.RawBytes
		values[i] = &v
	}
	err = rows.Scan(values...)
	if err != nil {
		fmt.Println("scan error ", err.Error())
		return nil
	}
	metrics := []Metric{}
	for i, name := range columns {
		bp := values[i].(*sql.RawBytes)
		vs := string(*bp)
		if dt, ok := slaveStatusKey[name]; ok {
			var v float64
			if vs == "Yes" {
				v = 1
			} else {
				v, _ = strconv.ParseFloat(vs, 64)
			}
			metrics = append(metrics, Metric{
				Metric:   fmt.Sprintf("mysql.%s", strings.ToLower(name)),
				DataType: dt,
				Value:    v,
				Tags: map[string]string{
					"port": strconv.Itoa(port),
				},
			})
		}
	}
	return metrics
}
