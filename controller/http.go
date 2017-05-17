package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func InitHttpServer() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	go startHttpServer()
	return nil
}

func startHttpServer() error {
	router := NewRouter()
	return http.ListenAndServe(GlobalConfig.HTTP_SERVER, router)
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandleFunc
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"nodeStatus",
		"GET",
		"/node_status",
		nodeStatus,
	},
	Route{
		"queueStatus",
		"GET",
		"/queue_status",
		queueStatus,
	},
	Route{
		"clearQueue",
		"POST",
		"/clear_queue",
		clearQueue,
	},
	Route{
		"sendSwitch",
		"POST",
		"/send_switch",
		sendSwitch,
	},
	Route{
		"sendInterval",
		"POST",
		"/set_interval",
		setInterval,
	},
}

func nodeStatus(w http.ResponseWriter, r *http.Request) {
	for _, node := range controller.nodePool.Nodes {
		fmt.Fprintln(w, string(node.Encode())+",")
	}
}

func queueStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>`)
	fmt.Fprintln(w, fmt.Sprintf("<br>邮件队列: %d</br>", controller.mailQueue.Size()))
	fmt.Fprintln(w, fmt.Sprintf("<br>短信队列: %d</br>", controller.smsQueue.Size()))
	fmt.Fprintln(w, fmt.Sprintf("<br>微信队列: %d</br>", controller.wechatQueue.Size()))
	fmt.Fprintln(w, fmt.Sprintf("<br>电话队列: %d</br>", controller.callQueue.Size()))
	fmt.Fprintln(w, fmt.Sprintf("<br>自定义队列: %d</br>", controller.actionQueue.Size()))
	fmt.Fprintln(w, fmt.Sprintf("<br>发送时间间隔: %d 毫秒</br>", GlobalConfig.SEND_INTERVAL))
	fmt.Fprintln(w, fmt.Sprintf("<br>报警功能开启状态: %t</br>", GlobalConfig.SEND_SWITCH))
	fmt.Fprintln(w, `<br>
			     <form action="/clear_queue" method="post" enctype="multipart/form-data" >
			         <input type="submit" name="" value="清空队列" /> 
		             </form>`)
	fmt.Fprintln(w, `    <form action="/send_switch" method="post" enctype="multipart/form-data" >
			         <input type="hidden" name="switch" value="on" /> 
			         <input type="submit" name="" value="开启报警" /> 
		             </form>`)
	fmt.Fprintln(w, `    <form action="/send_switch" method="post" enctype="multipart/form-data" >
			         <input type="hidden" name="switch" value="off" /> 
			         <input type="submit" name="" value="关闭报警" /> 
		             </form>`)
	fmt.Fprintln(w, `    <form action="/set_interval" method="post" enctype="multipart/form-data" >
			         <input type="text" name="interval" value="" /> 毫秒  <input type="submit" name="" value="设定发送间隔" /> 
		             </form>
			 </br>`)
	fmt.Fprintln(w, `</html>`)
}

func clearQueue(w http.ResponseWriter, r *http.Request) {
	controller.mailQueue.Clear()
	controller.smsQueue.Clear()
	controller.wechatQueue.Clear()
	controller.callQueue.Clear()
	controller.actionQueue.Clear()
	http.Redirect(w, r, "/queue_status", http.StatusFound)
}

func sendSwitch(w http.ResponseWriter, r *http.Request) {
	switch flag := r.FormValue("switch"); {
	case flag == "on":
		GlobalConfig.SEND_SWITCH = true
	case flag == "off":
		GlobalConfig.SEND_SWITCH = false
	}
	http.Redirect(w, r, "/queue_status", http.StatusFound)
}

func setInterval(w http.ResponseWriter, r *http.Request) {
	interval := r.FormValue("interval")
	interval_int, err := strconv.Atoi(interval)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	if interval_int > 0 {
		GlobalConfig.SEND_INTERVAL = interval_int
	}
	http.Redirect(w, r, "/queue_status", http.StatusFound)
}
