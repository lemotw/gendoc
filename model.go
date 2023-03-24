package gendoc

import "github.com/lemotw/gendoc/model"

type ConflunceSetting struct {
	Auth     model.ConflunceAccount
	Domain   string
	SpaceKey string
}

type ApiInfo struct {
	ApiURL    string
	ApiMethod string
	Intro     string
}

type Param struct {
	Data       interface{}
	JsonRender bool
}

type APIRegistry struct {
	ID       string
	ParentId string
	Info     *ApiInfo
	Req      *Param
	Res      *Param
}
