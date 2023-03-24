package gendoc

import (
	"log"
	"regexp"

	"github.com/lemotw/gendoc/analysis"
	"github.com/lemotw/gendoc/conflunce"
	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

type IDocService interface {
	RegisterAPI(pageId string, apiInfo *ApiInfo, req *Param, res *Param)

	GenDoc() error
}

func NewDocService(setting *ConflunceSetting) (IDocService, error) {
	descMap, err := analysis.GetDescriptionMap()
	if err != nil {
		return nil, err
	}
	serv := &docService{
		ConflunceSetting: setting,
		Registry:         make(map[string][]*APIRegistry),
		descriptionMap:   descMap,
	}

	return serv, nil
}

type docService struct {
	ConflunceSetting *ConflunceSetting

	Registry       map[string][]*APIRegistry
	descriptionMap map[string]map[string]string
}

func (serv *docService) RegisterAPI(pageId string, apiInfo *ApiInfo, req *Param, res *Param) {
	registry := &APIRegistry{
		ParentId: pageId,
		Info:     apiInfo,
		Req:      req,
		Res:      res,
	}

	if _, ok := serv.Registry[pageId]; !ok {
		serv.Registry[pageId] = []*APIRegistry{registry}
	} else {
		serv.Registry[pageId] = append(serv.Registry[pageId], registry)
	}
}

func (serv *docService) GenDoc() error {
	for pageId, apis := range serv.Registry {
		alreadyExistApiPage := make(map[string]*APIRegistry) // api key -> api registry

		walkFunc := func(child *model.ChildPage) bool {
			regex := regexp.MustCompile("(\\[(.*)\\])\\ *((?i:GET|POST|PUT|DEL)):\\ *(\\/.*)")
			titleInfo := regex.FindStringSubmatch(child.Title)
			if titleInfo == nil {
				return false
			}

			if len(titleInfo) == 5 {
				alreadyExistApiPage[titleInfo[3]+":"+titleInfo[4]] = &APIRegistry{
					ID: child.ID,
					Info: &ApiInfo{
						ApiURL:    titleInfo[4],
						ApiMethod: titleInfo[3],
						Intro:     titleInfo[2],
					},
				}
			}

			return false
		}

		if err := conflunce.WalkAllChildPage(serv.ConflunceSetting.Domain, pageId, walkFunc, &serv.ConflunceSetting.Auth); err != nil {
			return err
		}

		for _, api := range apis {
			existPage, ok := alreadyExistApiPage[api.Info.ApiMethod+":"+api.Info.ApiURL]
			if ok {
				// update
				// 1. get page and parse into model.Doc
				p, err := conflunce.FetchConfluncePage(serv.ConflunceSetting.Domain, existPage.ID, &serv.ConflunceSetting.Auth)
				if err != nil {
					return err
				}

				// 2. reserve the old content except the api (req and res)
				doc, err := analysis.ParseDoc(p.Body.Storage.Value)
				if err != nil {
					return err
				}
				log.Println(doc)

				// 3. update the api (req and res)
				keyS := []string{}
				for i := 0; i < len(doc.KeySequence); i++ {
					keyS = append(keyS, doc.KeySequence[i].Text())
				}
				log.Println(keyS)

				// 4. update the page

			} else {
				// create doc
				doc := &model.Doc{}

				// set introdution
				introdution := model.NewNodeRenderable([]*html.Node{{
					Type: html.ElementNode,
					Data: "p",
				}})
				introdution.Append([]*html.Node{{Type: html.TextNode, Data: api.Info.Intro}})
				doc.AppendRow(model.NewTitleRenderable("Introduction"), introdution)

				// 1. create the api (req and res)

				// set req
				reqContentStruct, reqRelateStruct := analysis.ReflectStruct(api.Req.Data, serv.descriptionMap)
				reqContent := model.NewNodeRenderable(reqContentStruct.GetNodes())
				for i := 0; i < len(reqRelateStruct); i++ {
					reqContent.Append(reqRelateStruct[i].GetNodes())
				}

				reqKey := model.NewTitleRenderable("Request")
				doc.Append(reqKey, reqContent)
				if api.Req.JsonRender {
					doc.Append(reqKey, model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
					doc.Append(reqKey, &model.JsonContent{Data: api.Req.Data})
				}

				// set res
				resContentStruct, resRelateStruct := analysis.ReflectStruct(api.Res.Data, serv.descriptionMap)
				resContent := model.NewNodeRenderable(resContentStruct.GetNodes())
				for i := 0; i < len(resRelateStruct); i++ {
					resContent.Append(resRelateStruct[i].GetNodes())
				}

				resKey := model.NewTitleRenderable("Response")
				doc.Append(resKey, resContent)
				if api.Res.JsonRender {
					doc.Append(resKey, model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
					doc.Append(resKey, &model.JsonContent{Data: api.Res.Data})
				}

				// 2. create the page
				//model.Page{}
				page := model.NewPage(doc, pageId, serv.ConflunceSetting.SpaceKey)
				page.Title = "[" + api.Info.Intro + "] " + api.Info.ApiMethod + ": " + api.Info.ApiURL
				err := conflunce.NewConfluncePage(serv.ConflunceSetting.Domain, page, &serv.ConflunceSetting.Auth)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
