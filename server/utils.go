package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vanshavenger/goproxynginx/utils"
)

// FindMatchingRule finds a rule that matches the given URL
func (w *Worker) FindMatchingRule(url string) *utils.Rule {
	for _, rule := range w.config.Server.Rules {
		if MatchRule(rule.Path, url) {
			return &rule
		}
	}
	return nil
}

// GetNextUpstream gets the next upstream server for a given rule
func (w *Worker) GetNextUpstream(rule *utils.Rule) string {
	upstreamID := rule.Upstreams[w.upstreamIndex%len(rule.Upstreams)]
	w.upstreamIndex++
	return upstreamID
}

// FindUpstream finds an upstream server by ID
func (w *Worker) FindUpstream(id string) *utils.Upstream {
	for _, upstream := range w.config.Server.Upstreams {
		if upstream.ID == id {
			return &upstream
		}
	}
	return nil
}

// ForwardRequest forwards a request to an upstream server
func (w *Worker) ForwardRequest(upstream *utils.Upstream, url string) WorkerMessageReply {
	protocol := upstream.Protocol
	if protocol == "" {
		protocol = "http"
	}
	fullURL := fmt.Sprintf("%s://%s%s", protocol, upstream.URL, url)

	resp, err := http.Get(fullURL)
	if err != nil {
		return WorkerMessageReply{
			Error:      fmt.Sprintf("Error forwarding request: %v", err),
			StatusCode: 500,
		}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return WorkerMessageReply{
			Error:      fmt.Sprintf("Error reading response body: %v", err),
			StatusCode: 500,
		}
	}

	return WorkerMessageReply{
		Data:       string(data),
		StatusCode: resp.StatusCode,
	}
}

// MatchRule checks if a URL matches a rule path
func MatchRule(rulePath, url string) bool {
	return strings.HasPrefix(url, rulePath)
}
