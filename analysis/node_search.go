package analysis

import "golang.org/x/net/html"

func SearchNodes(node *html.Node, nodeType html.NodeType) []*html.Node {
	res := []*html.Node{}
	stack := []*html.Node{node}
	for len(stack) > 0 {
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node.Type == nodeType {
			res = append(res, node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return res
}

func FindNode(node *html.Node, target string) *html.Node {
	stack := []*html.Node{node}

	for len(stack) > 0 {
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node.Type == html.ElementNode && node.Data == target {
			return node
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return nil
}
