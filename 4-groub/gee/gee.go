package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	Engine struct {
		*RouterGroup // 指向嵌套的 RouterGroup，作为根路由
		router       *router
		groups       []*RouterGroup // store all groups
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (groub *RouterGroup) Group(prefix string) *RouterGroup {
	engine := groub.engine
	newGroup := &RouterGroup{
		parent: groub,
		prefix: groub.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (groub *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := groub.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	groub.engine.router.addRoute(method, pattern, handler)
}

func (groub *RouterGroup) GET(pattern string, handler HandlerFunc) {
	groub.addRoute("GET", pattern, handler)
}

func (groub *RouterGroup) POST(pattern string, handler HandlerFunc) {
	groub.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
