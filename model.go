package gendoc

import "github.com/lemotw/gendoc/model"

const (
	INTRODUTION_HEADER = "Introduction"
	REQUEST_HEADER     = "Request"
	RESPONSE_HEADER    = "Response"
)

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

func (info *ApiInfo) GetTitle() string {
	return "[" + info.Intro + "] " + info.ApiMethod + ": " + info.ApiURL
}

type Param struct {
	Data       interface{}
	JsonRender bool
}

type APIRegistry struct {
	ID       string
	Version  int
	ParentId string
	Info     *ApiInfo
	Req      *Param
	Res      *Param
}
