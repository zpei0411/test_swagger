package web

import (
	"github.com/jinzhu/gorm"
)

type SwaRouter struct {
	orm *gorm.DB
}

func (r *SwaRouter) Setup(orm *gorm.DB, ms map[string]MiddlewareFunc) {
	r.orm = orm
}

func (r *SwaRouter) Routes() []Route {
	routes := []Route{
		{
			Handler: r.CreateUser,
			Method:  "POST",
			Path:    "/user",
		},
		{
			Handler: r.GetUser,
			Method:  "GET",
			Path:    "/users/{id}",
		},
	}

	for i, route := range routes {
		routes[i].Handler = route.Handler
	}
	return routes
}
