package auth
import (
	"sync"
	"time"
	
)
func NewAuth() (*Router, *MiddleWare, *APIMiddleWare) {
	sm := &tokenMap{
		tuMap: make(map[string]*userSession),
	}
	return &Router{tmap: sm}, &MiddleWare{tmap: sm}, &APIMiddleWare{nameMap: map[string]model.User{}}
}

// token to model.User map
type tokenMap struct {
	lock  sync.RWMutex
	tuMap map[string]*userSession

	cleanTime time.Time
}