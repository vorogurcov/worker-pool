package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
	"traceroute-optimised/internal"
	"traceroute-optimised/internal/domain"
	"traceroute-optimised/internal/metrics"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	workerCnt  = 3
	jobQueue   chan domain.Job
	queueSize  = 30
	jobCounter int64
)

func main() {
	m := metrics.NewMetrics()

	jobQueue = make(chan domain.Job, queueSize)

	for i := 0; i < workerCnt; i++ {
		go internal.Worker(i, jobQueue)
	}

	if err := godotenv.Load("./cmd/.env"); err != nil {
		log.Println("No .env file found, using system env")
	}

	dataRepository, err := internal.NewDataRepository()
	if err != nil {
		log.Fatal(err)
	}

	dataService := internal.DataService{DataRepo: dataRepository}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/get", func(w http.ResponseWriter, r *http.Request) {
		enqueueJob(w, r, internal.NewHandleGetData(&dataService, r), m)
	})

	mux.HandleFunc("/v1/bulk", func(w http.ResponseWriter, r *http.Request) {
		enqueueJob(w, r, internal.NewHandleGetBulkData(&dataService, r), m)
	})

	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	fmt.Println("Listening on :30099")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func enqueueJob(w http.ResponseWriter, r *http.Request, exec domain.JobExecFunc, m *metrics.Metrics) {
	jobID := atomic.AddInt64(&jobCounter, 1)
	ctx := r.Context()

	job := domain.Job{
		ID:   jobID,
		Ctx:  ctx,
		Resp: make(chan domain.Result, 1),
		Exec: exec,
	}

	select {
	case jobQueue <- job:
		// задача принята
	default:
		http.Error(w, "server overloaded, try later", http.StatusServiceUnavailable)
		return
	}

	select {
	case res := <-job.Resp:
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
		m.IncRequest(res.Path, strconv.Itoa(res.Status))

		writeJSON(w, res.Status, res.Body)
	case <-ctx.Done():
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
	case <-time.After(3 * time.Second):
		http.Error(w, "processing timeout", http.StatusGatewayTimeout)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
