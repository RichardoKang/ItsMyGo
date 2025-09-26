package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 将路由解析成字符串数组
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") // 按照 / 分割路由

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果遇到通配符，则停止继续分割
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) // 解析路由，比如 /p/:lang -> [p, :lang]

	key := method + "-" + pattern // 生成路由对应的键，比如 GET-/p/:lang
	_, ok := r.roots[method]      // 获取该方法对应的根节点
	if !ok {
		r.roots[method] = &node{} // 如果不存在，则创建一个新的根节点
	}
	r.roots[method].insert(pattern, parts, 0) // 插入节点
	r.handlers[key] = handler                 // 记录路由和处理函数的映射
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 解析请求路径，比如 /p/python -> [p, python]
	params := make(map[string]string) // 用于存储路由参数，比如 {lang: python}

	root, ok := r.roots[method] // 获取该方法对应的根节点
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // 搜索匹配的节点

	if n != nil {
		parts := parsePattern(n.pattern) // 解析路由，比如 /p/:lang -> [p, :lang]
		// 遍历路由的每一部分，提取路由参数
		for index, part := range parts {
			if part[0] == ':' { // 如果是参数，比如 :lang
				params[part[1:]] = searchParts[index] // 提取参数值，比如 {lang: python}
			}
			if part[0] == '*' && len(part) > 1 { // 如果是通配符，比如 *filepath
				// 提取通配符参数值，比如 {filepath: a/b/c}
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params // 返回匹配的节点和路由参数
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path) // 获取匹配的节点和路由参数
	if n != nil {
		c.Params = params // 设置路由参数
		// 根据请求方法和路由生成键，获取对应的处理函数
		key := c.Method + "-" + n.pattern
		r.handlers[key](c) // 调用处理函数
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
