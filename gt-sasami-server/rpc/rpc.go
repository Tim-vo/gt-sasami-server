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
	grStore gtsasamiserver.GTStore
}

func Setup(router chi.Router, grStore gtsasamiserver.GTStore) error {

	server := &Server{
		logger:  zap.S().With("package", "accountrpc"),
		router:  router,
		grStore: grStore,
	}

	// Base Functions
	server.router.Route("/api", func(router chi.Router) {
		router.Post("/accounts", server.AccountSave())
		router.Get("/accounts/{id}", server.AccountGetByID())
		router.Delete("/accounts/{id}", server.AccountDeleteByID())
		router.Get("/accounts", server.AccountsFind())
	})

	return nil

}
