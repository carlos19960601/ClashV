package trie

import "strings"

type Node[T any] struct {
	childMap  map[string]*Node[T]
	childNode *Node[T]
	childStr  string
	inited    bool
	data      T
}

func newNode[T any]() *Node[T] {
	return &Node[T]{}
}

func (n *Node[T]) isEmpty() bool {
	if n == nil || !n.inited {
		return true
	}

	return false
}

func (n *Node[T]) getChild(s string) *Node[T] {
	if n.childMap == nil {
		if n.childNode != nil && n.childStr == s {
			return n.childNode
		}
		return nil
	}

	return n.childMap[s]
}

func (n *Node[T]) getOrNewChild(s string) *Node[T] {
	node := n.getChild(s)
	if node == nil {
		node = newNode[T]()
		n.addChild(s, node)
	}

	return node
}

func (n *Node[T]) addChild(s string, child *Node[T]) {
	if n.childMap == nil {
		if n.childNode == nil {
			n.childStr = s
			n.childNode = child
			return
		}
		n.childMap = map[string]*Node[T]{}
		if n.childNode != nil {
			n.childMap[n.childStr] = n.childNode
		}
		n.childStr = ""
		n.childNode = nil
	}

	n.childMap[s] = child
}

func (n *Node[T]) setData(data T) {
	n.data = data
	n.inited = true
}

func (n *Node[T]) optimize() {
	if len(n.childStr) > 0 {
		n.childStr = strClone(n.childStr)
	}
	if n.childNode != nil {
		n.childNode.optimize()
	}

	if n.childMap == nil {
		return
	}

	switch len(n.childMap) {
	case 0:
		n.childMap = nil
		return
	case 1:
		for key := range n.childMap {
			n.childStr = key
			n.childNode = n.childMap[key]
		}
		n.childMap = nil
		n.optimize()
		return
	}

	children := make(map[string]*Node[T], len(n.childMap))
	for key := range n.childMap {
		child := n.childMap[key]
		if child == nil {
			continue
		}

		key = strClone(key)
		children[key] = child
		child.optimize()
	}
	n.childMap = children
}

func strClone(key string) string {
	switch key {
	case wildcard:
		key = wildcard
	case dotWildcard:
		key = dotWildcard
	case complexWildcard:
		key = complexWildcard
	case domainStep:
		key = domainStep
	default:
		key = strings.Clone(key)
	}
	return key
}
