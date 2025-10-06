package routes

import (
	_ "Effective_Mobile/docs"
	"Effective_Mobile/internal/httpserver/handlers/del"
	"Effective_Mobile/internal/httpserver/handlers/get"
	"Effective_Mobile/internal/httpserver/handlers/post"
	"Effective_Mobile/internal/httpserver/handlers/put"
	log "Effective_Mobile/internal/httpserver/middleware/logger"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/storage/pg"
	"errors"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"strings"
)

var MethodNotAllowed = errors.New("method not allowed")

type Router struct {
	getHandler     http.HandlerFunc
	postHandler    http.HandlerFunc
	putHandler     http.HandlerFunc
	deleteHandler  http.HandlerFunc
	swaggerHandler http.Handler
}

func New(storage *pg.Storage) *Router {
	return &Router{
		getHandler:     log.Middleware(get.New(storage)),
		postHandler:    log.Middleware(post.New(storage)),
		putHandler:     log.Middleware(put.New(storage, storage)),
		deleteHandler:  log.Middleware(del.New(storage)),
		swaggerHandler: httpSwagger.WrapHandler,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const op = "httpserver.routes.serveHTTP"

	if req.URL.Path == "/swagger" || strings.HasPrefix(req.URL.Path, "/swagger/") {
		r.swaggerHandler.ServeHTTP(w, req)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	logger.Info("%s: → %s %s | IP: %s | UA: %s", op, req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent())

	if path == "ping" || path == "health" {
		logger.Debug("%s: health check request", op)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	switch {
	case len(parts) == 1 && parts[0] == "people":
		logger.Debug("%s: matched route /people", op)
		r.handlePeople(w, req)
	case len(parts) == 2 && parts[0] == "people":
		logger.Debug("%s: matched route /people/{id}", op)
		r.handlePeopleWithID(w, req)
	default:
		logger.Error("%s: unknown path %q", op, req.URL.Path)
		http.Error(w, fmt.Sprintf("%s: %s", op, MethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (r *Router) handlePeople(w http.ResponseWriter, req *http.Request) {
	const op = "httpserver.routes.handlePeople"

	switch req.Method {
	case http.MethodGet:
		logger.Debug("%s: GET /people", op)
		r.getHandler(w, req)
	case http.MethodPost:
		logger.Debug("%s: POST /people", op)
		r.postHandler(w, req)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		logger.Error("%s: method %s not allowed", op, req.Method)
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, fmt.Sprintf("%s: %s", op, MethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (r *Router) handlePeopleWithID(w http.ResponseWriter, req *http.Request) {
	const op = "httpserver.routes.handlePeopleWithID"

	switch req.Method {
	case http.MethodPut:
		logger.Debug("%s: PUT /people/{id}", op)
		r.putHandler(w, req)
	case http.MethodDelete:
		logger.Debug("%s: DELETE /people/{id}", op)
		r.deleteHandler(w, req)
	case http.MethodGet:
		logger.Debug("%s: GET /people", op)
		r.getHandler(w, req)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		logger.Error("%s: method %s not allowed", op, req.Method)
		w.Header().Set("Allow", "PUT, DELETE")
		http.Error(w, fmt.Sprintf("%s: %s", op, MethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

//func RegisterSwagger() {
//	http.Handle("/swagger/", httpSwagger.WrapHandler)
//}
