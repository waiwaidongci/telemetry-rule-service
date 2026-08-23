package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	protocolTextDecoder "github.com/example/telemetry-rule-service/internal/adapter/protocol"
	"github.com/example/telemetry-rule-service/internal/application/events"
	"github.com/example/telemetry-rule-service/internal/application/ingest"
	"github.com/example/telemetry-rule-service/internal/application/rules"
	"github.com/example/telemetry-rule-service/internal/application/subscriptions"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"github.com/example/telemetry-rule-service/internal/domain/source"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

type Server struct {
	sources       source.Repository
	metrics       metric.Repository
	rules         rule.Repository
	ingest        *ingest.Service
	evaluate      *rules.Service
	events        *events.Service
	subscriptions *subscriptions.Service
	logger        *slog.Logger
	requests      atomic.Uint64
}

func NewServer(src source.Repository, met metric.Repository, rs rule.Repository, es *events.Service, ss *subscriptions.Service, ig *ingest.Service, ev *rules.Service, l *slog.Logger) *Server {
	return &Server{sources: src, metrics: met, rules: rs, events: es, subscriptions: ss, ingest: ig, evaluate: ev, logger: l}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { write(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) { write(w, 200, map[string]string{"status": "ready"}) })
	mux.HandleFunc("GET /metrics", s.metricsHandler)
	mux.HandleFunc("/api/v1/sources", s.sourcesHandler)
	mux.HandleFunc("/api/v1/metrics", s.metricsHandlerAPI)
	mux.HandleFunc("/api/v1/rules", s.rulesHandler)
	mux.HandleFunc("/api/v1/rules/", s.rulesHandler)
	mux.HandleFunc("/api/v1/telemetry", s.telemetryHandler)
	mux.HandleFunc("/api/v1/events", s.eventsHandler)
	mux.HandleFunc("/api/v1/events/", s.eventsHandler)
	mux.HandleFunc("/api/v1/subscriptions", s.subscriptionsHandler)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.requests.Add(1)
		w.Header().Set("X-Request-ID", fmt.Sprintf("req-%d", s.requests.Load()))
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
		s.logger.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "telemetry_http_requests_total %d\n", s.requests.Load())
}
func (s *Server) sourcesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		v, e := s.sources.List(r.Context())
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
	case "POST":
		var v source.DataSource
		if !decode(r, &v) {
			fail(w, fmt.Errorf("invalid source"))
			return
		}
		v.CreatedAt = time.Now().UTC()
		if e := s.sources.Create(r.Context(), v); e != nil {
			fail(w, e)
			return
		}
		write(w, 201, v)
	default:
		http.NotFound(w, r)
	}
}
func (s *Server) metricsHandlerAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		v, e := s.metrics.List(r.Context())
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
	case "POST":
		var v metric.Definition
		if !decode(r, &v) {
			fail(w, fmt.Errorf("invalid metric"))
			return
		}
		v.CreatedAt = time.Now().UTC()
		if e := s.metrics.Create(r.Context(), v); e != nil {
			fail(w, e)
			return
		}
		write(w, 201, v)
	default:
		http.NotFound(w, r)
	}
}
func (s *Server) rulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := s.rules.List(r.Context())
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	if r.Method == "POST" {
		var v rule.Definition
		if !decode(r, &v) {
			fail(w, fmt.Errorf("invalid rule"))
			return
		}
		v.CreatedAt = time.Now().UTC()
		v.Enabled = true
		if e := s.rules.Create(r.Context(), v); e != nil {
			fail(w, e)
			return
		}
		write(w, 201, v)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) >= 4 && r.Method == "PATCH" {
		id := parts[3]
		v, e := s.rules.Get(r.Context(), id)
		if e != nil {
			fail(w, e)
			return
		}
		var in struct {
			Enabled *bool `json:"enabled"`
		}
		if !decode(r, &in) {
			fail(w, fmt.Errorf("invalid rule"))
			return
		}
		if in.Enabled != nil {
			v.Enabled = *in.Enabled
		}
		if e = s.rules.Update(r.Context(), v); e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	http.NotFound(w, r)
}
func (s *Server) telemetryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	if r.ContentLength > 2<<20 {
		fail(w, fmt.Errorf("payload too large"))
		return
	}
	body, e := io.ReadAll(io.LimitReader(r.Body, 2<<20+1))
	if e != nil {
		fail(w, e)
		return
	}
	if len(body) > 2<<20 {
		fail(w, fmt.Errorf("payload too large"))
		return
	}
	var v metric.Sample
	if strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		v, e = protocolTextDecoder.TextDecoder{}.Decode(body)
	} else {
		v, e = protocolTextDecoder.JSONDecoder{}.Decode(body)
	}
	if e != nil {
		fail(w, e)
		return
	}
	if e = s.ingestSample(r.Context(), v); e != nil {
		fail(w, e)
		return
	}
	events, err := s.evaluate.Evaluate(r.Context(), v)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, 202, map[string]any{"sample": v, "events": events})
}
func protocolJSON(b []byte) (metric.Sample, error) {
	var v metric.Sample
	if e := json.Unmarshal(b, &v); e != nil {
		return v, e
	}
	if v.Timestamp.IsZero() {
		v.Timestamp = time.Now().UTC()
	}
	return v, nil
}
func protocolText(b []byte) (metric.Sample, error) {
	p := strings.Split(strings.TrimSpace(string(b)), ",")
	if len(p) < 4 {
		return metric.Sample{}, fmt.Errorf("expected source,metric,value,timestamp")
	}
	v := metric.Sample{SourceID: p[0], MetricID: p[1]}
	if _, e := fmt.Sscanf(p[2], "%f", &v.Value); e != nil {
		return v, e
	}
	t, e := time.Parse(time.RFC3339, p[3])
	v.Timestamp = t
	return v, e
}
func (s *Server) ingestSample(ctx context.Context, v metric.Sample) error {
	if e := v.Validate(); e != nil {
		return e
	}
	return s.metrics.Record(ctx, v)
}
func (s *Server) eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		limit := 50
		if x, _ := strconv.Atoi(r.URL.Query().Get("limit")); x > 0 {
			limit = x
		}
		v, e := s.events.List(r.Context(), event.Status(r.URL.Query().Get("status")), limit)
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) >= 4 && r.Method == "PATCH" {
		id := parts[3]
		var in struct {
			Status event.Status `json:"status"`
		}
		if !decode(r, &in) {
			fail(w, fmt.Errorf("invalid event"))
			return
		}
		v, e := s.events.ChangeStatus(r.Context(), id, in.Status)
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	http.NotFound(w, r)
}
func (s *Server) subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := s.subscriptions.List(r.Context())
		if e != nil {
			fail(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	if r.Method == "POST" {
		var v subscription.Subscription
		if !decode(r, &v) {
			fail(w, fmt.Errorf("invalid subscription"))
			return
		}
		v.CreatedAt = time.Now().UTC()
		v.Enabled = true
		if e := s.subscriptions.Create(r.Context(), v); e != nil {
			fail(w, e)
			return
		}
		write(w, 201, v)
		return
	}
	http.NotFound(w, r)
}
func decode(r *http.Request, v any) bool { return json.NewDecoder(r.Body).Decode(v) == nil }
func write(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, e error) {
	write(w, http.StatusBadRequest, map[string]string{"error": e.Error()})
}
