package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/vanshavenger/goproxynginx/utils"
)

// CreateServer creates a new server
func CreateServer(config *utils.RootConfig) {
	workersCount := *(config).Server.Workers
	port := config.Server.Listen
	if workersCount == 0 {
		workersCount = runtime.NumCPU()
	}
	workerPool := make(chan *Worker, workersCount)
    defer close(workerPool)
	for i := 0; i < workersCount; i++ {
		worker := newWorker(config)
		workerPool <- worker
		log.Printf("Worker %d started", i)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		worker := <-workerPool
		defer func() { workerPool <- worker }()

		payload := WorkerMessage{
			RequestType: "HTTP",
			Body:        r.Body,
			URL:         r.URL.Path,
		}

		reply := worker.processMessage(payload)

		w.WriteHeader(reply.StatusCode)
		if reply.Error != "" {
			fmt.Fprint(w, reply.Error)
		} else if reply.Data != "" {
			fmt.Fprint(w, reply.Data)
		} else {
			fmt.Fprint(w, "No content")
		}
	})

	log.Printf("Server listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

func newWorker(config *utils.RootConfig) *Worker {
	return &Worker{
		config:        config,
		upstreamIndex: 0,
	}
}

func (w *Worker) processMessage(message WorkerMessage) WorkerMessageReply {
	w.mu.Lock()
	defer w.mu.Unlock()

	rule := w.FindMatchingRule(message.URL)
	if rule == nil {
		return WorkerMessageReply{
			Error:      "Rule not found",
			StatusCode: 404,
		}
	}

	upstreamID := w.GetNextUpstream(rule)
	if upstreamID == "" {
		return WorkerMessageReply{
			Error:      "Upstream not found",
			StatusCode: 500,
		}
	}

	upstream := w.FindUpstream(upstreamID)
	if upstream == nil {
		return WorkerMessageReply{
			Error:      "Upstream not found",
			StatusCode: 500,
		}
	}

	return w.ForwardRequest(upstream, message.URL)
}
