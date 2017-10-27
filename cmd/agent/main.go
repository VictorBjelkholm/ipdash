package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

// Stats is the collection of statistics that will be sent to people dashboards
type Stats struct {
	ID            string `json:"ID"`
	Name          string
	NumberOfPeers int

	TotalIn  int     `json:"TotalIn"`
	TotalOut int     `json:"TotalOut"`
	RateIn   float64 `json:"RateIn"`
	RateOut  float64 `json:"RateOut"`

	Diskinfo struct {
		FreeSpace  int64  `json:"free_space"`
		Fstype     string `json:"fstype"`
		TotalSpace int64  `json:"total_space"`
	} `json:"diskinfo"`
	IpfsCommit  string `json:"ipfs_commit"`
	IpfsVersion string `json:"ipfs_version"`
	Memory      struct {
		Swap int `json:"swap"`
		Virt int `json:"virt"`
	} `json:"memory"`
	Net struct {
		Online bool `json:"online"`
	} `json:"net"`
	Runtime struct {
		Arch          string `json:"arch"`
		Compiler      string `json:"compiler"`
		Gomaxprocs    int    `json:"gomaxprocs"`
		Numcpu        int    `json:"numcpu"`
		Numgoroutines int    `json:"numgoroutines"`
		Os            string `json:"os"`
		Version       string `json:"version"`
	} `json:"runtime"`
}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func getBandwidth(stats *Stats) error {
	err := getJSON("http://localhost:5001/api/v0/stats/bw", stats)
	if err != nil {
		return err
	}
	return nil
}

func getSystem(stats *Stats) error {
	err := getJSON("http://localhost:5001/api/v0/diag/sys", stats)
	if err != nil {
		return err
	}
	return nil
}

func getID(stats *Stats) error {
	err := getJSON("http://localhost:5001/api/v0/id", stats)
	if err != nil {
		return err
	}
	return nil
}

func publishStats(stats *Stats) {
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	res, err := myClient.Get("http://localhost:5001/api/v0/pubsub/pub?arg=ipfs-dashboard&arg=" + string(statsJSON))
	if err != nil {
		panic(err)
	}
	if res.Status != "200 OK" {
		fmt.Println(res.Status)
		if b, err := ioutil.ReadAll(res.Body); err == nil {
			log.Fatal(string(b))
		}
	}
}

func getAndPublishStats() {
	fmt.Println("Getting stats...")
	stats := &Stats{}
	err := getID(stats)
	if err != nil {
		panic(err)
	}
	err = getBandwidth(stats)
	if err != nil {
		panic(err)
	}
	err = getSystem(stats)
	if err != nil {
		panic(err)
	}
	publishStats(stats)
	fmt.Println("Published stats!")
}

func main() {
	stats := &Stats{}
	err := getID(stats)
	if err != nil {
		log.Fatal("Could not connect to your ipfs daemon at localhost:5001 - Make sure it's running and accept API requests!")
	}
	for {
		getAndPublishStats()
		time.Sleep(1000 * time.Millisecond)
	}
}
