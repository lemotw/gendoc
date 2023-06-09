package model

import (
	"strings"

	"golang.org/x/net/html"
)

type NodeRenderable struct {
	nodes    []*html.Node
	text     string
	textFlag bool

	content     string
	contentFlag bool
}

func (n *NodeRenderable) GetNodes() []*html.Node {
	return n.nodes
}

func NewNodeRenderable(node []*html.Node) *NodeRenderable {
	return &NodeRenderable{
		nodes:       node,
		text:        "",
		textFlag:    false,
		content:     "",
		contentFlag: false,
	}
}

func NewTitleRenderable(title string) *NodeRenderable {
	strong := &html.Node{Type: html.ElementNode, Data: "strong"}
	strong.AppendChild(&html.Node{Type: html.TextNode, Data: title})
	node := &html.Node{Type: html.ElementNode, Data: "h1"}
	node.AppendChild(strong)
	return NewNodeRenderable([]*html.Node{node})
}

func NewParamRenderable(sd []*StructDef, colorset *ColorSet, prefix string) *NodeRenderable {
	if sd == nil || len(sd) == 0 {
		return nil
	}

	for _, s := range sd {
		if s == nil {
			continue
		}
		// spec color to relate struct
		colorset.Get(s.Name)
	}

	// header
	sd[0].Prefix = prefix
	ret := NewNodeRenderable(sd[0].GetNodes(colorset))
	for i := 1; i < len(sd); i++ {
		ret.Append(sd[i].GetNodes(colorset))
	}
	return ret
}

func (n *NodeRenderable) Rerender() {
	contentBuilder := strings.Builder{}

	for i := 0; i < len(n.nodes); i++ {
		html.Render(&contentBuilder, n.nodes[i])
	}

	n.contentFlag = true
	n.content = contentBuilder.String()
}

func (n *NodeRenderable) Render() string {
	if !n.contentFlag {
		n.Rerender()
	}
	return n.content
}

func (n *NodeRenderable) RenderWithBuilder(b *strings.Builder) {
	if !n.contentFlag {
		n.Rerender()
	}
	b.WriteString(n.content)
}

func (n *NodeRenderable) Text() string {
	if n.textFlag {
		return n.text
	}

	// build string by search text nodes
	keyBuilder := strings.Builder{}
	for i := 0; i < len(n.nodes); i++ {
		var node *html.Node
		stack := []*html.Node{n.nodes[i]}
		for len(stack) > 0 {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if node.Type == html.TextNode {
				keyBuilder.WriteString(node.Data)
			}

			// append to stack
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				stack = append(stack, child)
			}
		}
	}

	n.textFlag = true
	n.text = strings.TrimSpace(keyBuilder.String())

	return n.text
}

func (n *NodeRenderable) Append(node []*html.Node) {
	n.nodes = append(n.nodes, node...)
	n.textFlag = false
	n.contentFlag = false
}
