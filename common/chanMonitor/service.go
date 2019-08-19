package chanMonitor

import (
	"encoding/json"
	"log"
	"net/http"
)

type Service struct {
	url  string
	name string
}

func New(serviceName string, url string) *Service {

	return &Service{
		url:  url,
		name: serviceName,
	}
}

func (this *Service) Start() {
	http.HandleFunc("/channels", this.chanHandler)
	go func() {
		if err := this.start(); err != nil {
			panic(err)
		}
	}()
}

func (this *Service) start() error {
	return http.ListenAndServe(this.url, nil)
}

func (this *Service) chanHandler(w http.ResponseWriter, r *http.Request) {
	chStats := GetAll()

	resp := &ServiceChannelsStatus{
		Service:  this.name,
		Channels: chStats,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Printf("Error: %#v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonResp)
}

type ServiceChannelsStatus struct {
	Service  string                `json:"service"`
	Channels map[string]*ChanState `json:"channels"`
}

type Config struct {
	Name string
	Url  string
}
