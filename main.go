package main

import (
	"encoding/json"
	"net/http"
	"time"
	"runtime"
	"fmt"
)

const INTERVAL_PERIOD time.Duration = 3 * time.Hour

//const INTERVAL_PERIOD time.Duration = 30 * time.Second
var repository *RemoteJobRepository

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))

	parser := NewRocketJobHtmlParser()

	scheduler := NewScheduler(INTERVAL_PERIOD)
	scheduler.Runner = func() {
		parser.TotalGet()
	}
	scheduler.MinTick = 0
	scheduler.SecTick = 0
	scheduler.NsecTick = 0
	scheduler.Run()

	repository = NewRemoteJobRepository()
	repository.Open()
	defer repository.Close()
	http.HandleFunc("/api/getRemoteJobInfo", getRemoteJobInfo)
	http.HandleFunc("/api/batchNow", func(w http.ResponseWriter, r *http.Request) {
		parser.TotalGet()
		apiResult := new(APIResult)
		apiResult.ResultCode = 200
		returnVal, _ := json.Marshal(apiResult)
		w.Write(returnVal)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/views/", http.StatusFound)
	})
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("./views/"))))
	http.ListenAndServe(":4001", nil)
}

func getRemoteJobInfo(w http.ResponseWriter, r *http.Request) {
	apiResult := new(APIResult)
	apiResult.ResultCode = 200
	apiResult.ResultContent = repository.FindAll()
	returnVal, _ := json.Marshal(apiResult)
	w.Write(returnVal)
}

type APIResult struct {
	ResultCode    int
	ResultContent interface{}
}
