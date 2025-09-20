package httprequests

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/PingTower/internal/entities"
)

type HttpMetrics struct {
	URL            string  `json:"url"`
	Status         string  `json:"status"`
	Error          string  `json:"error,omitempty"`
	StatusCode     int     `json:"status_code,omitempty"`
	DNSMs          float64 `json:"dns_ms"`
	ConnectMs      float64 `json:"connect_ms"`
	TLSMs          float64 `json:"tls_ms"`
	TTFBMs         float64 `json:"ttfb_ms"`
	TotalMs        float64 `json:"total_ms"`
	ResponseLength int64   `json:"response_length,omitempty"`
}

// Функция для проверки конкретных url'ов из джобов
// привязывается к полю Handler
func HttpCheck(job entities.Job) (HttpMetrics, error) {
	var (
		dnsStart, dnsDone      time.Time
		connStart, connDone    time.Time
		tlsStart, tlsDone      time.Time
		gotConn, firstResponse time.Time
	)

	req, _ := http.NewRequest("GET", job.URL, nil)

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },

		ConnectStart: func(_, _ string) { connStart = time.Now() },
		ConnectDone:  func(_, _ string, _ error) { connDone = time.Now() },

		TLSHandshakeStart: func() { tlsStart = time.Now() },
		TLSHandshakeDone:  func(_ tls.ConnectionState, _ error) { tlsDone = time.Now() },

		GotConn:              func(_ httptrace.GotConnInfo) { gotConn = time.Now() },
		GotFirstResponseByte: func() { firstResponse = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start := time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	total := time.Since(start)

	metrics := HttpMetrics{
		URL:     job.URL,
		TotalMs: float64(total.Milliseconds()),
	}

	if !dnsStart.IsZero() && !dnsDone.IsZero() {
		metrics.DNSMs = dnsDone.Sub(dnsStart).Seconds() * 1000
	}
	if !connStart.IsZero() && !connDone.IsZero() {
		metrics.ConnectMs = connDone.Sub(connStart).Seconds() * 1000
	}
	if !tlsStart.IsZero() && !tlsDone.IsZero() {
		metrics.TLSMs = tlsDone.Sub(tlsStart).Seconds() * 1000
	}
	if !gotConn.IsZero() && !firstResponse.IsZero() {
		metrics.TTFBMs = firstResponse.Sub(gotConn).Seconds() * 1000
	}

	if err != nil {
		metrics.Status = "fail"
		metrics.Error = err.Error()
		return metrics, err
	}
	defer resp.Body.Close()

	metrics.Status = "ok"
	metrics.StatusCode = resp.StatusCode

	if resp.ContentLength > 0 {
		metrics.ResponseLength = resp.ContentLength
	}

	return metrics, nil
}
