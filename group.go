package gendoc

type APIGroup struct {
	ParentId string
	serv     IDocService
}

func (group *APIGroup) RegisterAPI(method, url, intro string, req *Param, res *Param) {
	group.serv.RegisterAPI(group.ParentId, &ApiInfo{
		ApiURL:    url,
		ApiMethod: method,
		Intro:     intro,
	}, req, res)
}
