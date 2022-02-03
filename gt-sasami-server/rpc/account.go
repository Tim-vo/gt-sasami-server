package rpc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/snowzach/queryp"

	gtSasamiServer "github.com/Tim-vo/gt-sasami-server/gt-sasami-server"
	"github.com/Tim-vo/gt-sasami-server/server"
	"github.com/Tim-vo/gt-sasami-server/store"
)

// AccountSave saves an account
//
// @ID AccountSave
// @Tags accounts
// @Summary Save account
// @Description Save a account
// @Param account body gtSasamiServer.Account true "Account"
// @Success 200 {object} gtSasamiServer.Account
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /accounts [post]
func (s *Server) AccountSave() http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()

		var account = new(gtSasamiServer.Account)
		if err := server.DecodeJSON(request.Body, account); err != nil {
			server.RenderErrInvalidRequest(writer, err)
			return
		}

		err := s.grStore.AccountSave(ctx, account)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(writer, serr.ErrorForOp(store.ErrorOpSave))
			} else {
				errID := server.RenderErrInternalWithID(writer, nil)
				s.logger.Errorw("AccountSave error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(writer, http.StatusOK, account)
	}

}

// AccountGetByID saves a account
//
// @ID AccountGetByID
// @Tags accounts
// @Summary Get account
// @Description Get a account
// @Param id path string true "ID"
// @Success 200 {object} gtSasamiServer.Account
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 404 {object} server.ErrResponse "Not Found"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /accounts/{id} [get]
func (s *Server) AccountGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		account, err := s.grStore.AccountGetByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "account")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpGet))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("AccountGetByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, account)
	}

}

// AccountDeleteByID saves a account
//
// @ID AccountDeleteByID
// @Tags accounts
// @Summary Delete account
// @Description Delete a account
// @Param id path string true "ID"
// @Success 204 "Success"
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 404 {object} server.ErrResponse "Not Found"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /accounts/{id} [delete]
func (s *Server) AccountDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.grStore.AccountDeleteByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "account")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpDelete))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("AccountDeleteByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderNoContent(w)

	}

}

// AccountsFind saves a account
//
// @ID AccountsFind
// @Tags accounts
// @Summary Find accounts
// @Description Find accounts
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param offset query int false "offset"
// @Param limit query int false "limit"
// @Param sort query string false "query"
// @Success 200 {array} gtSasamiServer.Account
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /accounts [get]
func (s *Server) AccountsFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		qp, err := queryp.ParseRawQuery(r.URL.RawQuery)
		if err != nil {
			server.RenderErrInvalidRequest(w, err)
		}

		accounts, count, err := s.grStore.AccountsFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpFind))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("AccountsFind error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, store.Results{Count: count, Results: accounts})

	}

}
