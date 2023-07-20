# Gee Web 框架实现

## net/http 库

go 标准库 `net/http` 提供了基础的Web功能，即监听端口，解析HTTP报文，映射静态路由。
但是，他并不支持某些功能:

1. 动态路由：例如 `hello/:name`，`hello/*`这类的规则。
2. 分组和鉴权：没有分组/统一鉴权的能力，需要在每个路由映射的 `handler` 中实现

## 实现 http.Handler 接口

首先看看 `http.Handler` 接口源码:

```go
package http

type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```

其中， `ListenAndServe()` 要求实现 `http.Handler` 接口，就可以调用此方法

```go
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// ...
}
```

## Gee 的实现

先简单实现一个雏形:

```go
package gee

import (
    "fmt"
    "net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP. 
// 采用 hashmap 来实现映射，此时不支持动态路由
type Engine struct {
    router map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func New() *Engine {
    return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    key := method + "-" + pattern
    engine.router[key] = handler
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    key := req.Method + "-" + req.URL.Path
    if handler, ok := engine.router[key]; ok {
        handler(w, req)
    } else {
        fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
    }
}
```
```go
// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
    engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
    engine.addRoute("POST", pattern, handler)
}
```


## 上下文 Context
对Web服务来说，无非是根据请求 `*http.Request`，构造响应 `http.ResponseWriter`。  
但是这两个对象提供的接口粒度太细，比如我们要构造一个完整的响应，需要考虑消息头(Header)和消息体(Body)，
而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等几乎每次请求都需要设置的信息。  
因此，如果不进行有效的封装，那么框架的用户将需要写大量重复，繁杂的代码，而且容易出错。  
举个例子:   
```go
r := gee.New()
r.POST("/login", func(c *gee.Context) {
	// 函数内容
    c.JSON(http.StatusOK, gee.H{
        "username": c.PostForm("username"),
        "password": c.PostForm("password"),
    })
})
```
对于上面这段代码，如果没有封装 `context`，那么就会是下面这样:  
```go
obj := map[string]interface{}{
    "name": "geektutu",
    "password": "1234",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
encoder := json.NewEncoder(w)
if err := encoder.Encode(obj); err != nil {
    http.Error(w, err.Error(), 500)
}
```

另外，对于框架来说，还需要支撑额外的功能。
例如，将来解析动态路由/hello/:name，参数:name的值放在哪呢？
再比如，框架需要支持中间件，那中间件产生的信息放在哪呢？  
Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载。  
因此，设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。路由的处理函数，以及将要实现的中间件，
参数都统一使用 Context 实例， Context 就像一次会话的百宝箱，可以找到任何东西。
```go
type H map[string]interface{}

type Context struct {
    // origin objects
    Writer http.ResponseWriter
    Req    *http.Request
    // request info
    Path   string
    Method string
    // response info
    StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
    return &Context{
        Writer: w,
        Req:    req,
        Path:   req.URL.Path,
        Method: req.Method,
    }
}

// Status 设置响应 code
func (c *Context) Status(code int) {
    c.StatusCode = code
    c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
    c.Writer.Header().Set(key, value)
}
```
```go
// 提供了访问 Query 和 PostForm 参数的方法。
// PostForm 根据 key 获得表单值
func (c *Context) PostForm(key string) string {
    return c.Req.FormValue(key)
}

// Query 根据 key 获得动态路由中的参数值
func (c *Context) Query(key string) string {
    return c.Req.URL.Query().Get(key)
}
```
```go
// 提供了快速构造String/Data/JSON/HTML响应的方法。
func (c *Context) String(code int, format string, values ...interface{}) {
    c.SetHeader("Content-Type", "text/plain")
    c.Status(code)
    c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
    c.SetHeader("Content-Type", "application/json")
    c.Status(code)
    encoder := json.NewEncoder(c.Writer)
    if err := encoder.Encode(obj); err != nil {
        http.Error(c.Writer, err.Error(), 500)
    }
}

func (c *Context) Data(code int, data []byte) {
    c.Status(code)
    c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
    c.SetHeader("Content-Type", "text/html")
    c.Status(code)
    c.Writer.Write([]byte(html))
}
```
那么使用时的效果:  
```go
func main() {
    r := gee.New()
    r.GET("/", func(c *gee.Context) {
        c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
    })
    r.GET("/hello", func(c *gee.Context) {
        // expect /hello?name=geektutu
        c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
    })
    
    r.POST("/login", func(c *gee.Context) {
        c.JSON(http.StatusOK, gee.H{
            "username": c.PostForm("username"),
            "password": c.PostForm("password"),
        })
    })
    
    r.Run(":9999")
}
```

### 加入动态路由后的修改
需要对 Context 对象增加一个属性和方法，来提供对路由参数的访问。
将解析后的参数存储到 `Params` 中，通过 `c.Param()` 的方式获取到对应的值。  
```go
type Context struct {
    // origin objects
    Writer http.ResponseWriter
    Req    *http.Request
    // request info
    Path   string
    Method string
    Params map[string]string  // new attribution
    // response info
    StatusCode int
}

func (c *Context) Param(key string) string {
    value, _ := c.Params[key]
    return value
}
```




## 路由 Router
将和路由相关的方法和结构提取出来，方便以后的扩展:  
```go
type router struct {
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    log.Printf("Route %4s - %s", method, pattern)
    key := method + "-" + pattern
    r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
    key := c.Method + "-" + c.Path
    if handler, ok := r.handlers[key]; ok {
        handler(c)
    } else {
        c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
    }
}
```
注意，此时在 gee.go 同样进行修改:  
```go
// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)
// type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
    router *router
    // 原来的如下：
    // router map[string]HandlerFunc
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := newContext(w, req)
    engine.router.handle(c)
}
```

### 动态路由(基于前缀树)
![img](https://geektutu.com/post/gee-day3/trie_router.jpg)  

HTTP请求的路径恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。
我们通过树结构查询，如果中间某一层的节点都不满足条件，那么就说明没有匹配到的路由，查询结束。  

接下来我们实现的动态路由具备以下两个功能。  
* 参数匹配 `:` 。例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
* 通配 `*` 。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器，能够递归地匹配子路径。

**前缀树 (Trie 树)**  
```go
type node struct {
    pattern  string      // 待匹配路由，例如 /p/:lang, 也就是从根到此节点的完整路径
    part     string      // 路由中的一部分，例如 :lang, 也就是此节点应该保存的值
    children []*node     // 子节点，例如 [doc, tutorial, intro]
    isWild   bool        // 是否精确匹配，part 含有 : 或 * 时为true
}
```
关于前缀树节点的实现，很容易想到其数据结构中应该有以下几项:  
* 节点内保存的值，例如: `pattern`、`part`
* 存储子节点的指针的数据结构: `children`
* 另外，还需要一个 `isWild` 参数来实现动态路由匹配  

Trie 树需要支持节点的插入与查询。
* 插入: 递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
* 查询: 同样也是递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点。
  有一点需要注意，`/p/:lang/doc`只有在第三层节点，即doc节点，pattern才会设置为`/p/:lang/doc`。
  p和:lang节点的pattern属性皆为空。因此，当匹配结束时，我们可以使用`n.pattern == ""`来判断路由规则是否匹配成功。
  比如说: 我注册了一个路由 `/p/:lang/doc`，那么前缀树会创建3个节点: `/p, /:lang, /doc`，那么这两个节点是不能被匹配到的。
  因为实际上并没有创建相应的路由。因此，他们的 `pattern = ""`。当我使用没注册的路由 `/p/python` 虽能成功匹配到 `:lang`，但`:lang`的pattern值为空，因此匹配失败。

```go
// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}
// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
    nodes := make([]*node, 0)
    for _, child := range n.children {
        if child.part == part || child.isWild {
            nodes = append(nodes, child)
        }
    }
    return nodes
}
```
```go
func (n *node) insert(pattern string, parts []string, height int) {
    if len(parts) == height {
        n.pattern = pattern
        return
    }
    part := parts[height]
    child := n.matchChild(part)
    if child == nil {
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }
    child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.part, "*") {
        if n.pattern == "" {
            return nil
        }
        return n
    }
    part := parts[height] 
    // 这里能匹配到多个是因为除了指定的节点，还有可能匹配到带有参数的节点
    children := n.matchChildren(part)
    for _, child := range children {
        result := child.search(parts, height+1)
        if result != nil {
            return result
        }
    }
    return nil
}
```

至此，前缀树的插入与查找实现完成。  

**前缀树应用在 Router**
```go
type router struct {
    // roots key eg, roots['GET'] roots['POST']
    roots    map[string]*node
    // handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{
        roots:    make(map[string]*node),
        handlers: make(map[string]HandlerFunc),
    }
}
```
```go
// Only one * is allowed
func parsePattern(pattern string) []string {
    vs := strings.Split(pattern, "/")
    parts := make([]string, 0)
    for _, item := range vs {
        if item != "" {
            parts = append(parts, item)
            if item[0] == '*' {break}
        }
    }
    return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    parts := parsePattern(pattern)
    key := method + "-" + pattern
    _, ok := r.roots[method]
    if !ok {
        // 根节点，值为空字符串
        r.roots[method] = &node{}
    }
    r.roots[method].insert(pattern, parts, 0)
    r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
    searchParts := parsePattern(path)
    params := make(map[string]string)
    root, ok := r.roots[method]
    if !ok {return nil, nil}
  
    n := root.search(searchParts, 0)
    if n != nil {
        parts := parsePattern(n.pattern)
        // 检查是否动态路由，并且记录动态路由的参数键值对
        for index, part := range parts {
            if part[0] == ':' {
                params[part[1:]] = searchParts[index]
            }
            if part[0] == '*' && len(part) > 1 {
				// 将剩余 parts 拼接起来作为参数
                params[part[1:]] = strings.Join(searchParts[index:], "/")
                break
            }
        }
        return n, params
    }
    return nil, nil
}
```

之后对 `Context` 进行相应的修改。

## 分组控制 Route Group Control
考虑以下场景:  
* 以/post开头的路由匿名可访问。
* 以/admin开头的路由需要鉴权。
* 以/api开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权。  

因此，分组控制(Group Control)是 Web 框架应提供的基础功能之一。  

一个 Group 对象需要具备哪些属性呢？
* 首先是前缀(**prefix**)，比如/，或者/api；
* 要支持分组嵌套，那么需要知道当前分组的父亲(**parent**)是谁；
* 中间件要求应用在分组上的，那还需要存储应用在该分组上的中间件(**middlewares**)  
* Group对象还需要有访问Router的能力，为了方便，我们可以在Group中，保存一个指针，指向 **Engine**，
  整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了。

```go
RouterGroup struct {
    prefix      string
    middlewares []HandlerFunc // support middleware
    parent      *RouterGroup  // support nesting
    engine      *Engine       // all groups share a Engine instance
}
```
还可以进一步地抽象，将Engine作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力。  
```go
Engine struct {
    *RouterGroup
    router *router
    groups []*RouterGroup // store all groups
}
```

那我们就可以将和路由有关的函数，都交给RouterGroup实现了。
```go
// New is the constructor of gee.Engine
func New() *Engine {
    engine := &Engine{router: newRouter()}
    engine.RouterGroup = &RouterGroup{engine: engine}
    engine.groups = []*RouterGroup{engine.RouterGroup}
    return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
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

// 对于路由组而言，只需要关注在添加路由映射时，把前缀加上。而获得handler时则需要完整的path
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
    pattern := group.prefix + comp
    log.Printf("Route %4s - %s", method, pattern)
    group.engine.router.addRoute(method, pattern, handler)
}

// engine 也能用 GET，POST 方法
// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
    group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
    group.addRoute("POST", pattern, handler)
}
```

## 中间件的支持
对中间件而言，需要考虑2个比较关键的点：
* 插入点位置
* 中间件的输入是什么  

Gee 的中间件的定义与路由映射的 Handler 一致，处理的输入是Context对象。  
插入点是框架接收到请求初始化Context对象后，允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等  
另外通过调用(*Context).Next()函数，中间件可等待用户自己定义的 Handler处理结束后，再做一些额外的操作，例如计算本次处理所用时间等。  
还需要支持设置多个中间件，依次进行调用。  

综上，接收到请求后，应查找所有应作用于该路由的中间件，保存在Context中，依次进行调用。  
为此，给 `Context` 添加2个参数，并定义 `Next` 方法：  
```go
type Context struct {
    // origin objects
    Writer http.ResponseWriter
    Req    *http.Request
    // request info
    Path   string
    Method string
    Params map[string]string
    // response info
    StatusCode int
    // middleware
    handlers []HandlerFunc
    index    int
}

func (c *Context) Next() {
    c.index++
    s := len(c.handlers)
	// 手动调用，因为并不是所有中间件都有 Next()
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}
```

那么就实现了中间件的调用过程:  
假设我们应用了中间件 `A` 和 `B`，和路由映射的 `Handler`。`c.handlers` 是这样的 `[A, B, Handler]`，
`c.index` 初始化为 `-1`。调用 `c.Next()`，接下来的流程是这样的：
* c.index++，c.index 变为 0, 0 < 3，调用 c.handlers[0]，即 A
* A 执行完部分代码，调用 c.Next()，c.index++，c.index 变为 1，1 < 3，调用 c.handlers[1]，即 B
* B 执行完部分代码，调用 c.Next()，c.index++，c.index 变为 2，2 < 3，调用 c.handlers[2]，即Handler
* Handler 调用完毕，返回到 B 中的后续代码并执行
* B 执行完毕，返回到 A 中的后续代码并执行
* A 执行完毕，结束。  

之后在 `gee.go` 中定义 `Use` 函数，用于给路由组注册中间件:  
```go
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
    group.middlewares = append(group.middlewares, middlewares...)
}
```

不同的路由组有不同的中间件，因此，还需要实现这部分功能代码，我们在 `ServeHTTP` 中实现:  
```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    var middlewares []HandlerFunc
    for _, group := range engine.groups {
        if strings.HasPrefix(req.URL.Path, group.prefix) {
            middlewares = append(middlewares, group.middlewares...)
        }
    }
    c := newContext(w, req)
    c.handlers = middlewares
    engine.router.handle(c)
}
```

相应的，还需要修改 `engine.router.handle` 函数:  
```go
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
    // 这样就可以从中间件开始调用
    c.Next()
}
```

## 模板文件与静态文件  
### 静态文件  
要做到服务端渲染，第一步便是要支持 JS、CSS 等静态文件。还记得我们之前设计动态路由的时候，
支持通配符*匹配多级子路径。比如路由规则/assets/*filepath，可以匹配/assets/开头的所有的地址。
例如/assets/js/geektutu.js，匹配后，参数filepath就赋值为js/geektutu.js。  
gee 框架只需要解析请求的地址，而映射到服务器上文件的真实地址，交给 `http.FileServer` 处理就好了。  
```go
// create static handler, http.FileSystem 是接口类型，string 类型实现了这个接口
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
    absolutePath := path.Join(group.prefix, relativePath)
	// fileServer 是实现了 ServeHTTP 接口的 Handler
    fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
    return func(c *Context) {
        file := c.Param("filepath")
        // Check if file exists and/or if we have permission to access it
        if _, err := fs.Open(file); err != nil {
            c.Status(http.StatusNotFound)
            return
        }
        fileServer.ServeHTTP(c.Writer, c.Req)
    }
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	// http.Dir() 是类型转换
    handler := group.createStaticHandler(relativePath, http.Dir(root))
    urlPattern := path.Join(relativePath, "/*filepath")
    // Register GET handlers
    group.GET(urlPattern, handler)
}
```
给 `RouterGroup` 添加了2个方法，`Static` 这个方法是暴露给用户的。
用户可以将磁盘上的某个文件夹 `root` 映射到路由 `relativePath`。例如：  
```go
r := gee.New()
r.Static("/assets", "/usr/geektutu/blog/static")
// 或相对路径 r.Static("/assets", "./static")
r.Run(":9999")
```
用户访问 `localhost:9999/assets/js/geektutu.js`，最终返回 `/usr/geektutu/blog/static/js/geektutu.js`。  


### HTML 模板渲染
使用 `html/template` 提供的方法。
首先为 Engine 示例添加了 `*template.Template` 和 `template.FuncMap` 对象，前者将所有的模板加载进内存，后者是所有的自定义模板渲染函数。
另外，给用户分别提供了设置自定义渲染函数 `funcMap` 和加载模板的方法。  
```go
Engine struct {
    *RouterGroup
    router        *router
    groups        []*RouterGroup     // store all groups
    htmlTemplates *template.Template // for html render
    funcMap       template.FuncMap   // for html render
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
    engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
    engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
```
接下来，对原来的 (*Context).HTML()方法做了些小修改，使之支持根据模板文件名选择模板进行渲染。我们在 Context 中添加了成员变量 engine *Engine，
这样就能够通过 Context 访问 Engine 中的 HTML 模板。实例化 Context 时，还需要给 c.engine 赋值。
```go
type Context struct {
    // ...
    // engine pointer
    engine *Engine
}

func (c *Context) HTML(code int, name string, data interface{}) {
    c.SetHeader("Content-Type", "text/html")
    c.Status(code)
    if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
        c.Fail(500, err.Error())
    }
}
```

### 使用
项目结构:  
```
---gee/
---static/
   |---css/
        |---geektutu.css
   |---file1.txt
---templates/
   |---arr.tmpl
   |---css.tmpl
   |---custom_func.tmpl
---main.go
```
```
<html>
    <link rel="stylesheet" href="/assets/css/geektutu.css">
    <p>geektutu.css is loaded</p>
</html>
```
那么，main 函数这样写:  
```go
type student struct {
    Name string
    Age  int8
}

func FormatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
    r := gee.New()
    r.Use(gee.Logger())
    r.SetFuncMap(template.FuncMap{
        "FormatAsDate": FormatAsDate,
    })
    r.LoadHTMLGlob("templates/*")
    r.Static("/assets", "./static")
  
    stu1 := &student{Name: "Geektutu", Age: 20}
    stu2 := &student{Name: "Jack", Age: 22}
    r.GET("/", func(c *gee.Context) {
        c.HTML(http.StatusOK, "css.tmpl", nil)
    })
    r.GET("/students", func(c *gee.Context) {
        c.HTML(http.StatusOK, "arr.tmpl", gee.H{
            "title":  "gee",
            "stuArr": [2]*student{stu1, stu2},
        })
    })
    r.GET("/date", func(c *gee.Context) {
        c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
            "title": "gee",
            "now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
        })
    })
    r.Run(":9999")
}
```

## 错误恢复  
对一个 Web 框架而言，错误处理机制是非常必要的。可能是框架本身没有完备的测试，
导致在某些情况下出现空指针异常等情况。也有可能用户不正确的参数，触发了某些异常，例如数组越界，空指针等。
如果因为这些原因导致系统宕机，必然是不可接受的。    

需要在 gee 中添加一个非常简单的错误处理机制，即在此类错误发生时，向用户返回 Internal Server Error，
并且在日志中打印必要的错误信息，方便进行错误定位。我们之前实现了中间件机制，错误处理也可以作为一个中间件，增强 gee 框架的能力。  
```go
// print stack trace for debug
func trace(message string) string {
    var pcs [32]uintptr
    n := runtime.Callers(3, pcs[:]) // skip first 3 caller
  
    var str strings.Builder
    str.WriteString(message + "\nTraceback:")
    for _, pc := range pcs[:n] {
        fn := runtime.FuncForPC(pc)
        file, line := fn.FileLine(pc)
        str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
    }
    return str.String()
}

func Recovery() HandlerFunc {
    return func(c *Context) {
        defer func() {
            if err := recover(); err != nil {
                message := fmt.Sprintf("%s", err)
                log.Printf("%s\n\n", trace(message))
                c.Fail(http.StatusInternalServerError, "Internal Server Error")
            }
        }()
  
        c.Next()
    }
}
```
Recovery 的实现非常简单，使用 defer 挂载上错误恢复的函数，在这个函数中调用 *recover()*，
捕获 panic，并且将堆栈信息打印在日志中，向用户返回 Internal Server Error。  

trace() 函数，这个函数是用来获取触发 panic 的堆栈信息。在 trace() 中，
调用了 runtime.Callers(3, pcs[:])，Callers 用来返回调用栈的程序计数器, 
第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 defer func。
因此，为了日志简洁一点，我们跳过了前 3 个 Caller。
接下来，通过 runtime.FuncForPC(pc) 获取对应的函数，在通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中。

至此，gee 框架的错误处理机制就完成了。
