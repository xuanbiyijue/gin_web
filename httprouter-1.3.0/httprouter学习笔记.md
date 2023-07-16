# HttpRouter 源码解析
源码地址: https://github.com/julienschmidt/httprouter  
源码文件主要有3个: **path.go, router.go, tree.go**.

## path.go 
此文件仅有一个函数: 
```go
func CleanPath(p string) string
```
作用: 它返回规范的 URL 路径。具体规则如下:  
* 将多个斜杠替换为单个斜杠
* 清除每个 '.'(表示当前路径) 与 '..'(表示上级路径)


## router.go  
* `type Handle func(http.ResponseWriter, *http.Request, Params)`  
这其实是给函数 `func(http.ResponseWriter, *http.Request, Params)` 起一个别名叫: Handle。
Handle 是对应 method + path 的处理函数，当一个请求发送过来，实际上就是找到这个请求对应的 Handle 来执行。  
相比于 `http.HandlerFunc -> type HandlerFunc func(ResponseWriter, *Request)`，Handle 多一个参数 `Params`, 
允许携带参数的动态路由， `Params` 是键值对 `Param` 的数组，可以通过 index 或 `ByName` 获取参数值   


* 路由引擎 `Router`  
  路由引擎 `Router` 是一个结构体，包含的变量如下(通过 `httprouter.New()` 来创建一个路由引擎。):  
```go
type Router struct {
	// 按照 mothods 分别保存的基数树，例如: Get-基数树
	trees map[string]*node  

	// 是否将带有尾斜杠的路由重定向到不带有尾斜杠的路由，例如: '/foo/' -> '/foo'
	RedirectTrailingSlash bool

	// 是否修正并重定向路由，如果开启，首先会把多余的path删除，之后进行不区分大小写地匹配
	RedirectFixedPath bool

	// 是否当路由没匹配到时查看其他方法能否匹配到，如果开启且在其他方法匹配到则返回 405 状态码以及 "Method Not Allowed" 消息
	HandleMethodNotAllowed bool

	// 是否自动回复 OPTIONS 请求
	HandleOPTIONS bool

	// 可选的选项。与 HandleOPTIONS 搭配使用
	GlobalOPTIONS http.Handler

	// Cached value (缓存值) of global (*) allowed methods
	globalAllowed string

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}
```


* Methods 与 r.Handle
  * `func (r *Router) GET(path string, handle Handle)`
  * `func (r *Router) HEAD(path string, handle Handle)`
  * `func (r *Router) OPTIONS(path string, handle Handle)`
  * `func (r *Router) POST(path string, handle Handle)`
  * `func (r *Router) PUT(path string, handle Handle)`
  * `func (r *Router) PATCH(path string, handle Handle)`
  * `func (r *Router) DELETE(path string, handle Handle)`  

  以上 Motheds 实际上都是对 `func (r *Router) Handle(method, path string, handle Handle)` 的调用，
可以说是一种快捷方式。  
  * `func (r *Router) Handle(method, path string, handle Handle)`, 还有一些 `r.Handle` 的衍生:
    * `func (r *Router) Handler(method, path string, handler http.Handler)` 用来兼容 http.Handler 的
    * `func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc)` 用来兼容 http.HandlerFunc，函数内部调用 `r.Handler`  
  
    因此，HttpRouter.Router 能够兼容 http.HandlerFunc (http.HandlerFunc 实现了 http.Handler 接口)


* 前端模板文件  
通过函数 `func (r *Router) ServeFiles(path string, root http.FileSystem)` 来实现。
这个函数实际上是通过调用 `http.FileServer()` 来实现的，并且在之后为模板文件创建 Get 模式的路由。
`http.FileServer()` 创建并返回一个 `fileHandler` (它也实现了http.Handler 接口)。


* `func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)`  
这个函数主要就是来根据给定的 method 从对应的基数树，再根据给定的 path 获得 Handle 函数并调用。

## tree.go
* 前缀树  
  

* 基数树  
  基数树是空间优化的Trie，当某个节点的子节点唯一时，将子节点与该节点合并.  
  * 节点的定义
  ```go
  type node struct {
    path      string   // 包含的路径片段
    wildChild bool     // 子节点是否为参数节点
    nType     nodeType // 节点类型：静态（默认）、根、命名参数捕获、任意参数捕获
    maxParams uint8    // 最大参数个数
    priority  uint32   // 优先级
    indices   string   // 索引
    children  []*node  // 子节点
    handle    Handle   // 该节点所代表路径的 handle
  }
  ```

  * 插入一个节点  
  调用 `func (n *node) addRoute(path string, handle Handle)`，逻辑如下：
    1. 如果是空树，直接插入；否则，进行下一步
    2. 


## 启动服务
启动服务通过调用 `http.ListenAndServe(":8080", router)`

实际上 httprouter.Router 本质上就是重写的 http.Handler



## 从一个请求开始