package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/savioruz/bake/internal/builder"
	e "github.com/savioruz/bake/pkg/error"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	port string
	log  *logrus.Logger
	mux  *http.ServeMux
}

func NewServer(viper *viper.Viper, log *logrus.Logger) *Server {
	return &Server{
		port: viper.GetString("APP_PORT"),
		log:  log,
		mux:  http.NewServeMux(),
	}
}

// LoggingMiddleware wraps an http.HandlerFunc and logs request details
func (s *Server) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next(rw, r)

		// Log response details
		duration := time.Since(start)
		s.log.WithFields(logrus.Fields{
			"method":      r.Method + " " + r.URL.Path,
			"duration":    duration,
			"duration_ms": float64(duration.Milliseconds()),
		}).Info("Request completed")
	}
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) chainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (s *Server) RegisterRoutes(routes []builder.Routes) {
	// Create API router for /api/v1 group
	apiRouter := &Router{
		prefix: "/api/v1",
		routes: make(map[string]http.HandlerFunc),
		log:    s.log,
	}

	// Register API routes
	for _, route := range routes {
		// Skip Swagger routes
		if strings.Contains(route.Path, "/docs") {
			continue
		}

		path := strings.TrimPrefix(route.Path, "/api/v1")

		handler := s.chainMiddleware(
			route.Handler,
			s.LoggingMiddleware,
		)

		apiRouter.Handle(route.Method, path, handler)
	}

	// Register Swagger routes directly
	for _, route := range routes {
		if strings.Contains(route.Path, "/docs") {
			s.mux.HandleFunc(route.Path, route.Handler)
		}
	}

	// Register API handler
	s.mux.HandleFunc("/api/v1/", apiRouter.ServeHTTP)
}

// Router handles /api/v1 routes
type Router struct {
	prefix string
	routes map[string]http.HandlerFunc
	log    *logrus.Logger
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	key := method + ":" + path
	r.routes[key] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Remove prefix from path
	path := strings.TrimPrefix(req.URL.Path, r.prefix)
	key := req.Method + ":" + path

	r.log.WithFields(logrus.Fields{
		"method": req.Method,
		"path":   path,
		"key":    key,
	}).Debug("Looking for route")

	// Find handler
	if handler, ok := r.routes[key]; ok {
		handler(w, req)
		return
	}

	// No route found
	e.ErrorHandler(w, req, http.StatusNotFound, e.ErrNotFound)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	s.log.Info("Server starting on port ", s.port)

	return http.ListenAndServe(addr, s.mux)
}
