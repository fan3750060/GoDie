package route

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"goframe/app/controller"
	"goframe/app/middleware"
	"net/http"
)

type Routes interface {
	LoadRoute() http.Handler
}

type RoutesClass struct {
}

func (routesClass RoutesClass) LoadRoute() http.Handler {
	errorChain := alice.New(middleware.Process)
	r := mux.NewRouter()
	r = routelist(r)
	return errorChain.Then(r)
}

func routelist(r *mux.Router) *mux.Router {

	//路由组
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/index", controller.Index).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	return r
}
