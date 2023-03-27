package gendoc

import (
	"regexp"

	"github.com/lemotw/gendoc/analysis"
	"github.com/lemotw/gendoc/conflunce"
	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

type IDocService interface {
	NewGroup(pageId string) *APIGroup

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

func (serv *docService) NewGroup(pageId string) *APIGroup {
	return &APIGroup{
		ParentId: pageId,
		serv:     serv,
	}
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

func fetchRegexMatchPage(existPageMap map[string]*APIRegistry) func(child *model.ChildPage) error {
	return func(child *model.ChildPage) error {
		regex := regexp.MustCompile("(\\[(.*)\\])\\ *((?i:GET|POST|PUT|DEL)):\\ *(\\/.*)")
		titleInfo := regex.FindStringSubmatch(child.Title)
		if titleInfo == nil {
			return nil
		}

		if len(titleInfo) == 5 {
			existPageMap[titleInfo[3]+":"+titleInfo[4]] = &APIRegistry{
				ID:      child.ID,
				Version: child.Version.Number,
				Info: &ApiInfo{
					ApiURL:    titleInfo[4],
					ApiMethod: titleInfo[3],
					Intro:     titleInfo[2],
				},
			}
		}

		return nil
	}
}

func (serv *docService) createPage(parentID string, api *APIRegistry) error {
	// create doc
	doc := &model.Doc{}

	// set introdution
	introdution := model.NewNodeRenderable([]*html.Node{{
		Type: html.ElementNode,
		Data: "p",
	}})
	introdution.Append([]*html.Node{{Type: html.TextNode, Data: api.Info.Intro}})
	doc.AppendRow(model.NewTitleRenderable(INTRODUTION_HEADER), introdution)

	// set req
	reqKey := model.NewTitleRenderable(REQUEST_HEADER)

	req := model.NewParamRenderable(analysis.ReflectAny(api.Req.Data, serv.descriptionMap, "req"))
	doc.Append(reqKey, req)
	if api.Req.JsonRender {
		doc.Append(reqKey, model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
		doc.Append(reqKey, &model.JsonContent{Data: api.Req.Data})
	}

	// set res
	resKey := model.NewTitleRenderable(RESPONSE_HEADER)

	res := model.NewParamRenderable(analysis.ReflectAny(api.Res.Data, serv.descriptionMap, "res"))
	doc.Append(resKey, res)
	if api.Res.JsonRender {
		doc.Append(resKey, model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
		doc.Append(resKey, &model.JsonContent{Data: api.Res.Data})
	}

	// 2. create the page
	page := model.NewPage(doc, parentID, api.Info.GetTitle(), serv.ConflunceSetting.SpaceKey)
	err := conflunce.NewConfluncePage(serv.ConflunceSetting.Domain, page, &serv.ConflunceSetting.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (serv *docService) updatePage(parentID string, existPage, api *APIRegistry) error {
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

	// 3. update the api (req and res)
	var reqKey model.Renderable
	var resKey model.Renderable
	for i := 0; i < len(doc.KeySequence); i++ {
		// find key
		if doc.KeySequence[i].Text() == REQUEST_HEADER {
			reqKey = doc.KeySequence[i]
		}

		if doc.KeySequence[i].Text() == RESPONSE_HEADER {
			resKey = doc.KeySequence[i]
		}
	}

	// req set
	if reqKey == nil {
		doc.Append(model.NewTitleRenderable(REQUEST_HEADER), model.NewParamRenderable(analysis.ReflectAny(api.Req.Data, serv.descriptionMap, "req")))
	} else {
		doc.Contents[reqKey] = []model.Renderable{model.NewParamRenderable(analysis.ReflectAny(api.Req.Data, serv.descriptionMap, "req"))}
		if api.Req.JsonRender {
			doc.Contents[reqKey] = append(doc.Contents[reqKey], model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
			doc.Contents[reqKey] = append(doc.Contents[reqKey], &model.JsonContent{Data: api.Req.Data})
		}
	}

	// res set
	if resKey == nil {
		doc.Append(model.NewTitleRenderable(RESPONSE_HEADER), model.NewParamRenderable(analysis.ReflectAny(api.Res.Data, serv.descriptionMap, "res")))
	} else {
		doc.Contents[resKey] = []model.Renderable{model.NewParamRenderable(analysis.ReflectAny(api.Res.Data, serv.descriptionMap, "res"))}
		if api.Req.JsonRender {
			doc.Contents[resKey] = append(doc.Contents[resKey], model.NewNodeRenderable([]*html.Node{{Type: html.ElementNode, Data: "br"}}))
			doc.Contents[resKey] = append(doc.Contents[resKey], &model.JsonContent{Data: api.Res.Data})
		}
	}

	// 4. update the page
	page := model.NewPage(doc, parentID, api.Info.GetTitle(), serv.ConflunceSetting.SpaceKey)
	page.Version.Number = existPage.Version + 1
	err = conflunce.PutConfluncePage(serv.ConflunceSetting.Domain, existPage.ID, page, &serv.ConflunceSetting.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (serv *docService) GenDoc() error {
	for pageId, apis := range serv.Registry {
		// make exist page map
		existPageMap := make(map[string]*APIRegistry)
		if err := conflunce.WalkAllChildPage(serv.ConflunceSetting.Domain, pageId, fetchRegexMatchPage(existPageMap), &serv.ConflunceSetting.Auth); err != nil {
			return err
		}

		for _, api := range apis {
			existPage, ok := existPageMap[api.Info.ApiMethod+":"+api.Info.ApiURL]
			if ok {
				if err := serv.updatePage(pageId, existPage, api); err != nil {
					return err
				}
			}

			if err := serv.createPage(pageId, api); err != nil {
				return err
			}
		}
	}

	return nil
}
