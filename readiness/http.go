package readiness

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPReadiness struct {
	httpClient *http.Client
	host       string
	port       uint
	path       string
	interval   uint
}

var defaultHTTPReadiness = HTTPReadiness{
	host:     "localhost",
	port:     80,
	path:     "/",
	interval: 10,
}

type HTTPReadinessOption func(*HTTPReadiness)

func HTTPClient(client *http.Client) HTTPReadinessOption {
	return func(readiness *HTTPReadiness) {
		readiness.httpClient = client
	}
}

func HTTPHost(host string, port uint) HTTPReadinessOption {
	return func(readiness *HTTPReadiness) {
		readiness.host = host
		readiness.port = port
	}
}

func HTTPPath(path string) HTTPReadinessOption {
	return func(readiness *HTTPReadiness) {
		readiness.path = path
	}
}

func HTTPInterval(interval uint) HTTPReadinessOption {
	return func(readiness *HTTPReadiness) {
		readiness.interval = interval
	}
}

func NewHTTPReadiness(opts ...HTTPReadinessOption) *HTTPReadiness {
	readiness := defaultHTTPReadiness

	for _, opt := range opts {
		opt(&readiness)
	}

	if readiness.httpClient == nil {
		readiness.httpClient = &http.Client{}
	}

	return &readiness
}

func (h *HTTPReadiness) Host() string {
	return h.host
}

func (h *HTTPReadiness) Port() uint {
	return h.port
}

func (h *HTTPReadiness) Path() string {
	return h.path
}

func (h *HTTPReadiness) Interval() uint {
	return h.interval
}

func (h *HTTPReadiness) IsReady() bool {
	url := fmt.Sprintf("%s:%d%s", h.Host(), h.Port(), h.Path())

	resp, err := h.httpClient.Get(url)


	if err != nil {
		log.Printf("failed to get readiness(Ignore), %s\n", err.Error())
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("service response without health status", resp)
		return false
	}

	return true
}
