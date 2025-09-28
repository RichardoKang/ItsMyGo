package gee

import "strings"

type node struct {
	pattern  string // 待匹配的路由，例如 /p/:lang
	part     string // 路由中的一部分，例如 :lang
	children []*node
	isWild   bool // 是否模糊匹配，part 含有 : 或 * 时为true
}

// 寻找匹配的节点，如果有一致或者通配符的字段，则返回该节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 存储所有能够匹配的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // 递归终止条件，表示已经插入到最后一层
		n.pattern = pattern
		return
	}

	part := parts[height]       // 当前最末端节点
	child := n.matchChild(part) // 寻找匹配的节点
	if child == nil {
		isWild := false
		if len(part) > 0 {
			isWild = part[0] == ':' || part[0] == '*'
		}
		child = &node{part: part, isWild: isWild}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) // 递归插入节点
}

// 搜索匹配的节点，返回该节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // 如果已经搜索到最后一层，或者遇到通配符
		if n.pattern == "" {
			return nil // 如果没有匹配到节点，则返回nil
		}
		return n // 返回匹配到的节点
	}

	part := parts[height]             // 当前最末端节点
	children := n.matchChildren(part) // 寻找匹配的节点
	// 遍历所有匹配的节点，递归搜索

	for _, child := range children {
		result := child.search(parts, height+1) // 递归搜索

		if result != nil {
			return result // 如果找到匹配的节点，则返回该节点
		}
	}
	return nil
}
