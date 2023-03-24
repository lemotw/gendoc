package analysis

import (
	"errors"
	"strings"

	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

func ParseDoc(doc string) (*model.Doc, error) {
	ret := &model.Doc{}

	node, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		panic(err)
	}

	bodyNode := FindNode(node, "body")
	if bodyNode == nil {
		return nil, errors.New("body node not found")
	}

	var prevKey *model.NodeRenderable
	var contentBuf []*html.Node
	for n := bodyNode.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode && n.Data == "h1" {
			if prevKey != nil {
				// append doc
				ret.AppendRow(prevKey, model.NewNodeRenderable(contentBuf))
				contentBuf = []*html.Node{}
			}

			prevKey = model.NewNodeRenderable([]*html.Node{n})
		} else {
			contentBuf = append(contentBuf, n)
		}
	}

	// last key
	if len(contentBuf) > 0 {
		if prevKey == nil {
			prevKey = model.NewTitleRenderable("")
		}
		ret.AppendRow(prevKey, model.NewNodeRenderable(contentBuf))
	}

	return ret, nil
}
