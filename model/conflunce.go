package model

type ConflunceAccount struct {
	Username string
	Password string
}

// page
type Ancestor struct {
	ID string `json:"id"`
}

type Space struct {
	Key string `json:"key"`
}

type Body struct {
	Storage Storage `json:"storage"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Version struct {
	Number int `json:"number"`
}

type Page struct {
	Title     string     `json:"title"`
	Version   Version    `json:"version"`
	Type      *string    `json:"type"`
	Ancestors []Ancestor `json:"ancestors"`
	Space     Space      `json:"space"`
	Body      Body       `json:"body"`
}

func NewPage(doc *Doc, parentId, title, spaceKey string) *Page {
	t := "page"
	return &Page{
		Title: title,
		Type:  &t,
		Body: Body{
			Storage: Storage{
				Value:          doc.Render(),
				Representation: "storage",
			},
		},
		Version: Version{
			Number: 0,
		},
		Space: Space{
			Key: spaceKey,
		},
		Ancestors: []Ancestor{
			{ID: parentId},
		},
	}
}

//func (page *Page) ToDoc() Doc {
//	htmlFragment, err := html.ParseFragment(strings.NewReader(page.Body.Storage.Value), nil)
//	if err != nil {
//		panic(err)
//	}
//
//	var bodyNode *html.Node
//	func() {
//		for i := 0; i < len(htmlFragment); i++ {
//			if bodyNode = findElementWithDFS(htmlFragment[i], "body"); bodyNode != nil {
//				return
//			}
//		}
//	}()
//
//	ret := Doc{Contents: make(map[*DocKey]Renderable)}
//	var prevKey *DocKey
//	var contentBuf []*html.Node
//	for n := bodyNode.FirstChild; n != nil; n = n.NextSibling {
//		// every node
//		if n.Type == html.ElementNode && n.Data == "h1" {
//			// document key
//			if prevKey != nil {
//				ret.Contents[prevKey] = NewDocContent(contentBuf)
//				contentBuf = contentBuf[:0]
//			}
//
//			prevKey = &DocKey{}
//			prevKey.Reset(n)
//			ret.KeySequence = append(ret.KeySequence, prevKey)
//		} else {
//			// document content
//			contentBuf = append(contentBuf, n)
//		}
//	}
//
//	// last key
//	if len(contentBuf) > 0 {
//		ret.Contents[prevKey] = NewDocContent(contentBuf)
//	}
//
//	return ret
//}

// child
type ChildPageResponse struct {
	Results    []ChildPage `json:"results"`
	Size       int         `json:"size"`
	Limit      int         `json:"limit"`
	IsLastPage bool        `json:"isLastPage"`
	Start      int         `json:"start"`
}

type ChildPage struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Version Version `json:"version"`
	Title   string  `json:"title"`
	Status  string  `json:"status"`
}
