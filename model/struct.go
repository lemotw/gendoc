package model

import (
	"strings"

	"golang.org/x/net/html"
)

type StructField struct {
	Name string
	Type string
	Req  bool
	Desc []string
}

func (field *StructField) GetNode(colorset *ColorSet) *html.Node {
	trNode := &html.Node{Type: html.ElementNode, Data: "tr"}

	// name
	td1Node := &html.Node{Type: html.ElementNode, Data: "td"}
	td1Node.AppendChild(&html.Node{Type: html.TextNode, Data: field.Name})
	trNode.AppendChild(td1Node)

	// type
	colorAttr := []html.Attribute{}
	rawType := field.Type
	rawType = strings.Replace(rawType, "*", "", -1)
	rawType = strings.Replace(rawType, "[]", "", -1)
	typeColor := colorset.TryGet(rawType)
	if typeColor != nil {
		colorAttr = append(colorAttr, html.Attribute{Key: "style", Val: "color:" + typeColor.Hex()})
	}

	td2Node := &html.Node{Type: html.ElementNode, Data: "td", Attr: colorAttr}
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
	for i := 0; i < len(field.Desc); i++ {
		pNode := html.Node{Type: html.ElementNode, Data: "p"}
		pNode.AppendChild(&html.Node{Type: html.TextNode, Data: field.Desc[i]})
		td4Node.AppendChild(&pNode)
	}
	trNode.AppendChild(td4Node)

	return trNode
}

type StructTable []*StructField

type StructDef struct {
	Name   string
	Prefix string
	Fields StructTable
}

func (def *StructDef) HeaderStr() string {
	if len(def.Prefix) > 0 {
		return "[" + def.Prefix + "]" + def.Name
	}

	return def.Name
}

func (def *StructDef) GetNodes(colorset *ColorSet) []*html.Node {
	// color
	sColorAttr := []html.Attribute{}
	if len(def.Prefix) == 0 {
		sColor := colorset.Get(def.Name)
		if sColor != nil {
			sColorAttr = append(sColorAttr, html.Attribute{Key: "style", Val: "color:" + sColor.Hex()})
		}
	}

	// struct header
	h2 := &html.Node{Type: html.ElementNode, Data: "h2", Attr: sColorAttr}
	strong := &html.Node{Type: html.ElementNode, Data: "strong"}
	strong.AppendChild(&html.Node{Type: html.TextNode, Data: def.HeaderStr()})
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
		tableNode.AppendChild(field.GetNode(colorset))
	}

	return []*html.Node{h2, tableNode}
}
