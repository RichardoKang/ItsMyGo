package gee

import (
	"net/http"
	"strings"
)

type (
	HandlerFunc func(c *Context)
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		engine      *Engine
		parent      *RouterGroup
	}

	Engine struct {
		*RouterGroup
		groups []*RouterGroup
		router *router
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp // e.g. /v1 + /hello
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 把中间件添加到当前路由组
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 如果请求的路由前缀符合某个路由组的前缀，就将该路由组的中间件添加到中间件链中
		// e.g. /v1 和 /v1/hello
		// 只有当请求的路由是以该路由组的前缀开头时，才会添加该路由组的中间件
		// 这样就可以实现对不同路由组使用不同的中间件
		// 例如，可以为 /v1 路由组添加身份验证中间件，而为 /v2 路由组添加日志记录中间件
		// 这样，当请求的路由是 /v1/hello 时，就会执行身份验证中间件
		// 当请求的路由是 /v2/hello 时，就会执行日志记录中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	// 把中间件链赋值给 Context，以便在处理请求时依次调用这些中间件
	c.handlers = middlewares
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
