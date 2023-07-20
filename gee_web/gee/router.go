/*
将和路由相关的方法和结构提取了出来
方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
 */

package gee

import (
	"net/http"
	"strings"
)

type router struct {
	// roots 存储每种请求方式的 Trie 树根节点
	// roots key eg, roots['GET'] roots['POST']
	roots     map[string]*node
	handlers  map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// Only one * is allowed
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
//	key := method + "-" + pattern
//	r.handlers[key] = handler
//}

// addRoute 注册路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 从对应的树中插入路由
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 获得路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 例如：[p, go, doc]
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	// 如果有这条路径
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 拿到路由中的参数
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

//func (r *router) handle(c *Context) {
//	n, params := r.getRoute(c.Method, c.Path)
//	if n != nil {
//		c.Params = params
//		key := c.Method + "-" + n.pattern
//		r.handlers[key](c)
//	} else {
//		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
//	}
//}
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 开始执行
	c.Next()
}