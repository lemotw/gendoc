package model

import "strings"

type Renderable interface {
	Text() string
	Render() string
	RenderWithBuilder(b *strings.Builder)
}

type Doc struct {
	Contents    map[Renderable][]Renderable
	KeySequence []Renderable
}

func (doc *Doc) AppendRow(key Renderable, content Renderable) {
	if doc.Contents == nil {
		doc.Contents = make(map[Renderable][]Renderable)
		doc.Contents[key] = []Renderable{content}
	} else {
		doc.Contents[key] = []Renderable{content}
	}

	if doc.KeySequence == nil {
		doc.KeySequence = []Renderable{key}
	} else {
		doc.KeySequence = append(doc.KeySequence, key)
	}
}

func (doc *Doc) Append(key Renderable, content Renderable) {
	if _, ok := doc.Contents[key]; ok {
		doc.Contents[key] = append(doc.Contents[key], content)
	} else {
		doc.AppendRow(key, content)
	}
}

func (doc *Doc) Set(key Renderable, content Renderable) {
	if _, ok := doc.Contents[key]; ok {
		doc.Contents[key] = []Renderable{content}
	} else {
		doc.AppendRow(key, content)
	}
}

func (doc *Doc) Render() string {
	docBuilder := strings.Builder{}
	for i := 0; i < len(doc.KeySequence); i++ {
		if content, ok := doc.Contents[doc.KeySequence[i]]; ok {
			doc.KeySequence[i].RenderWithBuilder(&docBuilder)
			for j := 0; j < len(content); j++ {
				content[j].RenderWithBuilder(&docBuilder)
			}
		}
	}

	return docBuilder.String()
}
