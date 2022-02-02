package rpc

import (
	gtsasamiserver "github.com/Tim-vo/gt-sasami-server/gt-sasami-server"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Server is the API web server
type Server struct {
	logger  *zap.SugaredLogger
	router  chi.Router
	grStore gtsasamiserver.GRStore
}

func Setup(router chi.Router, grStore gtsasamiserver.GRStore) error {

	s := &Server{
		logger:  zap.S().With("package", "thingrpc"),
		router:  router,
		grStore: grStore,
	}

	// Base Functions
	s.router.Route("/api", func(r chi.Router) {

	})

	return nil

}
