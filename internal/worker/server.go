package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// HTTPServer exposes /healthz and /metrics endpoints.
type HTTPServer struct {
	srv *http.Server
	log *logrus.Logger
}

// NewHTTPServer constructs the HTTP server on the given port.
func NewHTTPServer(port string, log *logrus.Logger) *HTTPServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{srv: srv, log: log}
}

func (s *HTTPServer) Start() {
	go func() {
		s.log.WithField("addr", s.srv.Addr).Info("HTTP server started")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.WithError(err).Error("HTTP server error")
		}
	}()
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.WithError(err).Error("HTTP server shutdown error")
	}
}
