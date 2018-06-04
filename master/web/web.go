package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	log "github.com/Sirupsen/logrus"
)

type Server struct {
	routers []Router
}

type Route struct {
	Handler APIFunc

	Method string
	Path   string
	Prefix string
}
type Router interface {
	Routes() []Route
}
type Component interface {
	Router
	Setup(orm *gorm.DB, middlewares map[string]MiddlewareFunc)
}

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error

type MiddlewareFunc func(handler APIFunc) APIFunc

type Middleware interface {
	Setup(orm *gorm.DB)
	AsFunc(handler APIFunc) APIFunc
}

func NewServer() *Server {
	return &Server{
		routers: []Router{},
	}
}

func (s *Server) makeHTTPHandler(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("calling %s %s", r.Method, r.URL.Path)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := context.Background()
		if err := handler(ctx, w, r, mux.Vars(r)); err != nil {
			s.handlerError(err, w, r)
		}

	}
}

func errorResponse(err error) interface{} {
	type resp struct {
		Error struct {
			Message string
		} `json:"error"`
	}

	r := resp{}
	r.Error.Message = err.Error()
	return r
}

func WriteJSONResponse(w http.ResponseWriter, val interface{}) error {
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		log.Errorf("error in write json Response: %s", err)
	}
	return err
}

func (s *Server) handlerError(err error, w http.ResponseWriter, r *http.Request) {
	log.Errorf("handle for %s %s returned error: %s", r.Method, r.URL.Path, err)
	switch err {
	case gorm.RecordNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	WriteJSONResponse(w, errorResponse(err))
}

func (s *Server) CreateHandler() http.Handler {
	m := mux.NewRouter()

	for _, router := range s.routers {
		for _, route := range router.Routes() {
			handler := s.makeHTTPHandler(route.Handler)
			if route.Prefix != "" {
				log.Debugf("register prefix %s %s", route.Method, route.Prefix)
				m.PathPrefix(route.Prefix).Methods(route.Method).Handler(handler)
			} else {
				log.Debugf("register path %s %s", route.Method, route.Path)
				m.Path(route.Path).Methods(route.Method).Handler(handler)
			}
		}
	}

	return m
}

func (s *Server) RegisterRouter(router Router) {
	s.routers = append(s.routers, router)
}

func (s *Server) Run() {
	mux := s.CreateHandler()
	http.Handle("/", mux)
	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		panic(err)
	}
}
