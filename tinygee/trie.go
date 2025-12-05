package tinygee

// node 表示路由树的一个节点。
type node struct {
	pattern  string  // 完整匹配的路由, 如 /p/:lang
	part     string  // 路由中的一段，如 :lang
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配，: 或 *
}

// matchChild 找到第一个匹配的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 找到所有匹配的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 将 pattern 插入 trie
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search 根据路径搜索匹配的节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || (len(n.part) > 0 && n.part[0] == '*') {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		if result := child.search(parts, height+1); result != nil {
			return result
		}
	}
	return nil
}
