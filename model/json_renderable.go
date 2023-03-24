package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

const JSON_MARKDOWN_TEMPLATE = "<ac:structured-macro ac:name=\"markdown\"><ac:plain-text-body><![CDATA[```json \n%s\n```]]></ac:plain-text-body></ac:structured-macro>"

type JsonContent struct {
	Data interface{}
}

func (c *JsonContent) Text() string {
	jsonStr, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		jsonStr = []byte("")
	}
	return string(jsonStr)
}

func (c *JsonContent) Render() string {
	jsonStr, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		jsonStr = []byte("")
	}

	return fmt.Sprintf(JSON_MARKDOWN_TEMPLATE, string(jsonStr))
}
func (c *JsonContent) RenderWithBuilder(b *strings.Builder) {
	jsonStr, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		jsonStr = []byte("")
	}

	b.WriteString(fmt.Sprintf(JSON_MARKDOWN_TEMPLATE, string(jsonStr)))
}
