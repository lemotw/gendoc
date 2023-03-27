package analysis

import (
	"golang.org/x/net/html"
)

func SearchNodes(node *html.Node, nodeType html.NodeType, target string) []*html.Node {
	res := []*html.Node{}
	stack := []*html.Node{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		stack = append(stack, child)
	}

	for len(stack) > 0 {
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node.Type == nodeType {
			if target == "" {
				res = append(res, node)
			} else if node.Data == target {
				res = append(res, node)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return res
}

func FindNode(root *html.Node, target string) *html.Node {
	stack := []*html.Node{root}

	var node *html.Node
	for len(stack) > 0 {
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node == nil {
			continue
		}

		if node.Type == html.ElementNode && node.Data == target {
			return node
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return nil
}
