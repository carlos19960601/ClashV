package trie

import (
	"errors"
	"strings"
)

const (
	wildcard        = "*"
	dotWildcard     = ""
	complexWildcard = "+"
	domainStep      = "."
)

var ErrInvalidDomain = errors.New("invalid domain")

type DomainTrie[T any] struct {
	root *Node[T]
}

func New[T any]() *DomainTrie[T] {
	return &DomainTrie[T]{
		root: newNode[T](),
	}
}

// Support
// 1. www.example.com
// 2. *.example.com
// 3. subdomain.*.example.com
// 4. .example.com
// 5. +.example.com
func (t *DomainTrie[T]) Insert(domain string, data T) error {
	parts, valid := ValidAndSplitDomain(domain)
	if !valid {
		return ErrInvalidDomain
	}

	if parts[0] == complexWildcard {
		t.insert(parts[1:], data)
		parts[0] = dotWildcard
		t.insert(parts, data)
	} else {
		t.insert(parts, data)
	}

	return nil
}

func (t *DomainTrie[T]) Optimize() {
	t.root.optimize()
}

func (t *DomainTrie[T]) insert(parts []string, data T) {
	node := t.root
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		node = node.getOrNewChild(part)
	}

	node.setData(data)
}

func (t *DomainTrie[T]) Search(domain string) *Node[T] {
	parts, valid := ValidAndSplitDomain(domain)
	if !valid || parts[0] == "" {
		return nil
	}

	n := t.search(t.root, parts)
	if n.isEmpty() {
		return nil
	}

	return n
}

func (t *DomainTrie[T]) search(node *Node[T], parts []string) *Node[T] {
	if len(parts) == 0 {
		return node
	}

	if c := node.getChild(parts[len(parts)-1]); c != nil {
		if n := t.search(c, parts[:len(parts)-1]); !n.isEmpty() {
			return n
		}
	}

	if c := node.getChild(wildcard); c != nil {
		if n := t.search(c, parts[:len(parts)-1]); !n.isEmpty() {
			return n
		}
	}

	return node.getChild(dotWildcard)
}

func ValidAndSplitDomain(domain string) ([]string, bool) {
	if domain != "" && domain[len(domain)-1] == '.' {
		return nil, false
	}

	domain = strings.ToLower(domain)
	parts := strings.Split(domain, domainStep)
	if len(parts) == 1 {
		if parts[0] == "" {
			return nil, false
		}

		return parts, true
	}

	for _, part := range parts[1:] {
		if part == "" {
			return nil, false
		}
	}

	return parts, true
}
