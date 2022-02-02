package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/knadh/koanf"
	"github.com/snowzach/certtools"
	"github.com/snowzach/certtools/autocert"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	router chi.Router
	server *http.Server
}

//Setup API listener
func New(config *koanf.Koanf) (*Server, error) {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	if config.Bool("server.log_requests") {
		switch config.String("logger.encoding") {
		case "strackdriver":
			router.Use(loggerHTTPMiddlewareStackdriver(config.Bool("server.log_requests_body"), config.Strings(("server.log_disabled_http"))))
		default:
			router.Use(loggerHTTPMiddlewareDefault(config.Bool("server.log_requests_body"), config.Strings("server.log_disabled_http")))
		}
	}

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   config.Strings("server.cors.allowed_origins"),
		AllowedMethods:   config.Strings("server.cors.allowed_methods"),
		AllowedHeaders:   config.Strings("server.cors.allowed_headers"),
		AllowCredentials: config.Bool("server.cors.allowed_credentials"),
		MaxAge:           config.Int("server.cors.max_age"),
	}).Handler)

	server := &Server{
		logger: zap.S().With("package", "server"),
		router: router,
	}

	return server, nil
}

//Listen and serve
func (server *Server) ListenAndServe(config *koanf.Koanf) error {

	server.server = &http.Server{
		Addr:    net.JoinHostPort(config.String("server.host"), config.String("server.port")),
		Handler: server.router,
	}

	listener, err := net.Listen("tcp", server.server.Addr)

	if err != nil {
		return fmt.Errorf("could not listen on %s: %v", server.server.Addr, err)
	}

	// Enable TLS when turned on
	if config.Bool(("server.tls")) {
		var cert tls.Certificate

		if config.Bool("server.devcert") {
			server.logger.Warn("WARNING: This server is using an insecure development TLS certificate. This should only be used for development!")
			cert, err = autocert.New(autocert.InsecureStringReader("localhost"))

			if err != nil {
				return fmt.Errorf("could not autocert generate server certificate: %v", err)
			}
		} else {
			cert, err = tls.LoadX509KeyPair(config.String("server.certfile"), config.String("server.keyfile"))
			if err != nil {
				return fmt.Errorf("could not load server certificate: %v", err)
			}
		}

		server.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   certtools.SecureTLSMinVersion(),
			CipherSuites: certtools.SecureTLSCipherSuites(),
		}

		listener = tls.NewListener(listener, server.server.TLSConfig)
	}

	go func() {
		if err = server.server.Serve(listener); err != nil {
			server.logger.Fatalw("API Listener error", "error", err, "address", server.server.Addr)
		}
	}()

	server.logger.Infow("API Listening", "address", server.server.Addr, "tls", config.Bool("server.tls"))

	//Enable profiler
	if config.Bool("server.profiler_enabled") && config.String("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", config.String("server.profiler_path"))
		server.router.Mount(config.String("server.profiler_path"), middleware.Profiler())
	}

	return nil
}

func (server *Server) Router() chi.Router {
	return server.router
}
