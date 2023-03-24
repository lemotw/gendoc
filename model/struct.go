package model

import (
	"golang.org/x/net/html"
)

type StructField struct {
	Name string
	Type string
	Req  bool
	Desc string
}

func (field *StructField) GetNode() *html.Node {
	trNode := &html.Node{Type: html.ElementNode, Data: "tr"}

	// name
	td1Node := &html.Node{Type: html.ElementNode, Data: "td"}
	td1Node.AppendChild(&html.Node{Type: html.TextNode, Data: field.Name})
	trNode.AppendChild(td1Node)

	// type
	td2Node := &html.Node{Type: html.ElementNode, Data: "td"}
	td2Node.AppendChild(&html.Node{Type: html.TextNode, Data: field.Type})
	trNode.AppendChild(td2Node)

	// req
	td3Node := &html.Node{Type: html.ElementNode, Data: "td"}
	if field.Req {
		td3Node.AppendChild(&html.Node{Type: html.TextNode, Data: "Y"})
	} else {
		td3Node.AppendChild(&html.Node{Type: html.TextNode, Data: "N"})
	}
	trNode.AppendChild(td3Node)

	// desc
	td4Node := &html.Node{Type: html.ElementNode, Data: "td"}
	td4Node.AppendChild(&html.Node{Type: html.TextNode, Data: field.Desc})
	trNode.AppendChild(td4Node)

	return trNode
}

type StructTable []*StructField

type StructDef struct {
	Name   string
	Fields StructTable
}

func (def *StructDef) GetNodes() []*html.Node {
	// Structure Name
	h2 := &html.Node{Type: html.ElementNode, Data: "h2"}
	strong := &html.Node{Type: html.ElementNode, Data: "strong"}
	strong.AppendChild(&html.Node{Type: html.TextNode, Data: def.Name})
	h2.AppendChild(strong)

	// table
	attrs := []html.Attribute{{Key: "border", Val: "1"}, {Key: "cellspacing", Val: "0"}, {Key: "cellpadding", Val: "5"}}

	tableNode := &html.Node{Type: html.ElementNode, Data: "table", Attr: attrs}
	theadNode := &html.Node{Type: html.ElementNode, Data: "thead"}
	trNode := &html.Node{Type: html.ElementNode, Data: "tr"}

	th1Node := &html.Node{Type: html.ElementNode, Data: "th"}
	th1Node.AppendChild(&html.Node{Type: html.TextNode, Data: "Name"})
	trNode.AppendChild(th1Node)

	th2Node := &html.Node{Type: html.ElementNode, Data: "th"}
	th2Node.AppendChild(&html.Node{Type: html.TextNode, Data: "Type"})
	trNode.AppendChild(th2Node)

	th3Node := &html.Node{Type: html.ElementNode, Data: "th"}
	th3Node.AppendChild(&html.Node{Type: html.TextNode, Data: "Required"})
	trNode.AppendChild(th3Node)

	th4Node := &html.Node{Type: html.ElementNode, Data: "th"}
	th4Node.AppendChild(&html.Node{Type: html.TextNode, Data: "Description"})
	trNode.AppendChild(th4Node)

	theadNode.AppendChild(trNode)
	tableNode.AppendChild(theadNode)

	for _, field := range def.Fields {
		tableNode.AppendChild(field.GetNode())
	}

	return []*html.Node{h2, tableNode}
}
