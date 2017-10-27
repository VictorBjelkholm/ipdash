package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var myClient = &http.Client{}

// Stats is the collection of statistics that will be sent to people dashboards
type Stats struct {
	ID            string `json:"ID"`
	Name          string
	NumberOfPeers int

	TimeReceived time.Time `json:"time_received"`

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

type Message struct {
	From string `json:"From"`
	Data string `json:"Data"`
}

type CollectedStats struct {
	Stats []Stats
}

func createMessage(src []byte, target *Message) error {
	err := json.NewDecoder(bytes.NewReader(src)).Decode(target)
	if err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(target.Data)
	if err != nil {
		return err
	}
	target.Data = string(data)
	return nil
}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func saveStatsToDisk(path string, stats *CollectedStats) error {
	statsJson, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, statsJson, 0644)
	return err
}

func loadStatsFromDisk(path string, stats *CollectedStats) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.NewDecoder(bytes.NewReader(file)).Decode(stats)
	return err
}

func statsExistsOnDisk(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	collectedStats := &CollectedStats{}
	statePath := "./stats.json"

	// if cache exists, load
	if statsExistsOnDisk(statePath) {
		loadStatsFromDisk(statePath, collectedStats)
	}

	res, err := myClient.Get("http://localhost:5001/api/v0/pubsub/sub?arg=ipfs-dashboard")
	if err != nil {
		log.Fatal("Could not connect to your ipfs daemon at localhost:5001 - Make sure it's running and accept API requests!")
	}
	fmt.Println("Reading body")
	reader := bufio.NewReader(res.Body)
	go (func() {
		for {
			fmt.Println("Saving state")
			saveStatsToDisk(statePath, collectedStats)
			time.Sleep(1 * time.Second)
		}
	})()
	if res.Status != "200 OK" {
		fmt.Println(res.Status)
		if b, err := ioutil.ReadAll(res.Body); err == nil {
			log.Fatal(string(b))
		}
	}
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		msg := &Message{}
		createMessage(line, msg)
		// skipping the first message or if data is empty
		if msg.Data != "" {
			stats := &Stats{}
			err = json.NewDecoder(bytes.NewReader([]byte(msg.Data))).Decode(stats)
			// TODO make sure pubsub message is from correct node ID
			stats.TimeReceived = time.Now()
			if err != nil {
				panic(err)
			}
			collectedStats.Stats = append(collectedStats.Stats, *stats)
			fmt.Println("Collected " + strconv.Itoa(len(collectedStats.Stats)) + " stats so far")
		}
	}
	fmt.Println("After")
}
