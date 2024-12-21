package server

import (
	"sync"

	"github.com/vanshavenger/goproxynginx/utils"
)

// WorkerMessage represents a message sent to a worker
type WorkerMessage struct {
	RequestType string            `json:"requestType"`
	Headers     map[string]string `json:"headers"`
	Body        interface{}       `json:"body"`
	URL         string            `json:"url"`
}

// WorkerMessageReply represents a message sent by a worker
type WorkerMessageReply struct {
	Data       string `json:"data,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode"`
}

// Worker represents a worker
type Worker struct {
	config        *utils.RootConfig
	upstreamIndex uint32
	mu            sync.Mutex
}
