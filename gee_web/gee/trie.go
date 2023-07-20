package gee

import (
	"fmt"
	"strings"
)

/*
从 req 中获取的path是 /p/go/doc
设置的动态路由是 /p/:lang/doc
 */

// node 前缀树的数据结构
type node struct {
	pattern  string      // 待匹配路由，只有注册过的路由才会有值，而路由的中间节点值为空
	part     string      // 路由的一部分（当前节点路由）
	children []*node     // 子节点
	isWild   bool        // 是否精确匹配(用于匹配动态路由)，模糊匹配时为 true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert 注册路由
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果要注册的路由 /p/:lang/doc
	// 那么 parts = [p, :lang, doc]，height = 0

	// 终止条件，n.pattern 为从根节点到当前节点的路径拼接而成
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// part 为下一层路由，因为从根节点(/)开始
	part := parts[height]
	// 获得与 part 匹配的子节点
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 调用子节点的 insert
	child.insert(pattern, parts, height+1)
}

// search 路由匹配
func (n *node) search(parts []string, height int) *node {
	// 如果要匹配的路由 /p/go/doc
	// 那么 parts = [p, go, doc]，height = 0
	// 而在树中实际为 [p, :lang, doc]

	// strings.HasPrefix(n.part, "*") 检测 n.part 是否以 * 为前缀
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果 n.pattern == "" ，那就是还没注册过的路由
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// matchChild 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 如果 子节点的值为目标指 或者 模糊查找
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
