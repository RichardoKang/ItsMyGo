package gee

type node struct {
	pattern  string // 待匹配的路由，例如 /p/:lang
	part     string // 路由中的一部分，例如 :lang
	children []*node
	isWild   bool // 是否模糊匹配，part 含有 : 或 * 时为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {

	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		isWild := false
		if len(part) > 0 {
			isWild = part[0] == ':' || part[0] == '*'
		}
		child = &node{part: part, isWild: isWild}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)

}
