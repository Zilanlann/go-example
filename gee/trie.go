package gee

import "strings"

// A node in the router tree.
type node struct {
	// The pattern to match.
	pattern string

	// The part of the pattern that is being matched.
	part string

	// The child nodes.
	children []*node

	// Whether the node is a wildcard match.
	isWild bool
}

// Find the first child node that matches the given part.
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// Find all child nodes that match the given part.
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// Insert a new node into the tree.
func (n *node) insert(pattern string, parts []string, height int) {
	// If we've reached the end of the pattern, set the node's pattern and return.
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// Get the next part of the pattern.
	part := parts[height]

	// Find the child node that matches the part.
	child := n.matchChild(part)

	// If there is no child node that matches the part, create a new one.
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// Recursively insert the rest of the pattern into the child node.
	child.insert(pattern, parts, height+1)
}

// Search for a node in the tree that matches the given parts.
func (n *node) search(parts []string, height int) *node {
	// If we've reached the end of the parts, return the node.
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// Get the next part of the parts.
	part := parts[height]

	// Find all child nodes that match the part.
	children := n.matchChildren(part)

	// Recursively search for a node in the child nodes that matches the rest of the parts.
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	// No node was found that matches the parts.
	return nil
}
